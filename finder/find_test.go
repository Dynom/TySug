package finder

import (
	"context"
	"testing"
	"time"
)

func exampleAlgorithm(a, b string) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	if a[0] == b[0] {
		return 1
	}

	return 0
}

func TestOptExampleAlgorithm(t *testing.T) {
	alg := exampleAlgorithm

	if s := alg("", "apple juice"); s != 0 {
		t.Errorf("Expected the example algorithm to return 0 when an argument is empty.")
	}

	if s := alg("apple juice", ""); s != 0 {
		t.Errorf("Expected the example algorithm to return 0 when an argument is empty.")
	}

	if s := alg("apple", "juice"); s != 0 {
		t.Errorf("Expected the example algorithm to return 0 when the values don't match.")
	}

	if s := alg("tree", "trie"); s != 1 {
		t.Errorf("Expected the example algorithm to return 1 when the first letters match.")
	}
}

func TestNewWithCustomAlgorithm(t *testing.T) {
	sug, _ := New([]string{"b"}, WithAlgorithm(exampleAlgorithm))

	var score float64
	var exact bool

	_, score, exact = sug.Find("a")
	if exact {
		t.Errorf("Expected exact to be false, instead I got %t (the score is %f).", exact, score)
	}

	_, score, exact = sug.Find("b")
	if !exact {
		t.Errorf("Expected exact to be true, instead I got %t (the score is %f).", exact, score)
	}
}

func TestNoAlgorithm(t *testing.T) {
	_, err := New([]string{})

	if err != ErrNoAlgorithmDefined {
		t.Errorf("Expected an error to be returned when no algorithm was specified.")
	}
}

func TestNoInput(t *testing.T) {
	sug, _ := New([]string{}, WithAlgorithm(exampleAlgorithm))
	sug.Find("")
}

func TestContextCancel(t *testing.T) {
	sug, err := New([]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m"}, func(sug *Finder) {
		sug.Alg = func(a, b string) float64 {
			time.Sleep(10 * time.Millisecond)
			return 1
		}
	})

	if err != nil {
		t.Errorf("Error when constructing Finder, %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	timeStart := time.Now()
	sug.FindCtx(ctx, "john")
	timeEnd := time.Now()

	timeSpent := int(timeEnd.Sub(timeStart).Seconds() * 1000)

	if 50 > timeSpent || timeSpent >= 130 {
		t.Errorf("Expected the context to cancel after one iteration")
	}
}

func TestFind(t *testing.T) {
	refs := []string{
		"a", "b",
		"12", "23", "24", "25",
		"food", "foor", "fool", "foon",
		"bar", "baz", "ban", "bal",
	}

	mockAlg := func(a, b string) float64 {
		var left string
		var right string

		if len(a) > len(b) {
			left, right = a, b
		} else {
			right, left = a, b
		}

		return -1 * float64(len(left)-len(right))
	}

	f, _ := New(refs,
		WithAlgorithm(mockAlg),
		WithLengthTolerance(0),
	)

	f.Find("bat")
}

func TestMeetsLengthTolerance(t *testing.T) {
	testData := []struct {
		Expect    bool
		Input     string
		Reference string
		Tolerance float64
	}{
		{Expect: true, Input: "foo", Reference: "bar", Tolerance: -1},
		{Expect: true, Input: "foo", Reference: "bar", Tolerance: 0},
		{Expect: true, Input: "foo", Reference: "bar", Tolerance: 1},
		{Expect: false, Input: "foo", Reference: "bar", Tolerance: 2}, // erroneous situation

		{Expect: true, Input: "smooth", Reference: "smoothie", Tolerance: 0.2},
		{Expect: false, Input: "smooth", Reference: "smoothie", Tolerance: 0.1},

		{Expect: true, Input: "abc", Reference: "defghi", Tolerance: 0.9},
		{Expect: true, Input: "abc", Reference: "defg", Tolerance: 0.5},
	}

	for _, td := range testData {
		r := meetsLengthTolerance(td.Tolerance, td.Input, td.Reference)
		if r != td.Expect {
			t.Errorf("Expected the tolerance to be %t\n%+v", td.Expect, td)
		}
	}

}

func BenchmarkSliceOrMap(b *testing.B) {
	// With sets of more than 20 elements, maps become more efficient. (Not including setup costs)
	size := 20
	var hashMap = make(map[int]int, size)
	var list = make([]int, size)

	for i := size - 1; i > 0; i-- {
		hashMap[i] = i
		list[i] = i
	}

	b.Run("Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = hashMap[i]
		}
	})
	b.Run("List", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range list {
				_ = v
			}
		}
	})
}
