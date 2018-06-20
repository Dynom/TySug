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

func TestSetTolerance(t *testing.T) {
	veryPositiveAlg := func(a, b string) float64 {
		return 1
	}

	testData := []struct {
		Input  string
		Expect string
		Exact  bool
	}{
		// Too long
		{Input: "coffeeeee", Expect: "coffeeeee", Exact: false},

		// Within the threshold of "2" characters
		{Input: "coffeeee", Expect: "coffee", Exact: false},
		{Input: "coffeee", Expect: "coffee", Exact: false},
		{Input: "coffee", Expect: "coffee", Exact: true},
		{Input: "coffe", Expect: "coffee", Exact: false},

		// Too short
		{Input: "coff", Expect: "coff", Exact: false},
		{Input: "cof", Expect: "cof", Exact: false},
	}

	for _, td := range testData {
		sug, _ := New([]string{"coffee"}, WithAlgorithm(veryPositiveAlg), WithLengthTolerance(0.2))

		m, _, e := sug.Find(td.Input)
		if m != td.Expect {
			t.Errorf("Expected the best match to equal '%s', instead I got '%s'.", td.Expect, m)
		}

		if e != td.Exact {
			t.Errorf("Expected the result to match %t, instead I got %t.", td.Exact, e)
		}
	}
}

func TestSetToleranceDisable(t *testing.T) {
	veryPositiveAlg := func(a, b string) float64 {
		return 1
	}

	expect := "coffee"
	sug, _ := New([]string{expect}, WithAlgorithm(veryPositiveAlg), WithLengthTolerance(0))

	// The input is a string exceeding the threshold.
	m, _, _ := sug.Find("coffee" + "eeeeeeeee")
	if m != expect {
		t.Errorf("Expected the best match to equal '%s', instead I got '%s'.", expect, m)
	}
}
