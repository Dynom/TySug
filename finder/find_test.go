package finder

import (
	"context"
	"reflect"
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

func TestFinder_FindTopRankingPrefixCtx(t *testing.T) {
	refs := []string{
		"abcdef",
		"bcdef",
	}

	type args struct {
		input        string
		prefixLength uint
	}
	tests := []struct {
		name     string
		args     args
		wantList []string
		wantErr  bool
	}{

		// match
		{name: "prefix full size", args: args{input: "abcdef", prefixLength: 6}, wantList: refs[0:1]},
		{name: "prefix partial", args: args{input: "abcdef", prefixLength: 2}, wantList: refs[0:1]},

		// no match
		{name: "prefix miss-match", args: args{input: "monkey", prefixLength: 6}, wantList: []string{"monkey"}},

		// errors
		{wantErr: true, name: "len exceeds input", args: args{input: "abc", prefixLength: 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			finder, _ := New(refs, func(sug *Finder) {
				sug.Alg = func(a, b string) float64 {
					return 1
				}
			})

			ctx := context.Background()

			gotList, _, err := finder.FindTopRankingPrefixCtx(ctx, tt.args.input, tt.args.prefixLength)
			if (err != nil) != tt.wantErr {
				t1.Errorf("FindTopRankingPrefixCtx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// On failure the input is returned.
			if tt.wantErr {
				tt.wantList = []string{tt.args.input}
			}

			if !reflect.DeepEqual(gotList, tt.wantList) {
				t1.Errorf("FindTopRankingPrefixCtx() gotList = %v, want %v", gotList, tt.wantList)
			}
		})
	}
}

func TestFinder_RefreshWithBuckets(t *testing.T) {
	refs := []string{
		"aabb",
		"aabbcc",
		"aabbccdd",
		"aabcd",

		"bbcc",
		"bbccdd",
		"bbccddee",
		"bbcde",
	}

	finder, _ := New(
		refs,
		WithPrefixBuckets(true),
		WithAlgorithm(func(a, b string) float64 {
			return BestScoreValue
		}),
	)

	t.Run("bucket size", func(t1 *testing.T) {
		if bl := len(finder.referenceBucket); bl != 2 {
			t.Errorf("Expecting two buckets got: %d, want: %d", bl, 2)

			for chr := range finder.referenceBucket {
				t.Logf("Bucket chars: %c", chr)
			}

			return
		}
	})

	t.Run("testing bucket contents", func(t *testing.T) {
		if finder.bucketChars != 1 {
			t.Errorf("Expecting only single rune buckets, instead it's %d", finder.bucketChars)
			return
		}

		const want = "bbccddee"
		list := finder.referenceBucket[rune(want[0])]
		var match bool
		for _, v := range list {
			if v == want {
				match = true
				break
			}
		}

		if !match {
			t.Errorf("Expected to find %q in the reference bucket", want)
		}
	})

	t.Run("testing bucket similarity", func(t *testing.T) {
		if finder.bucketChars != 1 {
			t.Errorf("Expecting only single rune buckets, instead it's %d", finder.bucketChars)
			return
		}

		input := "beer"
		bucketRune := rune(input[0])

		// making the test a bit more robust
		var want = make([]string, 0)
		for _, v := range refs {
			if rune(v[0]) == bucketRune {
				want = append(want, v)
			}
		}

		// due to the very liberal "algorithm", anything matches, as long as the bucket prefix is respected
		got, _, _, _ := finder.findTopRankingCtx(context.Background(), input, 0)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected the reference bucket to be %+v, instead got: %+v", want, got)
		}
	})
}

func TestFinder_GetMatchingPrefix(t *testing.T) {
	refs := []string{
		"ada lovelace",
		"grace hopper",
		"ida rhodes",
		"sophie wilson",
		"aminata sana congo",
		"mary lou jepsen",
		"shafi goldwasser",
	}

	type args struct {
		prefix string
		max    uint
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "single", args: args{prefix: "a", max: 1}, want: []string{"ada lovelace"}},
		{name: "multiple", args: args{prefix: "a", max: 2}, want: []string{"ada lovelace", "aminata sana congo"}},
		{name: "no max == all", args: args{prefix: "a", max: 0}, want: []string{"ada lovelace", "aminata sana congo"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sug, _ := New(refs, WithAlgorithm(exampleAlgorithm))
			got, err := sug.GetMatchingPrefix(context.Background(), tt.args.prefix, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMatchingPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMatchingPrefix() got = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("Testing context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		sug, _ := New(refs, WithPrefixBuckets(true), WithAlgorithm(exampleAlgorithm))
		list, err := sug.GetMatchingPrefix(ctx, "a", 2)

		if err == nil {
			t.Errorf("GetMatchingPrefix() error = %v", err)
		}

		t.Logf("list: %+v", list)
	})
}

func TestFinder_getRefList(t *testing.T) {

	refs := []string{
		"balloon",
		"basketball",
		"sea lion",
		"celebration",
		"sunshine",
	}

	tests := []struct {
		name  string
		input string
		want  uint
	}{
		// input exists in ref list. With buckets enabled we should see 2 results starting with the same letter
		{name: "Selecting bucket ref list", input: "balloon", want: 2},

		// no match, the entire list should be returned
		{name: "Selecting full ref list on no match", input: "lion", want: 5},

		// no input, the entire list should be returned
		{name: "Selecting full ref list on empty input", input: "", want: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sug, err := New(refs, WithPrefixBuckets(true), WithAlgorithm(exampleAlgorithm))
			if err != nil {
				t.Errorf("b00m headshot %+v", err)
			}

			if got := sug.getRefList(tt.input); uint(len(got)) != tt.want {
				t.Errorf("getRefList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFinder_Refresh(t *testing.T) {
	tests := []struct {
		name    string
		refs    []string
		buckets uint
	}{
		// TODO: Add test cases.

		{name: "refs without empties", refs: []string{"a", "b"}, buckets: 2},
		{name: "refs with empties", refs: []string{"", "a"}, buckets: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sug, err := New([]string{}, WithPrefixBuckets(true), WithAlgorithm(exampleAlgorithm))

			if err != nil {
				t.Errorf("Didn't expect construction to fail %v", err)
				return
			}

			sug.Refresh(tt.refs)
			if got := uint(len(sug.referenceBucket)); got != tt.buckets {
				t.Errorf("Expected %d buckets, instead I got %d", tt.buckets, got)

				t.Logf("Reference Map: %+v", sug.referenceMap)
				t.Logf("Reference: %+v", sug.reference)
				t.Logf("Reference Bucket: %+v", sug.referenceBucket)
			}
		})
	}
}
