// Package wyrm input readers
package wyrm

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

// Regexp for time strings, HH:MM or HHMM
var reHourMin = regexp.MustCompile(`(\d{2}):?(\d{2})(.*)`)

// InputRune read a single rune
func InputRune(p string) (rune, error) {
	fmt.Printf(p)
	buf := make([]byte, 1)
	os.Stdin.Read(buf)
	fmt.Println("")
	r := rune(buf[0])
	if r == RuneEsc {
		return r, ErrAbort
	}
	return r, nil
}

// InputText prints prompt and reads input from user
func InputText(p string, def string) (input string, err error) {
	r, err := readline.New(p)
	if err != nil {
		return input, err
	}
	r.Operation.SetBuffer(def)

	input, err = r.Readline()
	switch {
	case err == readline.ErrInterrupt:
		return input, ErrAbort
	case err == io.EOF:
		return input, ErrDone
	case err != nil:
		return input, err
	case input == "":
		return input, ErrEmpty
	}

	return strings.TrimSpace(input), nil
}

// InputInt read an integer in the range 0 to max from the user
func InputInt(p, def string, max int) (i int, err error) {
	// Read number as string
	input, err := InputText(p, def)
	if err != nil {
		return i, err
	}

	// Convert to int
	i, err = strconv.Atoi(input)
	if err != nil {
		return i, ErrNoNumber
	}
	if i < 0 || i > max {
		return i, ErrOutOfRange
	}

	return i, nil
}

// InputTime read a time input formatted as HH:MM or HHMM
// Return time as a string and any additional characters as tail
func InputTime(p, def string) (time, tail string, err error) {
	// Read input
	input, err := InputText(p, def)
	if err != nil {
		return time, tail, err
	}

	// Parse hour and minute
	h, m, tail, err := parseHourMin(input)
	if err != nil {
		return time, tail, ErrNoTime
	}

	return fmt.Sprintf("%02d:%02d", h, m), tail, nil
}

// parseHourMin parses a string and return hour and minute as string and the tail as string.
func parseHourMin(s string) (int, int, string, error) {
	ms := reHourMin.FindAllStringSubmatch(s, 2)
	if len(ms) != 1 || len(ms[0]) != 4 {
		return 0, 0, "", ErrNoTime
	}

	// Validate hour and minute
	h, err := strconv.Atoi(ms[0][1])
	if err != nil || h < 0 || h > 23 {
		return 0, 0, "", ErrNoTime
	}
	m, err := strconv.Atoi(ms[0][2])
	if err != nil || m < 0 || m > 59 {
		return 0, 0, "", ErrNoTime
	}

	tail := strings.TrimSpace(ms[0][3])

	return h, m, tail, nil
}
