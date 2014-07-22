// Package zlatlong implements Microsoft's lat/long compression algorithm for Bing paths
package zlatlong

import (
	"bytes"
	"errors"
	"math"
)

type Point struct {
	Lat, Long float64
}

var safeCharacters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-")

// Unmarshal decodes a compressed set of lat/long points
func Unmarshal(value []byte) ([]Point, error) {

	// From http://jkebeck.wordpress.com/2013/06/25/retrieving-boundaries-from-the-bing-spatial-data-services-preview/

	var points []Point

	index := 0
	xsum := 0
	ysum := 0
	max := 4294967296

	for index < len(value) {

		n := 0
		k := uint(0)

		for {
			if index >= len(value) {
				return nil, nil
			}

			b := bytes.IndexByte(safeCharacters, value[index])
			index++
			if b == -1 {
				return nil, errors.New("invalid character")
			}

			tmp := (b & 31) * (1 << k)

			ht := tmp / max
			lt := tmp % max

			hn := n / max
			ln := n % max

			nl := (lt | ln)
			n = (ht|hn)*max + nl
			k += 5
			if b < 32 {
				break
			}
		}

		diagonal := int((math.Sqrt(8*float64(n)+5) - 1) / 2)

		n -= diagonal * (diagonal + 1) / 2
		ny := int(n)
		nx := diagonal - ny
		nx = (nx >> 1) ^ -(nx & 1)
		ny = (ny >> 1) ^ -(ny & 1)
		xsum += nx
		ysum += ny
		lat := float64(ysum) * 0.00001
		lon := float64(xsum) * 0.00001
		points = append(points, Point{lat, lon})
	}
	return points, nil
}

func round(f float64) int {
	if f < 0 {
		return int(f - 0.5)
	}

	return int(f + 0.5)
}

// Marshal encodes a series of lat/long points
func Marshal(points []Point) []byte {

	// From http://msdn.microsoft.com/en-us/library/jj158958.aspx

	latitude := 0
	longitude := 0
	var result []byte

	for _, point := range points {
		// step 2
		newLatitude := round(point.Lat * 100000)
		newLongitude := round(point.Long * 100000)

		// step 3
		dy := newLatitude - latitude
		dx := newLongitude - longitude
		latitude = newLatitude
		longitude = newLongitude

		// step 4 and 5
		dy = (dy << 1) ^ (dy >> 31)
		dx = (dx << 1) ^ (dx >> 31)

		// step 6
		index := ((dy + dx) * (dy + dx + 1) / 2) + dy

		for index > 0 {

			// step 7
			var rem = index & 31
			index = (index - rem) / 32

			// step 8
			if index > 0 {
				rem += 32
			}

			// step 9
			result = append(result, safeCharacters[rem])
		}
	}

	return result
}
