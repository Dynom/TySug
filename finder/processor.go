package finder

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type Processor func(string) string

func MarkingTypoProcessor() Processor {
	return func(input string) string {
		return markingTypoReplacer(input, 0)
	}
}

func markingTypoReplacer(input string, offset uint) string {
	const (
		marker       = "^"
		markerLength = len(marker)
	)

	if offset >= uint(len(input)) {
		return input
	}

	firstIdx := strings.Index(input[offset:], "^")

	// If it's not present, or if it's the last character, there is nothing to replace
	if firstIdx == -1 || firstIdx == len(input)-1 {
		return input
	}

	firstIdx += int(offset)
	offset = uint(firstIdx)

	marking, markingWidth := utf8.DecodeRune([]byte(input[firstIdx+1:]))
	switch marking {
	default: // Unrecognised marking, advancing the index
		offset += uint(markerLength)

	case 'H': // Replace a single character to the left
		_, width := utf8.DecodeLastRune([]byte(input[:firstIdx]))
		from := firstIdx + markerLength + markingWidth
		offset = uint(firstIdx - width)
		input = input[:offset] + input[from:]

	case 'W': // Replace entire word to the left, until the first non-letter character is found
		var beforeSpace int
		for i := firstIdx; i > 0; {
			l, width := utf8.DecodeLastRuneInString(input[:i])
			if !unicode.IsLetter(l) {
				beforeSpace = i
				break
			}

			i -= width
		}

		input = input[:beforeSpace] + input[firstIdx+markerLength+markingWidth:]
		offset = uint(beforeSpace)
	}

	return markingTypoReplacer(input, offset)
}

// atLeastZero returns 0 if `i` is negative, or `i` when it's positive
func atLeastZero(i int) int {
	if i < 0 {
		return 0
	}
	return i
}
