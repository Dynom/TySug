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
		{Input: "bee4", List: []string{"beef", "beer"}, Expect: "beer"},

		{Input: "bee5", List: []string{"beer", "beef"}, Expect: "beer"},
		{Input: "bee5", List: []string{"beef", "beer"}, Expect: "beer"},
		{Input: "bee5", List: []string{"beef", "beer", "beast"}, Expect: "beer"},
		{Input: "bee5", List: []string{"beef", "beer", "ben"}, Expect: "beer"},

		{Input: "beeg", List: []string{"beef", "beer", "ben"}, Expect: "beef"},
		{Input: "beev", List: []string{"beef", "beer", "ben"}, Expect: "beef"},
		{Input: "beeb", List: []string{"beef", "beer", "ben"}, Expect: "beef"},

		{Input: "applejuice", List: []string{"apple", "pqqamqzxom", "excellent"}, Expect: "apple"},

		{
			Input:  "bee5",
			Expect: "beer",
			List: []string{
				"veer", "neer", "heer", "geer", "bwer", "bser", "bder", "brer", "b4er", "b3er", "bewr", "besr", "bedr", "berr",
				"be4r", "be3r", "beee", "beed", "beer", "beef", "beet", "bee4", "eer", "ber", "ber", "eber", "bere", "bbeer",
				"beeer", "beeer", "beerr",
			},
		},
	}

	kd := New(QwertyUS)
	for _, td := range testData {
		result, distance := kd.FindNearest(td.Input, td.List)
		if td.Expect != result {
			t.Errorf("Expected '%s' to match '%s', instead I got '%s' with distance %f, %+v",
				td.Input, td.Expect, result, distance, td)
			for _, ref := range td.List {
				t.Logf("ref '%s', distance %f", ref, kd.CalculateDistance(td.Input, ref))
			}

		}
	}
}

func TestCalculateDistance(t *testing.T) {
	t.Skip("Not particularly useful, mostly for debugging")
	testData := []struct {
		Input    string
		Ref      string
		Distance float64
	}{
		// near misses
		{Input: "bee4", Ref: "beef", Distance: 2},
		{Input: "bee5", Ref: "beer", Distance: 1.41},
		{Input: "applejuice", Ref: "apple", Distance: 5}, // Only the letters are missing

		// Realistic typos
		{Input: "applejuice", Ref: "applejuis", Distance: 2.41},
		{Input: "applejuice", Ref: "applejuisc", Distance: 3.41},
		{Input: "applejuice", Ref: "qookehjyuse", Distance: 13.88}, // fingers misaligned with keyboard

		// exaggerated experiments
		{Input: "applejuice", Ref: "pqqamqzxom", Distance: 69.05}, // A large distance (on QuertyUS), same word lengths
		{Input: "applejuice", Ref: "excellent", Distance: 42.58},  // against a random word

	}

	kd := New(QwertyUS)
	for _, td := range testData {
		d := kd.CalculateDistance(td.Input, td.Ref)

		if math.Abs(d-td.Distance) > floatTolerance {
			t.Errorf("Expected %s vs. %s to have a score of %f, instead I got %f", td.Input, td.Ref, td.Distance, d)
		}
	}
}

func TestGetDistance(t *testing.T) {
	testData := []struct {
		A        coordinates
		B        coordinates
		Distance float64
	}{
		{A: coordinates{X: 0, Y: 0}, B: coordinates{X: 0, Y: 100}, Distance: 100},
		{A: coordinates{X: 0, Y: 0}, B: coordinates{X: 100, Y: 0}, Distance: 100},
		{A: coordinates{X: 1, Y: 2}, B: coordinates{X: 1, Y: 2}, Distance: 0},
		{A: coordinates{X: 10, Y: 20}, B: coordinates{X: 20, Y: 10}, Distance: 14.14},

		{A: coordinates{X: 20, Y: 0}, B: coordinates{X: 18, Y: 1}, Distance: 2.236}, // X: -2, Y +1
		{A: coordinates{X: 20, Y: 0}, B: coordinates{X: 19, Y: 2}, Distance: 2.236}, // X: -1, Y +2
	}

	for _, td := range testData {
		d := getDistance(td.A, td.B)

		if math.Abs(d-td.Distance) > floatTolerance {
			t.Errorf("Expected the distance to be %f, instead I got %f\n%v", td.Distance, d, td)
		}
	}
}

func TestFindNearestPenaltyScore(t *testing.T) {
	d := New(Default)
	_, s := d.FindNearest("foobar", []string{"f"})
	expected := 5 * missingCharPenalty

	if int(s) != expected {
		t.Errorf("Expected a penalty of %d, instead I got: %f", expected, s)
	}
}

func TestGenerateKeyDistance(t *testing.T) {
	grid := generateKeyGrid([]string{
		"abc",  // 00, 10, 20
		"def",  // 01, 11, 21
		"ghi",  // 02, 12, 22
		" jkl", // 03, 13, 23, 33 (leading space)
	})

	if c := grid["a"]; c.X != 0 || c.Y != 0 {
		t.Errorf("Expected the coords to be at 0,0 %+v", c)
	}

	if c := grid["e"]; c.X != 1 || c.Y != 1 {
		t.Errorf("Expected the coords to be at 1,1 %+v", c)
	}

	if c := grid["i"]; c.X != 2 || c.Y != 2 {
		t.Errorf("Expected the coords to be at 2,2 %+v", c)
	}

	if c := grid["j"]; c.X != 1 || c.Y != 3 {
		t.Errorf("Expected the coords to be at 3,1 %+v", c)
	}
}

func BenchmarkFindNearest(b *testing.B) {
	smallList := generateList(10)
	bigList := generateList(1000)
	b.Run("small-list", func(b *testing.B) {
		kd := New(Default)
		for i := 0; i < b.N; i++ {
			kd.FindNearest("minkey", smallList)
		}
	})

	b.Run("big-list", func(b *testing.B) {
		kd := New(Default)
		for i := 0; i < b.N; i++ {
			kd.FindNearest("minkey", bigList)
		}
	})
}

func BenchmarkAccessingStringCharacters(b *testing.B) {
	str := "42"

	b.ReportAllocs()

	// Reference test
	b.Run("staying a byte", func(b *testing.B) {
		var r uint8
		for i := 0; i < b.N; i++ {
			r = str[1]
		}

		_ = r
	})

	b.Run("with cast", func(b *testing.B) {
		var r string
		for i := 0; i < b.N; i++ {
			r = string(str[1])
		}

		_ = r
	})

	b.Run("with substr index", func(b *testing.B) {
		var r string
		for i := 0; i < b.N; i++ {
			r = str[1:2]
		}
		_ = r
	})
}

func BenchmarkCompByteWithStringIndex(b *testing.B) {
	base := "foo"
	ref := "foo"

	b.ReportAllocs()
	b.Run("comparing bytes", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if base[0] == ref[0] {
				continue
			}
			panic("shouldn't happen")
		}
	})
	b.Run("comparing strings", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if base[0:1] == ref[0:1] {
				continue
			}
			panic("shouldn't happen")
		}
	})
}

func generateList(size int) []string {
	list := make([]string, 0, size)

	for i := 0; i < size; i++ {
		list = append(list, "monkey")
	}

	return list
}
