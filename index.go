// Package wyrm indexing helper
package wyrm

var indices []rune

func init() {
	// 0, .., 9
	for i := 48; i <= 57; i++ {
		indices = append(indices, rune(i))
	}
	// a, ..., z
	for i := 97; i <= 122; i++ {
		indices = append(indices, rune(i))
	}
	// A, ..., Z
	for i := 65; i <= 90; i++ {
		indices = append(indices, rune(i))
	}
}

// GetIndexRunes returns the list of runes used as indices
func GetIndexRunes() []rune {
	return indices
}

// GetRuneIndex returns the index for the rune
func GetRuneIndex(r rune) (int, error) {
	for i, ri := range indices {
		if r == ri {
			return i, nil
		}
	}

	return 0, ErrOutOfRange
}

// GetIndexRune returns the rune at index i
func GetIndexRune(i int) (rune, error) {
	if i < 0 || i > len(indices)-1 {
		return ' ', ErrOutOfRange
	}

	return indices[i], nil
}
