package wyrm

import (
	"fmt"
)

// ErrAbort indicates that the user wants to abort input
var ErrAbort = fmt.Errorf("abort")

// ErrNoNumber is returned if number input isn't a number
var ErrNoNumber = fmt.Errorf("not a number")

// ErrOutOfRange is returned if number input is out of a specific range
var ErrOutOfRange = fmt.Errorf("out of range")

// ErrNoTime is returned of entered time that doesn't mach HH:MM or HHMM
var ErrNoTime = fmt.Errorf("not a valid time value")
