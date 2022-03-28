package finder

import (
	"context"
	"errors"
	"math"
	"strings"
	"sync"
)

// Finder is the type to find the nearest reference
type Finder struct {
	referenceMap    referenceMapType
	reference       []string
	referenceBucket referenceBucketType
	Alg             Algorithm
	LengthTolerance float64 // A number between 0.0-1.0 (percentage) to allow for length miss-match, anything outside this is considered not similar. Set to 0 to disable.
	lock            sync.RWMutex
	bucketChars     uint // @todo figure out what (type of) bucket approach to take. Prefix or perhaps using an ngram/trie approach
}

// Errors
var (
	ErrNoAlgorithmDefined = errors.New("no algorithm defined")
)

type (
	referenceMapType    map[string]struct{}
	referenceBucketType map[rune][]string
)

// These constants hold the value of the lowest and highest possible scores. Compatible with JSON serialization.
// It's not ideal to mix presentation with business logic but in this instance it was convenient and similarly
// effective as math.Inf(-1)
const (
	WorstScoreValue = -1 * math.MaxFloat32
	BestScoreValue  = math.MaxFloat32
)

// New creates a new instance of Finder. The order of the list is significant
func New(list []string, options ...Option) (*Finder, error) {
	i := &Finder{}

	for _, o := range options {
		o(i)
	}

	i.Refresh(list)

	if i.Alg == nil {
		return i, ErrNoAlgorithmDefined
	}

	return i, nil
}

// Refresh replaces the internal reference list.
func (t *Finder) Refresh(list []string) {
	rm := make(referenceMapType, len(list))
	rb := make(referenceBucketType, 26)

	for _, r := range list {

		if r == "" {
			continue
		}

		rm[r] = struct{}{}

		// @todo make the bucket prefix length configurable
		if t.bucketChars > 0 {
			l := rune(r[0])
			if _, ok := rb[l]; !ok {
				rb[l] = make([]string, 0, 16)
			}
			rb[l] = append(rb[l], r)
		}
	}

	t.lock.Lock()
	t.reference = append(t.reference[0:0], list...)
	t.referenceMap = rm
	t.referenceBucket = rb
	t.lock.Unlock()
}

// Exact returns true if the input is an exact match.
func (t *Finder) Exact(input string) bool {
	t.lock.RLock()
	_, ok := t.referenceMap[input]
	t.lock.RUnlock()

	return ok
}

// Find returns the best alternative a score and if it was an exact match or not.
// Since algorithms can define their own upper-bound, there is no "best" value.
func (t *Finder) Find(input string) (string, float64, bool) {
	matches, score, exact := t.FindTopRankingCtx(context.Background(), input)
	return matches[0], score, exact
}

// FindCtx is the same as Find, with context support.
func (t *Finder) FindCtx(ctx context.Context, input string) (string, float64, bool) {
	matches, score, exact := t.FindTopRankingCtx(ctx, input)
	return matches[0], score, exact
}

// FindTopRankingCtx returns a list (of at least one element) of references with the same "best" score
func (t *Finder) FindTopRankingCtx(ctx context.Context, input string) ([]string, float64, bool) {
	r, s, e, _ := t.findTopRankingCtx(ctx, input, 0)
	return r, s, e
}

// FindTopRankingPrefixCtx requires the references to have an exact prefix match on N characters of the input.
// prefixLength cannot exceed length of input
func (t *Finder) FindTopRankingPrefixCtx(ctx context.Context, input string, prefixLength uint) (list []string, exact bool, err error) {
	list, _, exact, err = t.findTopRankingCtx(ctx, input, prefixLength)
	return
}

// getRefList returns the appropriate list of references. getRefList does not deal with locks!
func (t *Finder) getRefList(input string) []string {
	if len(input) > 0 {
		r := rune(input[0])
		if _, ok := t.referenceBucket[r]; ok {
			return t.referenceBucket[r]
		}
	}

	return t.reference
}

// GetMatchingPrefix returns up to max ref's, that start with the prefix argument
func (t *Finder) GetMatchingPrefix(ctx context.Context, prefix string, max uint) ([]string, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	var (
		list   = t.getRefList(prefix)
		result = make([]string, 0, max)
	)

	for _, ref := range list {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		if strings.HasPrefix(ref, prefix) {
			result = append(result, ref)
		}

		if max > 0 && max == uint(len(result)) {
			return result, nil
		}
	}

	return result, nil
}

func (t *Finder) findTopRankingCtx(ctx context.Context, input string, prefixLength uint) ([]string, float64, bool, error) {
	hs := WorstScoreValue

	if prefixLength > 0 && uint(len(input)) < prefixLength {
		return []string{input}, WorstScoreValue, false, errors.New("prefix length exceeds input length")
	}

	t.lock.RLock()
	defer t.lock.RUnlock()

	// Exact matches
	if _, exists := t.referenceMap[input]; exists || len(input) == 0 {
		return []string{input}, BestScoreValue, true, nil
	}

	var (
		list      = t.getRefList(input)
		sameScore = []string{input}
	)

	for _, ref := range list {
		select {
		case <-ctx.Done():
			return []string{input}, WorstScoreValue, false, ctx.Err()
		default:
		}

		if !meetsPrefixLengthMatch(prefixLength, input, ref) {
			continue
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

	return sameScore, hs, false, nil
}

// meetsPrefixLengthMatch tests is the strings both match until the specified length. A 0 length returns true
func meetsPrefixLengthMatch(length uint, input, reference string) bool {
	if length > 0 {
		if uint(len(reference)) < length {
			return false
		}

		if pi := length - 1; input[0:pi] != reference[0:pi] {
			return false
		}
	}

	return true
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
