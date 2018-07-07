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
	referenceMap    map[string]struct{}
	reference       []string
	Alg             Algorithm
	LengthTolerance float64 // A number between 0.0-1.0 (percentage) to allow for length miss-match, anything outside this is considered not similar. Set to 0 to disable.
}

// Errors
var (
	ErrNoAlgorithmDefined = errors.New("no algorithm defined")
)

// These constants hold the value of the lowest and highest possible scores
const (
	WorstScoreValue = -1 * math.MaxFloat32
	BestScoreValue  = math.MaxFloat32
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

// Find returns the best alternative a score and if it was an exact match or not.
// Since algorithms can define their own upper-bound, there is no "best" value.
func (t Finder) Find(input string) (string, float64, bool) {
	return t.FindCtx(context.Background(), input)
}

// FindCtx is the same as Find, with context support.
func (t Finder) FindCtx(ctx context.Context, input string) (string, float64, bool) {
	// Initial value, compatible with JSON serialisation. It's not ideal to mix presentation with business logic
	// but in this instance it was convenient and similarly effective to math.Inf(-1)
	var hs = WorstScoreValue

	// Exact matches
	if _, exists := t.referenceMap[input]; exists {
		return input, BestScoreValue, true
	}

	var best = input
	for _, ref := range t.reference {
		select {
		case <-ctx.Done():
			return input, WorstScoreValue, false
		default:
		}

		// Test if the input length is much less, making it an unlikely typo.
		if !meetsLengthTolerance(t.LengthTolerance, input, ref) {
			continue
		}

		if score := t.Alg(input, ref); score > hs {
			hs = score
			best = ref
		}
	}

	return best, hs, false
}

// meetsLengthTolerance checks if the input meets the length tolerance criteria
func meetsLengthTolerance(t float64, input, reference string) bool {
	if t == 0 {
		return true
	}

	inputLen := len(input)
	refLen := len(reference)
	threshold := int(math.Ceil(float64(inputLen) * t))

	// The result is N% of the length or at least 1 (due to math.Ceil)
	return refLen-threshold <= inputLen && inputLen <= refLen+threshold
}
