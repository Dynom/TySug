package keyboard

import (
	"math"
)

type KeyboardLayout string

const (
	Default  KeyboardLayout = QwertyUS
	QwertyUS KeyboardLayout = "qwerty-us"
)

type keyGrid map[string]coordinates

var (

	// @todo this design currently ignores the possibility of pressing the shift key while typing
	// we might want to allow printable symbols with the same coordinates as their un-shifted counterfeit
	keyboardLayouts = map[KeyboardLayout][]string{
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

type KeyDist struct {
	grid keyGrid
}

func New(layout KeyboardLayout) KeyDist {
	return KeyDist{
		grid: generateKeyGrid(keyboardLayouts[layout]),
	}
}

func (kd KeyDist) FindNearest(input string, list []string) (string, float64) {
	var bestScore = math.Inf(1)
	var result string

	for _, ref := range list {

		var score float64

		// Scanning each letter of this ref
		for i := 0; i < len(input); i++ {
			if i >= len(ref) {
				// @todo missing characters should have a cost, decide on a correct punishment value
				score += 1
				continue
			}

			if input[i] != ref[i] {
				left, right := string(input[i]), string(ref[i])
				score += getDistance(kd.grid[left], kd.grid[right])
			}
		}

		if score < bestScore {
			bestScore = score
			result = ref
		}
	}

	return result, bestScore
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
