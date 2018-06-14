package finder

import (
	"context"
	"errors"
	"math"
)

// Algorithm the type to comply with to create your own algorithm
type Algorithm func(a, b string) float64

// Finder is the type to find the nearest reference
type Finder struct {
	referenceMap map[string]struct{}
	reference    []string
	Alg          Algorithm
}

// Errors
var (
	ErrNoAlgorithmDefined = errors.New("no algorithm defined")
)

// New creates a new instance of Finder. The order of the list is significant
func New(list []string, options ...Option) (*Finder, error) {
	i := &Finder{
		referenceMap: make(map[string]struct{}, len(list)),
		reference:    list,
	}

	for _, r := range list {
		i.referenceMap[r] = struct{}{}
	}

	for _, o := range options {
		o(i)
	}

	if i.Alg == nil {
		return i, ErrNoAlgorithmDefined
	}

	return i, nil
}

// Find returns the best alternative and a score. A score of 1 means a perfect match
func (t Finder) Find(input string) (string, float64) {
	return t.FindCtx(context.Background(), input)
}

// FindCtx is the same as Find, with context support
func (t Finder) FindCtx(ctx context.Context, input string) (string, float64) {

	// Exact matches
	if _, exists := t.referenceMap[input]; exists {
		return input, 1
	}

	var hs = math.Inf(-1)
	var best string
	for _, ref := range t.reference {
		select {
		case <-ctx.Done():
			return input, 0
		default:
		}

		if d := t.Alg(input, ref); d > hs {
			hs = d
			best = ref
		}
	}

	return best, hs
}
