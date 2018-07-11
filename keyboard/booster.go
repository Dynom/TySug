package keyboard

import (
	"math"
)

var (
	keyGrid map[string]coordinates

	// @todo this design currently ignores the possibility of pressing the shift key while typing
	// we might want to allow printable symbols with the same coordinates as their un-shifted counterfeit
	keyboardLayouts = map[string][]string{
		"qwerty-us": {
			"`1234567890-=",
			" qwertyuiop[]\\",
			" asdfghjkl;'",
			" zxcvbnm,./",
		},
	}
)

type coordinates struct {
	X float64
	Y float64
}

func init() {
	keyGrid = generateKeyGrid(keyboardLayouts["qwerty-us"])
}

func GetBestMatch(input string, list []string) (string, float64) {
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
				score += getDistance(keyGrid[left], keyGrid[right])
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

func generateKeyGrid(rows []string) map[string]coordinates {
	var keyMap = make(map[string]coordinates, len(rows))

	for rowIndex, row := range rows {
		for column := 0; column < len(row); column++ {
			character := string(row[column])

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
