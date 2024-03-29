package finder

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"
)

var inspirationalRefList = []string{
	"Waleed Abdalati",
	"Nerilie Abram",
	"Ernest Afiesimama",
	"Myles Allen",
	"Richard Alley",
	"Kevin Anderson",
	"James Annan",
	"Julie Arblaster",
	"David Archer",
	"Svante Arrhenius",
	"Sallie Baliunas",
	"Eric J. Barron",
	"Roger G. Barry",
	"Robin Bell",
	"Lennart Bengtsson",
	"André Berger",
	"Richard A. Betts",
	"Jacob Bjerknes",
	"Vilhelm Bjerknes",
	"Bert Bolin",
	"Gerard C. Bond",
	"Jason Box",
	"Raymond S. Bradley",
	"Keith Briffa",
	"Wallace Smith Broecker",
	"Harold E. Brooks",
	"Keith Browning",
	"Robert Cahalan",
	"Ken Caldeira",
	"Guy Stewart Callendar",
	"Mark Cane",
	"Anny Cazenave",
	"Robert D. Cess",
	"Jule G. Charney",
	"John Christy",
	"John A. Church",
	"Ralph J. Cicerone",
	"Danielle Claar",
	"Allison Crimmins",
	"Harmon Craig",
	"Paul J. Crutzen",
	"Heidi Cullen",
	"Balfour Currie",
	"Judith Curry",
	"Willi Dansgaard",
	"Scott Denning",
	"Andrew Dessler",
	"P. C. S. Devara",
	"Robert E. Dickinson",
	"Mark Dyurgerov",
	"Sylvia Earle",
	"Don Easterbrook",
	"Tamsin Edwards",
	"Arnt Eliassen",
	"Kerry Emanuel",
	"Matthew England",
	"Ian G. Enting",
	"Joe Farman",
	"Christopher Field",
	"Eunice Newton Foote",
	"Piers Forster",
	"Joseph Fourier",
	"Jennifer Francis",
	"Benjamin Franklin",
	"Chris Freeman",
	"Eigil Friis-Christensen",
	"Inez Fung",
	"Yevgraf Yevgrafovich Fyodorov",
	"Francis Galton",
	"Filippo Giorgi",
	"Peter Gleick",
	"Kenneth M. Golden",
	"Natalya Gomez",
	"Jonathan M. Gregory",
	"Jean Grove",
	"Joanna Haigh",
	"Edmund Halley",
	"Gordon Hamilton",
	"James E. Hansen",
	"Kenneth Hare",
	"Klaus Hasselmann",
	"Ed Hawkins",
	"Katharine Hayhoe",
	"Gabriele C. Hegerl",
	"Isaac Held",
	"Ann Henderson-Sellers",
	"Ellie Highwood",
	"David A. Hodell",
	"Ove Hoegh-Guldberg",
	"Greg Holland",
	"Brian Hoskins",
	"John T. Houghton",
	"Malcolm K. Hughes",
	"Mike Hulme",
	"Thomas Sterry Hunt",
	"Sherwood Idso",
	"Eystein Jansen",
	"Phil Jones",
	"Jean Jouzel",
	"Peter Kalmus",
	"Daniel Kammen",
	"Thomas R. Karl",
	"David Karoly",
	"Charles David Keeling",
	"Ralph Keeling",
	"David W. Keith",
	"Wilfrid George Kendrew",
	"Gretchen Keppel-Aleks",
	"Joseph B. Klemp",
	"Thomas Knutson",
	"Kirill Y. Kondratyev",
	"Bronwen Konecky",
	"Pancheti Koteswaram",
	"Shen Kuo",
	"John E. Kutzbach",
	"Dmitry Lachinov",
	"Hubert Lamb",
	"Kurt Lambeck",
	"Helmut Landsberg",
	"Christopher Landsea",
	"Mojib Latif",
	"Corinne Le Quéré",
	"Anders Levermann",
	"Richard Lindzen",
	"Diana Liverman",
	"Michael Lockwood",
	"Edward Norton Lorenz",
	"Claude Lorius",
	"James Lovelock",
	"Amanda Lynch",
	"Peter Lynch",
	"Michael MacCracken",
	"Gordon J. F. MacDonald",
	"Jerry D. Mahlman",
	"László Makra",
	"Syukuro Manabe",
	"Gordon Manley",
	"Michael E. Mann",
	"David Marshall",
	"Valerie Masson-Delmotte",
	"Gordon McBean",
	"James J. McCarthy",
	"Helen McGregor",
	"Christopher McKay",
	"Marcia McNutt",
	"Carl Mears",
	"Gerald A. Meehl",
	"Katrin Meissner",
	"Sebastian H. Mernild",
	"Patrick Michaels",
	"Milutin Milanković",
	"John F. B. Mitchell",
	"Fritz Möller",
	"Mario J. Molina",
	"Nils-Axel Mörner",
	"Richard H. Moss",
	"Richard A. Muller",
	"R. E. Munn FRSC",
	"Gerald North",
	"Hans Oeschger",
	"Atsumu Ohmura",
	"Cliff Ollier",
	"Abraham H. Oort",
	"Michael Oppenheimer",
	"Timothy Osborn",
	"Friederike Otto",
	"Tim Palmer CBE FRS",
	"Garth Paltridge",
	"David E. Parker",
	"Fyodor Panayev",
	"Graeme Pearman OA FAAS",
	"William Richard Peltier",
	"Jean Robert Petit",
	"David Phillips OC",
	"Roger A. Pielke",
	"Raymond Pierrehumbert",
	"Andrew Pitman",
	"Gilbert Plass",
	"Ian Plimer",
	"Henry Pollack",
	"Vicky Pope",
	"Detlef Quadfasel",
	"Stefan Rahmstorf",
	"Veerabhadran Ramanathan",
	"Michael Raupach",
	"Maureen Raymo",
	"David Reay",
	"Martine Rebetez",
	"Roger Revelle",
	"Lewis Fry Richardson",
	"Eric Rignot",
	"Alan Robock",
	"Joseph J. Romm",
	"Carl-Gustaf Rossby",
	"Frank Sherwood Rowland",
	"Cynthia E. Rosenzweig",
	"William Ruddiman",
	"Steve Running",
	"Murry Salby",
	"Jim Salinger",
	"Dork Sahagian",
	"Ben Santer",
	"Nicola Scafetta",
	"Hans Joachim Schellnhuber",
	"David Schindler",
	"Michael Schlesinger",
	"William H. Schlesinger",
	"Gavin A. Schmidt",
	"Stephen H. Schneider",
	"Daniel P. Schrag",
	"Stephen E. Schwartz",
	"Tom Segalstad",
	"Wolfgang Seiler",
	"John H. Seinfeld",
	"Mark Serreze",
	"Nicholas Shackleton",
	"Nir Shaviv",
	"J. Marshall Shepherd",
	"Drew Shindell",
	"Keith Shine",
	"Jagdish Shukla",
	"Joanne Simpson",
	"Fred Singer",
	"Julia Slingo",
	"Joseph Smagorinsky",
	"Susan Solomon",
	"Richard C. J. Somerville",
	"Willie Soon",
	"Kozma Spassky-Avtonomov",
	"Roy Spencer",
	"Konrad Steffen",
	"Will Steffen",
	"Thomas Stocker",
	"Hans von Storch",
	"Peter A. Stott",
	"Hans E. Suess",
	"Henrik Svensmark",
	"Kevin Russel Tate",
	"Simon Tett",
	"Peter Thejll",
	"Peter Thorne",
	"Liz Thomas",
	"Lonnie Thompson",
	"Axel Timmermann",
	"Micha Tomkiewicz",
	"Owen Toon",
	"Kevin E. Trenberth",
	"Susan Trumbore",
	"John Tyndall",
	"Jean-Pascal van Ypersele",
	"David Vaughan",
	"Jan Veizer",
	"Pier Vellinga",
	"Ricardo Villalba",
	"Peter Wadhams ScD",
	"Warren M. Washington",
	"John Michael Wallace",
	"Andrew Watson",
	"Sir Robert Watson",
	"Betsy Weatherhead",
	"Andrew J. Weaver",
	"Harry Wexler",
	"Penny Whetton",
	"Tom Wigley",
	"Josh Willis",
	"David Wratt",
	"Donald Wuebbles",
	"Carl Wunsch",
	"Olga Zolina",
	"Eduardo Zorita",
}

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
		sug.algorithm = func(a, b string) float64 {
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
				sug.algorithm = func(a, b string) float64 {
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
		want := make([]string, 0)
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
			t.Logf("list: %+v", list)
		}
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

func TestFinder_Exact(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{want: true, name: "exact match", input: "a"},

		{name: "empty input", input: ""},
		{name: "not exact input", input: "c"},
	}

	sug, err := New([]string{"a", "b", "z"}, WithAlgorithm(exampleAlgorithm))
	if err != nil {
		t.Errorf("Didn't expect construction to fail %v", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sug.Exact(tt.input); got != tt.want {
				t.Errorf("Exact() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_meetsPrefixLengthMatch(t *testing.T) {
	type args struct {
		length    uint
		input     string
		reference string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "0 skips the line, and will always return true",
			want: true,
			args: args{length: 0, input: "foo", reference: "foo"},
		},
		{
			name: "full length",
			want: true,
			args: args{length: 3, input: "foo", reference: "foo"},
		},
		{
			name: "input too short",
			args: args{length: 3, input: "fo", reference: "fee"},
		},
		{
			name: "ref too short",
			args: args{length: 3, input: "foo", reference: "fo"},
		},
		{
			name: "empty input and ref",
			args: args{length: 3, input: "", reference: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := meetsPrefixLengthMatch(tt.args.length, tt.args.input, tt.args.reference); got != tt.want {
				t.Errorf("meetsPrefixLengthMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkFindTopRankingCTXRace(b *testing.B) {
	sort.Strings(inspirationalRefList)
	f, err := New(
		inspirationalRefList[0:5],
		WithAlgorithm(exampleAlgorithm),
		WithLengthTolerance(0),
		WithPrefixBuckets(false),
	)
	if err != nil {
		b.Fatal("Setting up test failed")
	}

	ctx := context.Background()

	// Validating that we don't have race conditions
	b.Run("chaos", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.SetParallelism(10)

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _, _, _ = f.findTopRankingCtx(ctx, "a", 0)
				f.Refresh(inspirationalRefList)
			}
		})
	})
}

func BenchmarkFindTopRankingCTX(b *testing.B) {
	sort.Strings(inspirationalRefList)
	f, err := New(
		inspirationalRefList[0:5],
		WithAlgorithm(exampleAlgorithm),
		WithLengthTolerance(0),
		WithPrefixBuckets(false),
	)
	if err != nil {
		b.Fatal("Setting up test failed")
	}

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _, _ = f.findTopRankingCtx(ctx, "a", 0)
	}
}
