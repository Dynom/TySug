package TySug

import (
	"fmt"
	"math"
	"strings"
	"testing"
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
	cases := []struct {
		Input  string
		Expect string
	}{
		{Input: "example", Expect: "example"},
		{Input: "CaseSensitive", Expect: "CaseSensitive"},
	}

	for _, td := range cases {
		sug, _ := New([]string{"foo", "example", "CaseSensitive", "cASEsENSITIVE"})
		match, score := sug.Find(td.Input)

		if match != td.Expect {
			t.Errorf("Expected the input '%s' to result in '%s', however the best match is '%s'", td.Input, td.Expect, match)
		}

		if math.Abs(1-score) > floatTolerance {
			t.Errorf("Expected a score of ~1.0, instead it is: %f", score)
		}
	}
}

func TestApproximateMatch(t *testing.T) {
	cases := []struct {
		Input     string
		Reference string
	}{
		{Input: "exampel", Reference: "example"},
		{Input: "casesensitive", Reference: "CaseSensitive"},
	}

	for _, td := range cases {
		sug, _ := New([]string{td.Reference})
		match, _ := sug.Find(td.Input)

		if match != td.Reference {
			t.Errorf("Expected the input '%s' to result in '%s', however the best match '%s'", td.Input, td.Reference, match)
		}
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
	domains := []string{"gmail.com", "hotmail.com", "yahoo.com", "example.com"}

	// Typo in the TLD
	input := "yourusername@example.co"

	alt, _ := SuggestAlternative(input, domains)
	fmt.Printf("Perhaps you meant '%s' instead!", alt)
	// Output: Perhaps you meant 'yourusername@example.com' instead!
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
	alternative, score := sug.Find(strings.ToLower(hostname))

	if score > 0.9 {
		combined := localPart + "@" + alternative
		return combined, score
	}

	return email, score
}
