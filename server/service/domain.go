package service

import (
	"github.com/Dynom/TySug/finder"
	"github.com/sirupsen/logrus"
	"github.com/xrash/smetrics"
)

func NewDomainService(list []string, l *logrus.Logger, options ...finder.Option) (Domain, error) {
	defaults := []finder.Option{finder.OptSetAlgorithm(algJaroWinkler())}

	scorer, err := finder.New(list, append(defaults, options...)...)
	if err != nil {
		return Domain{}, err
	}

	return Domain{
		scorer: scorer,
		logger: l,
	}, nil
}

type Domain struct {
	scorer *finder.Scorer
	logger *logrus.Logger
}

func (ds Domain) Rank(input string) (string, float64) {
	suggestion, score := ds.scorer.Find(input)
	ds.logger.WithFields(logrus.Fields{
		"input":      input,
		"suggestion": suggestion,
		"score":      score,
	}).Debug("Completed new ranking request")

	return suggestion, score
}

func algJaroWinkler() finder.AlgWrapper {
	return func(a, b string) float64 {
		return smetrics.JaroWinkler(a, b, .7, 4)
	}
}
