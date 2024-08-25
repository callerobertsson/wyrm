package wyrm

import (
	"testing"
)

func TestParseHourMin(t *testing.T) {
	cases := []struct {
		s   string
		eh  int
		em  int
		t   string
		err error
	}{
		{"12:34", 12, 34, "", nil},
		{"12:34 ", 12, 34, "", nil},
		{" 12:34 ", 12, 34, "", nil},
		{"12:34 affen ", 12, 34, "affen", nil},
		{"12:34 one two", 12, 34, "one two", nil},
		{"HELO", 0, 0, "", ErrNoTime},
		{"12", 0, 0, "", ErrNoTime},
		{":34", 0, 0, "", ErrNoTime},
		{"", 0, 0, "", ErrNoTime},
	}

	for _, c := range cases {
		rh, rm, tail, err := parseHourMin(c.s)
		if err != c.err {
			t.Errorf("parseHourMin error %q, expected %q", err, c.err)
			continue
		}
		if rh != c.eh || rm != c.em || tail != c.t {
			t.Errorf("parseHourMin(%q) = (%d, %d, %q), expected (%d, %d, %q)\n",
				c.s, rh, rm, tail, c.eh, c.em, c.t)
		}
	}
}
