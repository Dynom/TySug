package TySug

import (
	"github.com/Dynom/TySug/finder"
	"github.com/xrash/smetrics"
)

func optSetDefaultAlgo() finder.Option {
	return func (sug *finder.TySug) {
		sug.Alg = func(a, b string) float64 {
			return smetrics.JaroWinkler(a, b, .7, 4)
		}
	}
}

// New creates a new instance of TySug
func New(list []string, options ...finder.Option) (*finder.TySug, error) {
	defaults := []finder.Option{optSetDefaultAlgo()}

	return finder.New(list, append(defaults, options...)...)
}
