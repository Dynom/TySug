package finder

import (
	"flag"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/xrash/smetrics"
)

var (
	flagWithJaroReferences       = flag.Bool("with-jaro-reference", false, "Include more verbose information around Jaro reference tests")
	flagFailOnReferenceMissMatch = flag.Bool("fail-on-reference", false, "Fail tests when they don't strictly match their reference implementation")
)

// Reference list from the official publication https://www.census.gov/srd/papers/pdf/rrs2006-02.pdf
var jaroReferenceList = []struct {
	a     string
	b     string
	score float64 // The Jaro (not the Jaro-Winkler) score
}{
	{a: "SHACKLEFORD", b: "SHACKELFORD", score: 0.970},
	{a: "DUNNINGHAM", b: "CUNNIGHAM", score: 0.896},
	{a: "NICHLESON", b: "NICHULSON", score: 0.926},
	{a: "JONES", b: "JOHNSON", score: 0.790},
	{a: "MASSEY", b: "MASSIE", score: 0.889},
	{a: "ABROMS", b: "ABRAMS", score: 0.889},
	{a: "HARDIN", b: "MARTINEZ", score: 0.000},
	{a: "ITMAN", b: "SMITH", score: 0.000},
	{a: "JERALDINE", b: "GERALDINE", score: 0.926},
	{a: "MARHTA", b: "MARTHA", score: 0.944},
	{a: "MICHELLE", b: "MICHAEL", score: 0.869},
	{a: "JULIES", b: "JULIUS", score: 0.889},
	{a: "TANYA", b: "TONYA", score: 0.867},
	{a: "DWAYNE", b: "DUANE", score: 0.822},
	{a: "SEAN", b: "SUSAN", score: 0.783},
	{a: "JOHN", b: "JON", score: 0.000},
	{a: "JON", b: "JOHN", score: 0.917},
	{a: "JON", b: "JAN", score: 0.000},
}

var homoPhones = [][]string{
	{"Ad", "Add"},
	{"Aerie", "Airy"},
	{"Ail", "Ale"},
	{"Air", "Heir"},
	{"Aisle", "Isle"},
	{"All", "Awl"},
	{"Allowed", "Aloud"},
	{"Altar", "Alter"},
	{"Arc", "Ark"},
	{"Ascent", "Assent"},
	{"Ate", "Eight"},
	{"Attendance", "Attendants"},
	{"Aural", "Oral"},
	{"Axes", "Axis"},
	{"Aye", "Eye"},
	{"Bail", "Bale"},
	{"Baited", "Bated"},
	{"Bald", "Bawled"},
	{"Ball", "Bawl"},
	{"Band", "Banned"},
	{"Bard", "Barred"},
	{"Bare", "Bear"},
	{"Baron", "Barren"},
	{"Base", "Bass"},
	{"Based", "Baste"},
	{"Be", "Bee"},
	{"Beach", "Beech"},
	{"Beat", "Beet"},
	{"Beau", "Bow"},
	{"Beer", "Bier"},
	{"Bell", "Belle"},
	{"Berry", "Bury"},
	{"Berth", "Birth"},
	{"Billed", "Build"},
	{"Blew", "Blue"},
	{"Bloc", "Block"},
	{"Boar", "Bore"},
	{"Board", "Bored"},
	{"Boarder", "Border"},
	{"Bold", "Bowled"},
	{"Bolder", "Boulder"},
	{"Bootie", "Booty"},
	{"Born", "Borne"},
	{"Burro", "Burrow"},
	{"Bough", "Bow"},
	{"Braid", "Brayed"},
	{"Brake", "Break"},
	{"Bread", "Bred"},
	{"Brewed", "Brood"},
	{"Brews", "Bruise"},
	{"Bridal", "Bridle"},
	{"Broach", "Brooch"},
	{"Brows", "Browse"},
	{"But", "Butt"},
	{"Buy", "Bye"},
	{"Caddie", "Caddy"},
	{"Callous", "Callus"},
	{"Canon", "Cannon"},
	{"Canter", "Cantor"},
	{"Canvas", "Canvass"},
	{"Capital", "Capitol"},
	{"Carat", "Carrot"},
	{"Carol", "Carrel"},
	{"Cast", "Caste"},
	{"Cede", "Seed"},
	{"Ceiling", "Sealing"},
	{"Cell", "Sell"},
	{"Cellar", "Seller"},
	{"Censor", "Sensor"},
	{"Cent", "Scent"},
	{"Cereal", "Serial"},
	{"Cession", "Session"},
	{"Chance", "Chants"},
	{"Chased", "Chaste"},
	{"Cheap", "Cheep"},
	{"Chews", "Choose"},
	{"Chic", "Sheik"},
	{"Chilli", "Chilly"},
	{"Choir", "Quire"},
	{"Choral", "Coral"},
	{"Chord", "Cored"},
	{"Chute", "Shoot"},
	{"Cite", "Sight"},
	{"Clause", "Claws"},
	{"Close", "Clothes"},
	{"Coarse", "Course"},
	{"Colonel", "Kernel"},
	{"Complement", "Compliment"},
	{"Council", "Counsel"},
	{"Coward", "Cowered"},
	{"Creak", "Creek"},
	{"Crewel", "Cruel"},
	{"Crews", "Cruise"},
	{"Cue", "Queue"},
	{"Currant", "Current"},
	{"Cygnet", "Signet"},
	{"Cymbal", "Symbol"},
	{"Dam", "Damn"},
	{"Days", "Daze"},
	{"Dear", "Deer"},
	{"Dense", "Dents"},
	{"Desert", "Dessert"},
	{"Dew", "Due"},
	{"Die", "Dye"},
	{"Disburse", "Disperse"},
	{"Discreet", "Discrete"},
	{"Doe", "Dough"},
	{"Does", "Doze"},
	{"Done", "Dun"},
	{"Dual", "Duel"},
	{"Dyeing", "Dying"},
	{"Earn", "Urn"},
	{"Eave", "Eve"},
	{"Eek", "Eke"},
	{"Eight", "Ate"},
	{"Ewe", "You"},
	{"Ewes", "Yews"},
	{"Eye", "Aye"},
	{"Eyelet", "Islet"},
	{"Faint", "Feint"},
	{"Fair", "Fare"},
	{"Faun", "Fawn"},
	{"Faze", "Phase"},
	{"Feat", "Feet"},
	{"Fined", "Find"},
	{"Fir", "Fur"},
	{"Fisher", "Fissure"},
	{"Flair", "Flare"},
	{"Flea", "Flee"},
	{"Flew", "Flue"},
	{"Floe", "Flow"},
	{"Flour", "Flower"},
	{"Foaled", "Fold"},
	{"For", "Four"},
	{"Foreword", "Forward"},
	{"Forth", "Fourth"},
	{"Foul", "Fowl"},
	{"Franc", "Frank"},
	{"Frays", "Phrase"},
	{"Freeze", "Frieze"},
	{"Friar", "Fryer"},
	{"Gaff", "Gaffe"},
	{"Gait", "Gate"},
	{"Gamble", "Gambol"},
	{"Genes", "Jeans"},
	{"Gibe", "Jibe"},
	{"Gild", "Guild"},
	{"Gilt", "Guilt"},
	{"Knew", "New"},
	{"Gofer", "Gopher"},
	{"Gorilla", "Guerrilla"},
	{"Gourd", "Gored"},
	{"Grate", "Great"},
	{"Grill", "Grille"},
	{"Grisly", "Grizzly"},
	{"Groan", "Grown"},
	{"Guessed", "Guest"},
	{"Guise", "Guys"},
	{"Hail", "Hale"},
	{"Hair", "Hare"},
	{"Hall", "Haul"},
	{"Hangar", "Hanger"},
	{"Hart", "Heart"},
	{"Hay", "Hey"},
	{"Heal", "Heel"},
	{"Hear", "Here"},
	{"Heard", "Herd"},
	{"Heir", "Air"},
	{"Heroin", "Heroine"},
	{"Hew", "Hue"},
	{"Hi", "High"},
	{"Him", "Hymn"},
	{"Ho", "Hoe"},
	{"Hoard", "Horde"},
	{"Hoarse", "Horse"},
	{"Hoes", "Hose"},
	{"Hole", "Whole"},
	{"Holy", "Wholly"},
	{"Hostel", "Hostile"},
	{"Hour", "Our"},
	{"Idle", "Idol"},
	{"In", "Inn"},
	{"Incidence", "Incidents"},
	{"Intense", "Intents"},
	{"Isle", "Aisle"},
	{"Islet", "Eyelet"},
	{"Jam", "Jamb"},
	{"Jeans", "Genes"},
	{"Jibe", "Gibe"},
	{"Kernel", "Colonel"},
	{"Knave", "Nave"},
	{"Knead", "Need"},
	{"Knew", "New"},
	{"Knight", "Night"},
	{"Knit", "Nit"},
	{"Knot", "Not"},
	{"Know", "No"},
	{"Knows", "Nose"},
	{"Lacks", "Lax"},
	{"Lain", "Lane"},
	{"Lama", "Llama"},
	{"Laps", "Lapse"},
	{"Lay", "Lei"},
	{"Leach", "Leech"},
	{"Lead", "Led"},
	{"Leak", "Leek"},
	{"Lean", "Lien"},
	{"Leased", "Least"},
	{"Lessen", "Lesson"},
	{"Levee", "Levy"},
	{"Lie", "Lye"},
	{"Links", "Lynx"},
	{"Lo", "Low"},
	{"Load", "Lode"},
	{"Loan", "Lone"},
	{"Locks", "Lox"},
	{"Loot", "Lute"},
	{"Made", "Maid"},
	{"Mail", "Male"},
	{"Main", "Mane"},
	{"Maize", "Maze"},
	{"Mall", "Maul"},
	{"Manner", "Manor"},
	{"Mantel", "Mantle"},
	{"Marshal", "Martial"},
	{"Mask", "Masque"},
	{"Massed", "Mast"},
	{"Meat", "Meet"},
	{"Medal", "Meddle"},
	{"Metal", "Mettle"},
	{"Mewl", "Mule"},
	{"Mews", "Muse"},
	{"Might", "Mite"},
	{"Mince", "Mints"},
	{"Mind", "Mined"},
	{"Miner", "Minor"},
	{"Missal", "Missile"},
	{"Missed", "Mist"},
	{"Moan", "Mown"},
	{"Moose", "Mousse"},
	{"Morning", "Mourning"},
	{"Muscle", "Mussel"},
	{"Mustard", "Mustered"},
	{"Naval", "Navel"},
	{"Nave", "Knave"},
	{"Nay", "Neigh"},
	{"Need", "Knead"},
	{"New", "Gnu"},
	{"Nicks", "Nix"},
	{"Night", "Knight"},
	{"Nit", "Knit"},
	{"No", "Know"},
	{"None", "Nun"},
	{"Nose", "Knows"},
	{"Not", "Knot"},
	{"Oar", "Ore"},
	{"Ode", "Owed"},
	{"Oh", "Owe"},
	{"One", "Won"},
	{"Oral", "Aural"},
	{"Our", "Hour"},
	{"Paced", "Paste"},
	{"Packed", "Pact"},
	{"Pail", "Pale"},
	{"Pain", "Pane"},
	{"Pair", "Pear"},
	{"Palette", "Pallet"},
	{"Passed", "Past"},
	{"Patience", "Patients"},
	{"Pause", "Paws"},
	{"Peace", "Piece"},
	{"Peak", "Peek"},
	{"Peal", "Peel"},
	{"Pearl", "Purl"},
	{"Pedal", "Peddle"},
	{"Peer", "Pier"},
	{"Phase", "Faze"},
	{"Phrase", "Frays"},
	{"Plain", "Plane"},
	{"Plait", "Plate"},
	{"Pleas", "Please"},
	{"Plum", "Plumb"},
	{"Pole", "Poll"},
	{"Pore", "Pour"},
	{"Praise", "Prays"},
	{"Presence", "Presents"},
	{"Pride", "Pried"},
	{"Pries", "Prize"},
	{"Primer", "Primmer"},
	{"Prince", "Prints"},
	{"Principal", "Principle"},
	{"Profit", "Prophet"},
	{"Quarts", "Quartz"},
	{"Queue", "Cue"},
	{"Quire", "Choir"},
	{"Rabbet", "Rabbit"},
	{"Rain", "Reign"},
	{"Raise", "Rays"},
	{"Rap", "Wrap"},
	{"Rapt", "Wrapped"},
	{"Real", "Reel"},
	{"Red", "Read"},
	{"Read", "Reed"},
	{"Reek", "Wreak"},
	{"Residence", "Residents"},
	{"Rest", "Wrest"},
	{"Retch", "Wretch"},
	{"Review", "Revue"},
	{"Right", "Write"},
	{"Ring", "Wring"},
	{"Road", "Rowed"},
	{"Roe", "Row"},
	{"Role", "Roll"},
	{"Roomer", "Rumour"},
	{"Root", "Route"},
	{"Rose", "Rows"},
	{"Rote", "Wrote"},
	{"Rouse", "Rows"},
	{"Rude", "Rued"},
	{"Rung", "Wrung"},
	{"Rye", "Wry"},
	{"Sac", "Sack"},
	{"Sail", "Sale"},
	{"Sane", "Seine"},
	{"Saver", "Savour"},
	{"Scene", "Seen"},
	{"Scent", "Sent"},
	{"Scull", "Skull"},
	{"Sea", "See"},
	{"Sealing", "Ceiling"},
	{"Seam", "Seem"},
	{"Seas", "Seize"},
	{"Seed", "Cede"},
	{"Sell", "Cell"},
	{"Seller", "Cellar"},
	{"Sensor", "Censor"},
	{"Serf", "Surf"},
	{"Serge", "Surge"},
	{"Serial", "Cereal"},
	{"Session", "Cession"},
	{"Sew", "Sow"},
	{"Shearn", "Sheer"},
	{"Sheik", "Chic"},
	{"Shoe", "Shoo"},
	{"Shone", "Shown"},
	{"Shoot", "Chute"},
	{"Sic", "Sick"},
	{"Side", "Sighed"},
	{"Sighs", "Size"},
	{"Sight", "Cite"},
	{"Signet", "Cygnet"},
	{"Slay", "Sleigh"},
	{"Sleight", "Slight"},
	{"Soar", "Sore"},
	{"Soared", "Sword"},
	{"Sole", "Soul"},
	{"Soled", "Sold"},
	{"Some", "Sum"},
	{"Son", "Sun"},
	{"Staid", "Stayed"},
	{"Stairs", "Stares"},
	{"Stake", "Steak"},
	{"Stationary", "Stationery"},
	{"Steal", "Steel"},
	{"Step", "Steppe"},
	{"Stile", "Style"},
	{"Straight", "Strait"},
	{"Succour", "Sucker"},
	{"Suite", "Sweet"},
	{"Symbol", "Cymbal"},
	{"Tacked", "Tact"},
	{"Tacks", "Tax"},
	{"Tail", "Tale"},
	{"Taper", "Tapir"},
	{"Taught", "Taut"},
	{"Tea", "Tee"},
	{"Team", "Teem"},
	{"Tear", "Tier"},
	{"Teas", "Tease"},
	{"Tense", "Tents"},
	{"Tern", "Turn"},
	{"Their", "There"},
	{"Threw", "Through"},
	{"Throes", "Throws"},
	{"Throne", "Thrown"},
	{"Thyme", "Time"},
	{"Tic", "Tick"},
	{"Tide", "Tied"},
	{"Too", "Two"},
	{"Toad", "Towed"},
	{"Toe", "Tow"},
	{"Told", "Tolled"},
	{"Tracked", "Tract"},
	{"Troop", "Troupe"},
	{"Trussed", "Trust"},
	{"Undo", "Undue"},
	{"Urn", "Earn"},
	{"Use", "Yews"},
	{"Vain", "Vein"},
	{"Vale", "Veil"},
	{"Vice", "Vise"},
	{"Wade", "Weighed"},
	{"Waist", "Waste"},
	{"Wait", "Weight"},
	{"Waive", "Wave"},
	{"Wares", "Wears"},
	{"Warn", "Worn"},
	{"Way", "Weigh"},
	{"We", "Wee"},
	{"Weak", "Week"},
	{"Whole", "Hole"},
	{"Wholly", "Holy"},
	{"Won", "One"},
	{"Wood", "Would"},
	{"Wrap", "Rap"},
	{"Wrapped", "Rapped"},
	{"Wreak", "Reek"},
	{"Wrest", "Rest"},
	{"Wretch", "Retch"},
	{"Wring", "Ring"},
	{"Write", "Right"},
	{"Wrote", "Rote"},
	{"Wrung", "Rung"},
	{"Wry", "Rye"},
	{"Yew", "You"},
	{"Yews", "Use"},
	{"Yoke", "Yolk"},
}

func equal(a, b float64) bool {
	const radix = 0.0005

	if a > b {
		return a-b < radix
	}

	return b-a < radix
}

func TestHomoPhoneJaroImplementations(t *testing.T) {
	if !*flagWithJaroReferences {
		t.Skip("Ignoring Jaro test, enable with -with-jaro-reference")
		return
	}

	const (
		smetricsJaro         = iota
		rosettaJaroV0        = iota
		rosettaJaroV1        = iota
		jaroDistanceMasatana = iota
	)

	for _, tt := range homoPhones {
		a := strings.ToLower(tt[0])
		b := strings.ToLower(tt[1])

		scores := make([]float64, 4)
		scores[smetricsJaro] = smetrics.Jaro(a, b)
		scores[rosettaJaroV0] = RosettaJaroV0(a, b)
		scores[rosettaJaroV1] = NewJaro()(a, b)
		scores[jaroDistanceMasatana] = func() float64 {
			s, _ := JaroDistanceMasatana(a, b)
			return s
		}()

		tmp := scores[smetricsJaro]
		for i, score := range scores {
			if !equal(score, tmp) {
				t.Errorf("%f != %f testing i %d (a: %q, b: %q)", tmp, score, i, a, b)

				t.Logf("input: %q/%q", a, b)
				t.Logf("smetrics.Jaro        %f", scores[smetricsJaro])
				t.Logf("RosettaJaroV0        %f", scores[rosettaJaroV0])
				t.Logf("RosettaJaroV1        %f", scores[rosettaJaroV1])
				t.Logf("JaroDistanceMasatana %f", scores[jaroDistanceMasatana])
			}

			tmp = score
		}
	}
}

func TestJaroImplementations(t *testing.T) {
	RosettaJaroV1 := NewJaro()

	t.Run("reference list", func(t *testing.T) {
		for _, tt := range jaroReferenceList {
			score := NewJaro()(tt.a, tt.b)

			// Skip the test if the reference has a score of 0
			if equal(tt.score, 0) && !*flagFailOnReferenceMissMatch {
				continue
			}

			if !equal(tt.score, score) {
				t.Errorf("Expected a score of %f, instead it was %f for input, a: %q, b: %q ", tt.score, score, tt.a, tt.b)

				t.Logf("%q vs. %q", tt.a, tt.b)
				t.Logf("smetrics.Jaro        %f", smetrics.Jaro(tt.a, tt.b))
				t.Logf("RosettaJaroV0        %f", RosettaJaroV0(tt.a, tt.b))
				t.Logf("RosettaJaroV1        %f", RosettaJaroV1(tt.a, tt.b))
				t.Logf("JaroDistanceMasatana %f", func() float64 {
					s, _ := JaroDistanceMasatana(tt.a, tt.b)
					return s
				}())
			}
		}
	})

	t.Run("specific variants", func(t *testing.T) {
		variants := []struct {
			a     string
			b     string
			score float64
		}{
			{a: "a", b: "b", score: 0},
			{a: "x", b: "x", score: 1},
		}

		for _, tt := range variants {
			score := NewJaro()(tt.a, tt.b)

			// Skip the test if the reference has a score of 0
			if equal(tt.score, 0) && !*flagFailOnReferenceMissMatch {
				continue
			}

			if !equal(tt.score, score) {
				t.Errorf("Expected a score of %f, instead it was %f for input, a: %q, b: %q ", tt.score, score, tt.a, tt.b)

				t.Logf("%q vs. %q", tt.a, tt.b)
				t.Logf("smetrics.Jaro        %f", smetrics.Jaro(tt.a, tt.b))
				t.Logf("RosettaJaroV0        %f", RosettaJaroV0(tt.a, tt.b))
				t.Logf("RosettaJaroV1        %f", RosettaJaroV1(tt.a, tt.b))
				t.Logf("JaroDistanceMasatana %f", func() float64 {
					s, _ := JaroDistanceMasatana(tt.a, tt.b)
					return s
				}())
			}
		}
	})
}

// From: github.com/masatana/go-textdistance
func JaroDistanceMasatana(s1, s2 string) (float64, int) {
	if s1 == s2 {
		return 1.0, 0.0
	}
	// compare length using rune slice length, as s1 and s2 are not necessarily ASCII-only strings
	longer, shorter := []rune(s1), []rune(s2)
	if len(longer) < len(shorter) {
		longer, shorter = shorter, longer
	}
	scope := int(math.Floor(float64(len(longer)/2))) - 1
	// m is the number of matching characters.
	m := 0
	matchFlags := make([]bool, len(longer))
	matchIndexes := make([]int, len(longer))
	for i := range matchIndexes {
		matchIndexes[i] = -1
	}

	for i := 0; i < len(shorter); i++ {
		k := Min(i+scope+1, len(longer))
		for j := Max(i-scope, 0); j < k; j++ {
			if matchFlags[j] || shorter[i] != longer[j] {
				continue
			}
			matchIndexes[i] = j
			matchFlags[j] = true
			m++
			break
		}
	}
	ms1 := make([]rune, m)
	ms2 := make([]rune, m)
	si := 0
	for i := 0; i < len(shorter); i++ {
		if matchIndexes[i] != -1 {
			ms1[si] = shorter[i]
			si++
		}
	}
	si = 0
	for i := 0; i < len(longer); i++ {
		if matchFlags[i] {
			ms2[si] = longer[i]
			si++
		}
	}

	t := 0
	for i, c := range ms1 {
		if c != ms2[i] {
			t++
		}
	}
	prefix := 0
	for i := 0; i < len(shorter); i++ {
		if longer[i] == shorter[i] {
			prefix++
		} else {
			break
		}
	}
	if m == 0 {
		return 0.0, 0.0
	}
	newt := float64(t) / 2.0
	newm := float64(m)
	return 1 / 3.0 * (newm/float64(len(shorter)) + newm/float64(len(longer)) + (newm-newt)/newm), prefix
}

func Min(is ...int) int {
	var min int
	for i, v := range is {
		if i == 0 || v < min {
			min = v
		}
	}
	return min
}

// Max returns the maximum number of passed int slices.
func Max(is ...int) int {
	var max int
	for _, v := range is {
		if max < v {
			max = v
		}
	}
	return max
}

func BenchmarkJaroImplementations(b *testing.B) {
	sets := []struct {
		a string
		b string
	}{
		{a: "gmilcon", b: "gmilno"},
		{a: "DIXON", b: "DICKSONX"},
		{a: "MARHTA", b: "martha"},
	}

	for _, set := range sets {
		b.Run(fmt.Sprintf("%q %q", set.a, set.b), func(b *testing.B) {
			b.Run("JaroDistanceMasatana", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_, _ = JaroDistanceMasatana(set.a, set.b)
				}
			})

			b.Run("RosettaJaro V0", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = RosettaJaroV0(set.a, set.b)
				}
			})

			b.Run("RosettaJaro V1", func(b *testing.B) {
				RosettaJaroV1 := NewJaro()
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = RosettaJaroV1(set.a, set.b)
				}
			})

			b.Run("smetrics.Jaro", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = smetrics.Jaro(set.a, set.b)
				}
			})
		})
	}
}

func BenchmarkRosettaJaro(b *testing.B) {
	sets := []struct {
		a string
		b string
	}{
		{a: "aaaaaa", b: "zzzzzz"},
		{a: "beer", b: "root"},
		{a: "beer", b: "been"},
		{a: "huffelpuf", b: "puffelhuf"},
		{a: "algorithm", b: "algoritm"},
		{a: "corn", b: "corm"},
	}

	b.Run("Double alloc", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range sets {
				_ = RosettaJaroV0(s.a, s.b)
			}
		}
	})
	b.Run("Single alloc", func(b *testing.B) {
		RosettaJaroV1 := NewJaro()
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, s := range sets {
				_ = RosettaJaroV1(s.a, s.b)
			}
		}
	})
}

// @see https://rosettacode.org/wiki/Jaro_distance#Go
func RosettaJaroV0(a, b string) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1
	}
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	matchDistance := len(a)
	if len(b) > matchDistance {
		matchDistance = len(b)
	}

	matchDistance = matchDistance/2 - 1
	aMatches := make([]bool, len(a))
	bMatches := make([]bool, len(b))

	var matches float64
	var transpositions float64
	for i := range a {
		start := i - matchDistance
		if start < 0 {
			start = 0
		}

		end := i + matchDistance + 1
		if end > len(b) {
			end = len(b)
		}

		for k := start; k < end; k++ {
			if bMatches[k] {
				continue
			}
			if a[i] != b[k] {
				continue
			}

			aMatches[i] = true
			bMatches[k] = true
			matches++
			break
		}
	}

	if matches == 0 {
		return 0
	}

	k := 0
	for i := range a {
		if !aMatches[i] {
			continue
		}

		for !bMatches[k] {
			k++
		}

		if a[i] != b[k] {
			transpositions++
		}

		k++
	}

	return (matches/float64(len(a)) +
		matches/float64(len(b)) +
		(matches-(transpositions/2))/matches) / 3
}
