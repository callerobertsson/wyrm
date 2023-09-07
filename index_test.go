package wyrm

import (
	"testing"
)

func TestGetIndexRunes(t *testing.T) {
	exp := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	res := GetIndexRunes()
	if string(res) != exp {
		t.Errorf("GetIndexRunes = %q, expected %q", string(res), exp)
	}
}

func TestGetRuneIndex(t *testing.T) {
	cases := []struct {
		r   rune
		i   int
		err error
	}{
		{'0', 0, nil},
		{'9', 9, nil},
		{'a', len("0123456789"), nil},
		{'z', 9 + len("abcdefghijklmnopqrstuvwxyz"), nil},
		{'A', 36, nil},
		{'Z', 35 + len("ABCDEFGHIJKLMNOPQRSTUVWXYZ"), nil},
		{'?', 0, ErrOutOfRange},
	}

	for _, c := range cases {
		res, err := GetRuneIndex(c.r)
		if err == c.err {
			continue
		}
		if err != c.err {
			t.Errorf("GetIndexRune error %q, expected %q", err, c.err)
			continue
		}
		if res != c.i {
			t.Errorf("GetRuneIndex(%s) = %d, expected %d\n", string(c.r), res, c.i)
		}
	}
}

func TestIndexRune(t *testing.T) {
	cases := []struct {
		i   int
		r   rune
		err error
	}{
		{1, '1', nil},
		{11, 'b', nil},
		{37, 'B', nil},
		{200, 'x', ErrOutOfRange},
		{-1, ' ', ErrOutOfRange},
	}

	for _, c := range cases {
		res, err := GetIndexRune(c.i)
		if err == c.err {
			continue
		}
		if err != c.err {
			t.Errorf("GetIndexRune error %q, expected %q", err, c.err)
			continue
		}
		if res != c.r {
			t.Errorf("GetIndexRune(%d) = %q, expected %q\n", c.i, string(res), string(c.r))
		}
	}
}
