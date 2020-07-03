package finder

import (
	"fmt"
	"math"
	"testing"

	"github.com/xrash/smetrics"
)

func equal(a, b float64) bool {
	const radix = 0.0005

	if a > b {
		return a-b < radix
	}

	return b-a < radix
}

func TestJaroImplementations(t *testing.T) {

	tests := []struct {
		a     string
		b     string
		score float64
	}{
		{a: "CRATE", b: "TRACE", score: 0.733333},
		{a: "MARTHA", b: "MARHTA", score: 0.944444},
		{a: "DIXON", b: "DICKSONX", score: 0.766666},
		{a: "gmilcon", b: "gmilno", score: 0.896825},
	}
	for _, tt := range tests {
		score := NewJaro()(tt.a, tt.b)
		if !equal(tt.score, score) {
			t.Errorf("Expected a score of %f, instead it was %f for input, a: %q, b: %q ", tt.score, score, tt.a, tt.b)
		}

		t.Logf("%q vs. %q", tt.a, tt.b)
		t.Logf("smetrics.Jaro        %f", smetrics.Jaro(tt.a, tt.b))
		t.Logf("RosettaJaroV0        %f", RosettaJaroV0(tt.a, tt.b))
		t.Logf("RosettaJaroV1        %f", RosettaJaroV1(tt.a, tt.b))
		t.Logf("JaroDistanceMasatana %f", func() float64 {
			s, _ := JaroDistanceMasatana(tt.a, tt.b)
			return s
		}())
	}
}

func TestComparingJaroImplementations(t *testing.T) {

	algos := []Algorithm{
		// Rosetta Jaro variants
		RosettaJaroV0,
		RosettaJaroV1,

		// smetrics impementatio, currently fails
		// smetrics.Jaro,

		// Masatana's implementation, slightly different API
		func(a, b string) float64 {
			s, _ := JaroDistanceMasatana(a, b)
			return s
		},
	}

	sets := []struct {
		a string
		b string
	}{
		{a: "aaaaaa", b: "zzzzzz"},
		{a: "beer", b: "root"},
		{a: "beer", b: "been"},
		{a: "huffelpuf", b: "puffelhuf"},
		{a: "algorithm", b: "algoritm"},
		{a: "corn", b: "corm"},
	}

	for _, set := range sets {

		var score float64 = 0
		for _, algo := range algos {
			subScore := algo(set.a, set.b)
			if score == 0 {
				score = subScore
				continue
			}

			if !equal(score, subScore) {
				t.Errorf("Algorithm disagreement for %q and %q, score %f subscore %f", set.a, set.b, score, subScore)
			}
		}
	}
}

// From: github.com/masatana/go-textdistance
func JaroDistanceMasatana(s1, s2 string) (float64, int) {
	if s1 == s2 {
		return 1.0, 0.0
	}
	// compare length using rune slice length, as s1 and s2 are not necessarily ASCII-only strings
	longer, shorter := []rune(s1), []rune(s2)
	if len(longer) < len(shorter) {
		longer, shorter = shorter, longer
	}
	scope := int(math.Floor(float64(len(longer)/2))) - 1
	// m is the number of matching characters.
	m := 0
	matchFlags := make([]bool, len(longer))
	matchIndexes := make([]int, len(longer))
	for i := range matchIndexes {
		matchIndexes[i] = -1
	}

	for i := 0; i < len(shorter); i++ {
		k := Min(i+scope+1, len(longer))
		for j := Max(i-scope, 0); j < k; j++ {
			if matchFlags[j] || shorter[i] != longer[j] {
				continue
			}
			matchIndexes[i] = j
			matchFlags[j] = true
			m++
			break
		}
	}
	ms1 := make([]rune, m)
	ms2 := make([]rune, m)
	si := 0
	for i := 0; i < len(shorter); i++ {
		if matchIndexes[i] != -1 {
			ms1[si] = shorter[i]
			si++
		}
	}
	si = 0
	for i := 0; i < len(longer); i++ {
		if matchFlags[i] {
			ms2[si] = longer[i]
			si++
		}
	}

	t := 0
	for i, c := range ms1 {
		if c != ms2[i] {
			t++
		}
	}
	prefix := 0
	for i := 0; i < len(shorter); i++ {
		if longer[i] == shorter[i] {
			prefix++
		} else {
			break
		}
	}
	if m == 0 {
		return 0.0, 0.0
	}
	newt := float64(t) / 2.0
	newm := float64(m)
	return 1 / 3.0 * (newm/float64(len(shorter)) + newm/float64(len(longer)) + (newm-newt)/newm), prefix
}

func Min(is ...int) int {
	var min int
	for i, v := range is {
		if i == 0 || v < min {
			min = v
		}
	}
	return min
}

// Max returns the maximum number of passed int slices.
func Max(is ...int) int {
	var max int
	for _, v := range is {
		if max < v {
			max = v
		}
	}
	return max
}

func BenchmarkJaroImplementations(b *testing.B) {
	sets := []struct {
		a string
		b string
	}{
		{a: "gmilcon", b: "gmilno"},
		{a: "DIXON", b: "DICKSONX"},
		{a: "MARHTA", b: "martha"},
	}

	for _, set := range sets {

		b.Run(fmt.Sprintf("%q %q", set.a, set.b), func(b *testing.B) {

			b.Run("JaroDistanceMasatana", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_, _ = JaroDistanceMasatana(set.a, set.b)
				}
			})

			b.Run("RosettaJaro V0", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = RosettaJaroV0(set.a, set.b)
				}
			})

			b.Run("RosettaJaro V1", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = RosettaJaroV1(set.a, set.b)
				}
			})

			b.Run("MyJaro", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = smetrics.Jaro(set.a, set.b)
				}
			})
		})
	}
}

func BenchmarkRosettaJaro(b *testing.B) {
	sets := []struct {
		a string
		b string
	}{
		{a: "aaaaaa", b: "zzzzzz"},
		{a: "beer", b: "root"},
		{a: "beer", b: "been"},
		{a: "huffelpuf", b: "puffelhuf"},
		{a: "algorithm", b: "algoritm"},
		{a: "corn", b: "corm"},
	}

	b.Run("Double alloc", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range sets {
				_ = RosettaJaroV0(s.a, s.b)
			}
		}
	})
	b.Run("Single alloc", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, s := range sets {
				_ = RosettaJaroV1(s.a, s.b)
			}
		}
	})
}

// @see https://rosettacode.org/wiki/Jaro_distance#Go
func RosettaJaroV0(a, b string) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1
	}
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	matchDistance := len(a)
	if len(b) > matchDistance {
		matchDistance = len(b)
	}

	matchDistance = matchDistance/2 - 1
	aMatches := make([]bool, len(a))
	bMatches := make([]bool, len(b))

	var matches float64
	var transpositions float64
	for i := range a {
		start := i - matchDistance
		if start < 0 {
			start = 0
		}

		end := i + matchDistance + 1
		if end > len(b) {
			end = len(b)
		}

		for k := start; k < end; k++ {
			if bMatches[k] {
				continue
			}
			if a[i] != b[k] {
				continue
			}

			aMatches[i] = true
			bMatches[k] = true
			matches++
			break
		}
	}

	if matches == 0 {
		return 0
	}

	k := 0
	for i := range a {
		if !aMatches[i] {
			continue
		}

		for !bMatches[k] {
			k++
		}

		if a[i] != b[k] {
			transpositions++
		}

		k++
	}

	return (matches/float64(len(a)) +
		matches/float64(len(b)) +
		(matches-(transpositions/2))/matches) / 3
}

// @see https://rosettacode.org/wiki/Jaro_distance#Go
// Changes:
// - Minor allocation improvement
func RosettaJaroV1(a, b string) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1
	}
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	matchDistance := len(a)
	if len(b) > matchDistance {
		matchDistance = len(b)
	}

	matchDistance = matchDistance/2 - 1
	matchesCollected := make([]bool, len(a)+len(b))

	var matches float64
	var transpositions float64
	for i := range a {
		start := i - matchDistance
		if start < 0 {
			start = 0
		}

		end := i + matchDistance + 1
		if end > len(b) {
			end = len(b)
		}

		for k := start; k < end; k++ {
			if matchesCollected[k+len(a)] {
				continue
			}
			if a[i] != b[k] {
				continue
			}

			matchesCollected[i] = true
			matchesCollected[k+len(a)] = true
			matches++
			break
		}
	}

	if matches == 0 {
		return 0
	}

	k := 0
	for i := range a {
		if !matchesCollected[i] {
			continue
		}

		for !matchesCollected[k+len(a)] {
			k++
		}

		if a[i] != b[k] {
			transpositions++
		}

		k++
	}

	return (matches/float64(len(a)) +
		matches/float64(len(b)) +
		(matches-(transpositions/2))/matches) / 3
}
