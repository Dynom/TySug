package main

import (
	"github.com/Dynom/TySug/finder"
	"github.com/xrash/smetrics"
)

func algJaroWinkler() finder.AlgWrapper {
	return func(a, b string) float64 {
		return smetrics.JaroWinkler(a, b, .7, 4)
	}
}

// New creates a new instance of Scorer
func New(list []string, options ...finder.Option) (*finder.Scorer, error) {
	defaults := []finder.Option{finder.OptSetAlgorithm(algJaroWinkler())}

	return finder.New(list, append(defaults, options...)...)
}
