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
	sug, _ := New([]string{"b"}, OptSetAlgorithm(exampleAlgorithm))

	var score float64

	_, score = sug.Find("a")
	if score != 0 {
		t.Errorf("Expected the score to be 0, instead I got %f.", score)
	}

	_, score = sug.Find("b")
	if score != 1 {
		t.Errorf("Expected the score to be 1, instead I got %f.", score)
	}
}

func TestNoAlgorithm(t *testing.T) {
	_, err := New([]string{})

	if err != ErrNoAlgorithmDefined {
		t.Errorf("Expected an error to be returned when no algorithm was specified.")
	}
}

func TestNoInput(t *testing.T) {
	sug, _ := New([]string{}, OptSetAlgorithm(exampleAlgorithm))
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancel()

	timeStart := time.Now()
	sug.FindCtx(ctx, "john")
	timeEnd := time.Now()

	timeSpent := int(timeEnd.Sub(timeStart).Seconds() * 100)

	if timeSpent != 1 {
		t.Errorf("Expected the context to cancel after one iteration")
	}
}

func BenchmarkSliceOrMap(b *testing.B) {
	size := 50
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
