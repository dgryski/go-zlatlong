package zlatlong

import (
	"math"
	"testing"
)

func TestEncode(t *testing.T) {

	const output = "vx1vilihnM6hR7mEl2Q"

	points := []Point{
		{35.894309002906084, -110.72522000409663},
		{35.893930979073048, -110.72577999904752},
		{35.893744984641671, -110.72606003843248},
		{35.893366960808635, -110.72661500424147},
	}

	r := Marshal(points)

	if string(r) != output {
		t.Fatalf("Marshall(points)=%q, want %q\n", string(r), output)
	}

	p2, _ := Unmarshal(r)

	for i := range points {
		if math.Abs(points[i].Lat-p2[i].Lat) > 1e-5 || math.Abs(points[i].Long-p2[i].Long) > 1e-5 {
			t.Errorf("Failed unpacking: got=%v want=%v", p2[i], points[i])
		}
	}
}
