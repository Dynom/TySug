package keyboard

import (
	"math"
)

var (
	keyDistances    map[string]Coordinates
	keyboardLayouts = map[string][]string{
		"qwerty-us": {
			"`1234567890-=",
			" qwertyuiop[]\\",
			" asdfghjkl;'",
			" zxcvbnm,./",
		},
	}
)

type Coordinates struct {
	X float64
	Y float64
}

func init() {
	keyDistances = generateKeyDistance(keyboardLayouts["qwerty-us"])
}

func GetBestMatch(input string, list []string) (string, float64) {
	var scores = make([]float64, len(list))
	for listOffset, ref := range list {
		for i := 0; i < len(input); i++ {
			if i >= len(ref) {
				// @todo missing characters, we should add a penalty
				scores[listOffset] += 1
				continue
			}

			left, right := string(input[i]), string(ref[i])
			scores[listOffset] += getDistance(keyDistances[left], keyDistances[right])
		}
	}

	var bestScore = math.Inf(1)
	var offset int
	for listOffset, score := range scores {
		if score < bestScore {
			bestScore = score
			offset = listOffset
		}
	}

	return list[offset], scores[offset]
}

func getDistance(a, b Coordinates) float64 {
	return math.Sqrt(
		math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2),
	)
}

func generateKeyDistance(rows []string) map[string]Coordinates {
	var keyMap = make(map[string]Coordinates, len(rows))

	for rowIndex, row := range rows {
		for column := 0; column < len(row); column++ {
			character := string(row[column])
			keyMap[character] = Coordinates{
				X: float64(column),
				Y: float64(rowIndex),
			}
		}
	}

	return keyMap
}
