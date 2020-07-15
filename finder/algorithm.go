package finder

import (
	"math"

	"github.com/alextanhongpin/stringdist"
	"github.com/xrash/smetrics"
)

// Algorithm the type to comply with to create your own algorithm
// Note that the return value must be greater than WorstScoreValue and less than BestScoreValue
type Algorithm func(a, b string) float64

// NewJaroWinklerDefaults returns the Jaro Winkler algorithm with 0.7 boost threshold and a prefix length of 4
func NewJaroWinklerDefaults() Algorithm {
	return NewJaroWinkler(0.7, 4)
}

// NewJaroWinkler returns the JaroWinkler algorithm
func NewJaroWinkler(boostThreshold float64, prefixLength int) Algorithm {

	// @see smetrics.Jaro() Duplicated here to reference a different local Jaro implementation
	return func(a, b string) float64 {
		j := NewJaro()(a, b)

		if j <= boostThreshold {
			return j
		}

		prefixLength = int(math.Min(float64(len(a)), math.Min(float64(prefixLength), float64(len(b)))))

		var prefixMatch float64
		for i := 0; i < prefixLength; i++ {
			if a[i] == b[i] {
				prefixMatch++
			}
		}

		return j + 0.1*prefixMatch*(1.0-j)
	}
}

// NewDamerauLevenshtein returns the DamerauLevenshtein algorithm
func NewDamerauLevenshtein() Algorithm {
	var dl = stringdist.NewTrueDamerauLevenshtein()
	return func(a, b string) float64 {
		return float64(-dl.Calculate(a, b))
	}
}

// NewWagnerFischer returns the NewWagnerFischer algorithm, sensible defaults are: i:1, d:3, s:1
func NewWagnerFischer(insert, delete, substitution int) Algorithm {
	return func(a, b string) float64 {
		return float64(-smetrics.WagnerFischer(a, b, insert, delete, substitution))
	}
}

// NewJaro returns the default Jaro algorithm
// @see https://rosettacode.org/wiki/Jaro_distance#Go
//nolint:gocyclo
func NewJaro() Algorithm {
	return func(a, b string) float64 {
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
}
