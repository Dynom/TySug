package TySug

import (
	"testing"
	"math"
)

const floatTolerance = 0.000001

func TestNew(t *testing.T) {
	expect := "example"
	sug, _ := New([]string{expect, "ample"})
	alt, _ := sug.Find("exampel")

	if alt != expect {
		t.Errorf("Expected '%s' to be '%s'.", alt, expect)
	}
}


func TestTestExactMatch(t *testing.T) {
	sug, _ := New([]string{"foo", "example", "dissipation"})
	input := "example"
	match, score := sug.Find(input)

	if match != input {
		t.Errorf("Expected the input '%s' to equal the best match '%s'", input, match)
	}

	if math.Abs(1-score) > floatTolerance {
		t.Errorf("Expected a score of ~1.0, instead it is: %f", score)
	}
}

func TestApproximateMatch(t *testing.T) {
	reference := "example"
	sug, _ := New([]string{"foo", reference, "dissipation"})
	match, _ := sug.Find("exampel")

	if match != reference {
		t.Errorf("Expected the input '%s' to equal the best match '%s'", reference, match)
	}
}

func BenchmarkBasicUsage(b *testing.B) {
	sug, _ := New([]string{"foo", "abr", "pulp"})

	b.Run("Direct match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = sug.Find("foo")
		}
	})

	b.Run("Non direct match, low score", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = sug.Find("popl")
		}
	})

	b.Run("Non direct match, high score", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = sug.Find("bar")
		}
	})

}
