package service

import (
	"github.com/Dynom/TySug/finder"
	"github.com/sirupsen/logrus"
	"github.com/xrash/smetrics"
)

// NewDomainService creates a new service
func NewDomainService(references []string, logger *logrus.Logger, options ...finder.Option) (Domain, error) {
	defaults := []finder.Option{finder.OptSetAlgorithm(algJaroWinkler())}

	scorer, err := finder.New(references, append(defaults, options...)...)
	if err != nil {
		return Domain{}, err
	}

	return Domain{
		scorer: scorer,
		logger: logger,
	}, nil
}

// Domain is the service type
type Domain struct {
	scorer *finder.Scorer
	logger *logrus.Logger
}

// Rank returns the nearest reference
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
