package service

import (
	"github.com/Dynom/TySug/finder"
	"github.com/xrash/smetrics"
)

func NewDomainService(list []string, options ...finder.Option) (Domain, error) {
	defaults := []finder.Option{finder.OptSetAlgorithm(algJaroWinkler())}

	scorer, err := finder.New(list, append(defaults, options...)...)
	if err != nil {
		return Domain{}, err
	}

	return Domain{
		scorer: scorer,
	}, nil
}

type Domain struct {
	scorer *finder.Scorer
}

func (ds Domain) Rank(input string) (string, float64) {
	return ds.scorer.Find(input)
}

func algJaroWinkler() finder.AlgWrapper {
	return func(a, b string) float64 {
		return smetrics.JaroWinkler(a, b, .7, 4)
	}
}
