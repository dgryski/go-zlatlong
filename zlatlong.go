// Package zlatlong implements Microsoft's lat/long compression algorithm for Bing paths
package zlatlong

import (
	"errors"
	"math"
)

type Point struct {
	Lat, Long float64
}

var safeCharacters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-")

var safeIdx [256]byte

func init() {
	for i := range safeIdx {
		safeIdx[i] = 255
	}

	for i, c := range safeCharacters {
		safeIdx[c] = byte(i)
	}
}

// Unmarshal decodes a compressed set of lat/long points
func Unmarshal(value []byte) ([]Point, error) {

	// From https://docs.microsoft.com/en-us/bingmaps/spatial-data-services/geodata-api

	var points []Point

	index := 0
	xsum := int64(0)
	ysum := int64(0)
	max := int64(4294967296)

	for index < len(value) {

		n := int64(0)
		k := uint(0)

		for {
			if index >= len(value) {
				return nil, nil
			}

			b := int64(safeIdx[value[index]])
			index++
			if b == 255 {
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

		diagonal := int64((math.Sqrt(8*float64(n)+5) - 1) / 2)

		n -= diagonal * (diagonal + 1) / 2
		ny := n
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

func round(f float64) int64 {
	if f < 0 {
		return int64(f - 0.5)
	}

	return int64(f + 0.5)
}

// Marshal encodes a series of lat/long points
func Marshal(points []Point) []byte {

	// From http://msdn.microsoft.com/en-us/library/jj158958.aspx

	latitude := int64(0)
	longitude := int64(0)
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
