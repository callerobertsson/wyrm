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
	Commands    map[rune]*Command
	Parent      *Command
	Function    func() error
}

// state struct holds the internal current state of Wyrm
type state struct {
	cmd *Command // current command
	key rune     // pressed key
}

// New creates a new wyrm
func New(rootCommand *Command) *Wyrm {
	w := Wyrm{
		rootCommand,
		state{
			cmd: rootCommand,
			key: rune(' '),
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
	for k := range w.state.cmd.Commands {
		if isSpecialKey(k) {
			continue
		}
		keys = append(keys, string(k))
	}

	sort.Strings(keys)
	return keys
}

// Run starts the command line interface
func (w *Wyrm) Run() {

	// Disable buffering and set no display
	// Ugly hack because Macos (Darwin) needs -f and Linux -F
	f := "-F"
	if runtime.GOOS != "linux" {
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

			// Execute function if present
			if w.state.cmd.Function != nil {
				err := w.state.cmd.Function()
				if err == ErrAbort {
					w.state.cmd = w.state.cmd.Parent
					if w.state.cmd == nil {
						w.state.cmd = w.rootCommand
					}
					continue
				}
				if err != nil {
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
