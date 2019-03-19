package finder

import (
	"github.com/alextanhongpin/stringdist"
	"github.com/xrash/smetrics"
)

// Algorithm the type to comply with to create your own algorithm
// Note that the return value must be greater than WorstScoreValue and less than BestScoreValue
type Algorithm func(a, b string) float64

func NewJaroWinklerDefaults() Algorithm {
	return NewJaroWinkler(0.7, 4)
}

// NewJaroWinkler returns the JaroWinkler algorithm
func NewJaroWinkler(boostThreshold float64, prefixLength int) Algorithm {
	return func(a, b string) float64 {
		return smetrics.JaroWinkler(a, b, boostThreshold, prefixLength)
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
func NewJaro() Algorithm {
	return smetrics.Jaro
}
