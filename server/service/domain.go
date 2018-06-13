package service

import (
	"github.com/Dynom/TySug/finder"
	"github.com/sirupsen/logrus"
	"github.com/xrash/smetrics"
)

// NewDomain creates a new service
func NewDomain(references []string, logger *logrus.Logger, options ...finder.Option) (Service, error) {
	defaults := []finder.Option{finder.OptSetAlgorithm(algJaroWinkler())}

	scorer, err := finder.New(references, append(defaults, options...)...)
	if err != nil {
		return Service{}, err
	}

	return Service{
		scorer,
		logger,
	}, nil
}

// Service is the service type
type Service struct {
	scorer *finder.Finder
	logger *logrus.Logger
}

// Rank returns the nearest reference
func (s Service) Rank(input string) (string, float64) {
	suggestion, score := s.scorer.Find(input)
	s.logger.WithFields(logrus.Fields{
		"input":      input,
		"suggestion": suggestion,
		"score":      score,
	}).Debug("Completed new ranking request")

	return suggestion, score
}

func algJaroWinkler() finder.Algorithm {
	return func(a, b string) float64 {
		return smetrics.JaroWinkler(a, b, .7, 4)
	}
}
