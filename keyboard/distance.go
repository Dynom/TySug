package keyboard

import (
	"math"
)

// Layout is the type used to define keyboard layouts
type Layout string

// Predefined keyboard layouts
const (
	Default  Layout = QwertyUS
	QwertyUS Layout = "qwerty-us"
)

const missingCharPenalty = 3 // @todo arbitrary penalty, find a proper basis for this

type keyGrid map[string]coordinates

var (

	// @todo this design currently ignores the possibility of pressing the shift key while typing
	// we might want to allow printable symbols with the same coordinates as their un-shifted counterfeit
	keyboardLayouts = map[Layout][]string{
		QwertyUS: {
			"`1234567890-=",
			" qwertyuiop[]\\",
			" asdfghjkl;'",
			" zxcvbnm,./",
		},
		/*
			"azerty-fr": {
				"&é\"'(-è çà)=",
				"azertyuiop $",
				"qsdfghjklmù*",
				"<wxcvbn,;:!",
			},
		*/
	}
)

type coordinates struct {
	X float64
	Y float64
}

// KeyDist is the type that allows to find the best alternative based on keyboard layouts
type KeyDist struct {
	grid keyGrid
}

// New produces a new instance of KeyDist, based on the keyboard layout you choose
func New(l Layout) KeyDist {
	return KeyDist{
		grid: generateKeyGrid(keyboardLayouts[l]),
	}
}

// FindNearest finds the item in the list that is nearest to the input, based on the keyboard layout
func (kd KeyDist) FindNearest(input string, list []string) (string, float64) {
	var bestScore = math.Inf(1)
	var result string

	for _, ref := range list {
		score := kd.CalculateDistance(input, ref)
		if score < bestScore {
			bestScore = score
			result = ref
		}
	}

	return result, bestScore
}

// CalculateDistance calculates the total distances of the reference to the input
func (kd KeyDist) CalculateDistance(input, ref string) float64 {
	var score float64

	// Scanning each letter of this ref
	for i := 0; i < len(input); i++ {
		if i >= len(ref) {

			// @todo missing characters should have a cost, decide on a correct punishment value
			score += float64(missingCharPenalty * (len(input) - len(ref)))
			break
		}

		if input[i] == ref[i] {
			continue
		}

		left, right := input[i:i+1], ref[i:i+1]
		score += getDistance(kd.grid[left], kd.grid[right])

	}

	return score
}

func getDistance(a, b coordinates) float64 {
	return math.Sqrt(
		math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2),
	)
}

func generateKeyGrid(rows []string) keyGrid {
	var keyMap = make(keyGrid, len(rows))

	for rowIndex, row := range rows {
		for column, char := range row {
			character := string(char)
			if character == " " {
				continue
			}

			keyMap[character] = coordinates{
				X: float64(column),
				Y: float64(rowIndex),
			}
		}
	}

	return keyMap
}
