package keyboard

import (
	"math"
	"testing"
)

const (
	floatTolerance = 0.01
)

func TestGetBestMatch(t *testing.T) {
	testData := []struct {
		Input  string
		List   []string
		Expect string
	}{
		{Input: "bee4", List: []string{"beer", "beef"}, Expect: "beer"},
		{Input: "bee5", List: []string{"beef", "beer"}, Expect: "beer"},
		{Input: "bee5", List: []string{"beef", "beer", "beast"}, Expect: "beer"},
		{Input: "bee5", List: []string{"beef", "beer", "ben"}, Expect: "beer"},
	}

	for _, td := range testData {
		result, distance := GetBestMatch(td.Input, td.List)
		if td.Expect != result {
			t.Errorf("Expected '%s' to match '%s', instead I got '%s' with distance %f, %+v",
				td.Input, td.Expect, result, distance, td)
		}
	}
}

func TestGetDistance(t *testing.T) {
	testData := []struct {
		A        Coordinates
		B        Coordinates
		Distance float64
	}{
		{A: Coordinates{X: 0, Y: 0}, B: Coordinates{X: 0, Y: 100}, Distance: 100},
		{A: Coordinates{X: 0, Y: 0}, B: Coordinates{X: 100, Y: 0}, Distance: 100},
		{A: Coordinates{X: 1, Y: 2}, B: Coordinates{X: 1, Y: 2}, Distance: 0},
		{A: Coordinates{X: 10, Y: 20}, B: Coordinates{X: 20, Y: 10}, Distance: 14.14},
	}

	for _, td := range testData {
		d := getDistance(td.A, td.B)

		if math.Abs(d-td.Distance) > floatTolerance {
			t.Errorf("Expected the distance to be %f, instead I got %f\n%v", td.Distance, d, td)
		}
	}

}

func TestGenerateKeyDistance(t *testing.T) {
	table := generateKeyDistance([]string{
		"abc", // 00, 01, 02
		"def", // 10, 11, 12
		"ghi", // 20, 21, 22
	})

	if table["a"].X != 0 || table["a"].Y != 0 {
		t.Errorf("Expected the coords to be at 0,0 %+v", table["a"])
	}

	if table["e"].X != 1 || table["e"].Y != 1 {
		t.Errorf("Expected the coords to be at 1,1 %+v", table["i"])
	}

	if table["i"].X != 2 || table["i"].Y != 2 {
		t.Errorf("Expected the coords to be at 2,2 %+v", table["i"])
	}
}
