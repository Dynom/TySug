package finder

import (
	"testing"
	"time"
	"context"
)

func TestNewWithCustomAlgorithm(t *testing.T) {
	sug, _ := New([]string{"b"}, OptExampleAlgorithm)

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

func TestContextCancel(t *testing.T) {
	sug, err := New([]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m"}, func(sug *TySug) {
		sug.Alg = func(a, b string) float64 {
			time.Sleep(10 * time.Millisecond)
			return 1
		}
	})

	if err != nil {
		t.Errorf("Error when constructing TySug, %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Millisecond)
	defer cancel()

	timeStart := time.Now()
	sug.FindCtx("john", ctx)
	timeEnd := time.Now()

	timeSpent := int(timeEnd.Sub(timeStart).Seconds() * 100)

	if timeSpent != 1 {
		t.Errorf("Expected the context to cancel after one iteration")
	}
}

func BenchmarkSliceOrMap(b *testing.B) {
	size := 50
	var hashmap = make(map[int]int, size)
	var list = make([]int,size)

	for i := size - 1; i > 0; i-- {
		hashmap[i] = i
		list[i] = i
	}

	b.Run("Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = hashmap[i]
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