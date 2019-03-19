package finder

import (
	"context"
	"errors"
	"math"
)

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

// These constants hold the value of the lowest and highest possible scores. Compatible with JSON serialisation.
// It's not ideal to mix presentation with business logic but in this instance it was convenient and similarly
// effective as math.Inf(-1)
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
	matches, score, exact := t.FindTopRankingCtx(context.Background(), input)
	return matches[0], score, exact
}

// FindCtx is the same as Find, with context support.
func (t Finder) FindCtx(ctx context.Context, input string) (string, float64, bool) {
	matches, score, exact := t.FindTopRankingCtx(ctx, input)
	return matches[0], score, exact
}

// FindTopRankingCtx returns a list (of at least one element) of references with the same "best" score
func (t Finder) FindTopRankingCtx(ctx context.Context, input string) ([]string, float64, bool) {
	var hs = WorstScoreValue

	// Exact matches
	if _, exists := t.referenceMap[input]; exists {
		return []string{input}, BestScoreValue, true
	}

	var sameScore = []string{input}
	for _, ref := range t.reference {
		select {
		case <-ctx.Done():
			return []string{input}, WorstScoreValue, false
		default:
		}

		// Test if the input length differs too much from the reference, making it an unlikely typo.
		if !meetsLengthTolerance(t.LengthTolerance, input, ref) {
			continue
		}

		score := t.Alg(input, ref)
		if score > hs {
			hs = score
			sameScore = []string{ref}
		} else if score == hs {
			sameScore = append(sameScore, ref)
		}
	}

	return sameScore, hs, false
}

// meetsLengthTolerance checks if the input meets the length tolerance criteria. The percentage is based on `input`
func meetsLengthTolerance(t float64, input, reference string) bool {
	if t <= 0 {
		return true
	}

	if t > 1 {
		return false
	}

	inputLen := len(input)
	refLen := len(reference)
	threshold := int(math.Ceil(float64(inputLen) * t))

	// The result is N% of the length or at least 1 (due to math.Ceil)
	return refLen-threshold <= inputLen && inputLen <= refLen+threshold
}
