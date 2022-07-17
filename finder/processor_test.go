package finder

import (
	"testing"
	"unicode/utf8"
)

func TestMarkingTypeProcessor(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "Hellow^H", want: "Hello"},
		{input: "Hellow^H Universe", want: "Hello Universe"},
		{input: "Hello Universe^WWorld", want: "Hello World"},
	}
	processor := MarkingTypoProcessor()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			if got := processor(tt.input); got != tt.want {
				t.Errorf("MarkingTypoProcessor(%q)\ngot : %q\nwant: %q", tt.input, got, tt.want)
			}
		})
	}
}

func Test_markingTypoReplacer(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "^H", want: ""},
		{input: "a^H", want: ""},
		{input: "a^Hb", want: "b"},
		{input: "ba^H", want: "b"},
		{input: "^Hb", want: "b"},
		{input: "ba^Hb", want: "bb"},

		{input: "^W", want: ""},
		{input: "a^W", want: ""},
		{input: "a^Wb", want: "b"},
		{input: "b a^W", want: "b "},
		{input: "^Wb", want: "b"},
		{input: "b a^Wb", want: "b b"},
		{input: "烊0^H^H0", want: "0"}, // Marking order precedece

		{input: "foo bE^HE^Ht^Wfoo", want: "foo foo"},

		{input: "longer-without-space^W", want: "longer-without-"},
		{input: "World Wide Mess^WWeb", want: "World Wide Web"},
		{input: "When the Hare eats the ^, or when the ^Eats^WThe^WHare", want: "When the Hare eats the ^, or when the ^Hare"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			if got := markingTypoReplacer(tt.input, 0); got != tt.want {
				t.Errorf("markingTypoReplacer(%q)\ngot : %q,\nwant: %q", tt.input, got, tt.want)
			}
		})
	}

	// Fuzzer found patterns
	regressions := []struct {
		name  string
		input string
		want  string
	}{
		{name: "single carrot", input: "^", want: "^"},
		{name: "Losing euros", input: "€\x80^H\u0080^H€^H", want: "€"},
		{name: "Starting unicode characters", input: "뽥^H", want: ""},
		{name: "Starting unicode character, multiple Marks", input: "烊0^H^H0", want: "0"},
		{name: "Removing an unicode character", input: "^뽥^H", want: "^"},
		{name: "Removing an unicode character", input: "뽥^H뽥^H", want: ""},
		{name: "Removing an unicode character", input: "鮉鮉^W", want: ""},
	}
	for _, tt := range regressions {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := markingTypoReplacer(tt.input, 0); got != tt.want {
				t.Errorf("markingTypoReplacer(%q)\ngot : %q,\nwant: %q", tt.input, got, tt.want)
			}
		})
	}
}

var r1 string

func BenchmarkMarkingTypoReplacer(b *testing.B) {
	tests := []string{
		"^H",
		"烊.^H",
		"烊Multibyte^W",
		"Nothing to see ha^Here, moving along!",
		"World Wide Mess^WWeb",
		"When the Hare eats the ^, or when the ^Eats^WThe^WHare",
	}

	for _, tt := range tests {
		tt := tt
		b.Run(tt, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r1 = markingTypoReplacer(tt, 0)
			}
		})
	}

	_ = r1
}

func FuzzMarkingTypoReplacer(f *testing.F) {
	tests := []string{
		"^",
		"\x80", // €
		"^H",
		"^뽥^H",
		"鮉鮉^W",
		"뽥^H뽥^H",
		"烊0^H^H0",
		"bogus!^H",
		"Nothing to see ha^Here, moving along!",
		"World Wide Mess^WWeb",
		"When the Hare eats the ^, or when the ^Eats^WThe^WHare",
	}

	for _, tt := range tests {
		f.Add(tt)
	}

	f.Fuzz(func(t *testing.T, input string) {
		result := markingTypoReplacer(input, 0)

		if !utf8.Valid([]byte(result)) && utf8.Valid([]byte(input)) {
			t.Errorf("Result from %q isn't valid utf-8: %q", input, result)
		}
	})
}
