package finder

import "testing"

func TestSetAlgorithm(t *testing.T) {
	veryPositiveAlg := func(a, b string) float64 {
		return 1
	}

	sug, err := New([]string{}, WithAlgorithm(veryPositiveAlg))

	if sug.Alg == nil || err == ErrNoAlgorithmDefined {
		t.Errorf("Expected the algorithm to be set")
	}
}
