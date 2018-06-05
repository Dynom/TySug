package TySug

import (
	"testing"
	"math"
	"strings"
	"fmt"
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
	sug, _ := New([]string{"foo", "abr", "butterfly"})

	b.Run("Direct match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = sug.Find("foo")
		}
	})

	b.Run("Non direct match, low score", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = sug.Find("juice")
		}
	})

	b.Run("Non direct match, high score", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = sug.Find("butterfyl")
		}
	})
}

func ExampleNew() {
	input := "yourusername@hotmail.co"
	domains := []string{"gmail.com", "hotmail.com", "yahoo.com"}

	alt, _ := SuggestAlternative(input, domains)
	fmt.Printf("Perhaps you meant '%s' instead!", alt)
	// Output: Perhaps you meant 'example@hotmail.com' instead!
}

func SuggestAlternative(email string, domains []string) (string, float64) {

	i := strings.LastIndex(email, "@")
	if i <= 0 || i >= len(email) {
		return email, 0
	}

	// Extracting the local and domain parts
	localPart := email[:i]
	hostname := email[i+1:]

	sug, _ := New(domains)
	alternative, score := sug.Find(hostname)

	if score > 0.9 {
		combined := localPart + "@" + alternative
		return combined, score
	}

	return email, score
}
