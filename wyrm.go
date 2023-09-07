// Package wyrm defines the wyrm struct
package wyrm

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
)

// Wyrm is the quick command handler
type Wyrm struct {
	rootCommand *Command // the root of all evil
	state       state    // the current state
	prompter    Prompter // prompt printer interface
}

// Command has a description, function and a map of sub commands.
// Used to build up a command hiarchy.
type Command struct {
	Title       string
	Description string
	Sort        int
	Commands    map[rune]*Command
	Parent      *Command
	Function    func() error
	Pre         func() error
	Post        func() error
}

// state struct holds the internal current state of Wyrm
type state struct {
	key rune     // pressed key
	cmd *Command // current command
}

type stateByOrder []state

func (a stateByOrder) Len() int      { return len(a) }
func (a stateByOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a stateByOrder) Less(i, j int) bool {
	// fmt.Printf("comparing %v (%v) and %v (%v): ", a[i].cmd.Sort, string(a[i].key), a[j].cmd.Sort, string(a[j].key))
	if a[i].cmd.Sort == 0 {
		return false // 0 is never less than anything
	}
	if a[j].cmd.Sort == 0 {
		return true // 0 is never less than anything
	}
	if a[i].cmd.Sort < a[j].cmd.Sort {
		return true // normal compare
	}
	return false
}

// New creates a new wyrm
func New(rootCommand *Command) *Wyrm {
	w := Wyrm{
		rootCommand,
		state{
			key: rune(' '),
			cmd: rootCommand,
		},
		nil,
	}

	return &w
}

// SetPrompter sets the Prompter to use instead of the default
func (w *Wyrm) SetPrompter(p Prompter) {
	w.prompter = p
}

// GetCurrentCommand returns the active command
func (w *Wyrm) GetCurrentCommand() *Command {
	return w.state.cmd
}

// GetCurrentKey returns the current key pressed
func (w *Wyrm) GetCurrentKey() rune {
	return w.state.key
}

// GetCurrentKeyStrings returns the keys, no special keys, of the Command as stings in alphabetically order
func (w *Wyrm) GetCurrentKeyStrings() []string {
	keys := []string{}

	states := []state{}
	for r, c := range w.state.cmd.Commands {
		states = append(states, state{r, c})
	}

	sort.Sort(stateByOrder(states))

	for _, s := range states {
		if isSpecialKey(s.key) {
			continue
		}
		keys = append(keys, string(s.key))
	}

	return keys
}

// Run starts the command line interface
func (w *Wyrm) Run() {

	// Disable buffering and set no display
	f := "-F"
	if runtime.GOOS == "darwin" {
		// Ugly hack because Macos (Darwin) needs -f iso -F
		f = "-f"
	}

	exec.Command("stty", f, "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", f, "/dev/tty", "-echo").Run()

	// Loop until quit
	for {
		// Prompt
		input, err := InputRune(w.CommandPrompt())
		if err == ErrAbort {
			w.state.cmd = w.state.cmd.Parent
			if w.state.cmd == nil {
				w.state.cmd = w.rootCommand
			}
			continue
		}

		// Get sub command of current command
		cmd, ok := w.state.cmd.Commands[input]
		if ok {
			w.state.key = input // remember key

			// Switch to new command
			cmd.Parent = w.state.cmd
			w.state.cmd = cmd

			// Execute Pre if present
			if w.state.cmd.Pre != nil {
				if err := w.state.cmd.Pre(); err != nil {
					fmt.Printf("Error: %s\n", err)
					w.state.cmd = w.rootCommand
					continue
				}
			}

			// Execute function if present
			if w.state.cmd.Function != nil {
				err := w.state.cmd.Function()
				switch {
				case err == ErrAbort:
					w.state.cmd = w.state.cmd.Parent
					if w.state.cmd == nil {
						w.state.cmd = w.rootCommand
					}
					continue
				case err == nil:
					if w.state.cmd.Post != nil {
						if err := w.state.cmd.Post(); err != nil {
							fmt.Printf("Error: %s\n", err)
							w.state.cmd = w.rootCommand
							continue
						}
					}
				default:
					fmt.Printf("Error: %s\n", err)
				}

				// Return to root command, if no sub commands
				if len(w.state.cmd.Commands) < 1 {
					w.state.cmd = w.rootCommand
					continue
				}
			}

			// If no function continue processing the new command
			continue
		}

		// Check global commands (can be override above)
		if cmd, ok = w.getGlobalCommands()[input]; ok {
			if cmd.Function == nil {
				fmt.Println("No function defined")
				continue
			}

			err := cmd.Function()
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			continue
		}

		fmt.Printf("Unknown command %s\n", string(input))
	}
}
