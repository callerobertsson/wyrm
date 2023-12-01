// Package wyrm global commands
package wyrm

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Special keys/runes
const (
	RuneEnter  = '\n'
	RuneSpace  = ' '
	RuneScream = '!'
	RuneSlash  = '/'
	RuneClear  = '\f'
	RuneEsc    = '\x1b'
	RuneQue    = '?'
	RuneQuit   = 'q'
)

// globalKeyInfo hold help for the global command keys
var globalKeyInfo = map[rune][2]string{
	RuneSpace:  {"space", "show key info for current command"},
	RuneEnter:  {"newline", "show key info recursivly"},
	RuneScream: {"!", "execute shell command"},
	RuneClear:  {"ctrl-l", "clear screen"},
	RuneEsc:    {"escape", "abort input"},
	RuneQue:    {"?", "display detailed help"},
	RuneQuit:   {"q", "quit program nicely"},
}

// getGlobalCommands returns a rune to Command map with the global commands
func (w *Wyrm) getGlobalCommands() map[rune]*Command {
	return map[rune]*Command{
		RuneSpace: { // key info
			Description: globalKeyInfo[RuneSpace][1],
			Function:    func() error { return w.commandsHelpCommand(false) },
		},
		RuneEnter: { // key info recursivly
			Description: globalKeyInfo[RuneEnter][1],
			Function:    func() error { return w.commandsHelpCommand(true) },
		},
		RuneQue: { // question mark for help
			Description: globalKeyInfo[RuneQue][1],
			Function:    func() error { return w.detailedHelpCommand() },
		},
		RuneClear: { // ctrl-l to clear screen
			Description: globalKeyInfo[RuneClear][1],
			Function:    func() error { fmt.Printf(string("\x1b[2J\x1b[H")); return nil },
		},
		RuneEsc: { // esc for abort
			Description: globalKeyInfo[RuneEsc][1],
			Function:    func() error { w.state.cmd = w.rootCommand; return nil },
		},
		RuneScream: { // exclamation mark to execute shell command
			Description: globalKeyInfo[RuneScream][1],
			Function:    w.shellCommand,
		},
		RuneQuit: { // q to exit
			Description: globalKeyInfo[RuneQuit][1],
			Function:    w.quitCommand,
		},
	}
}

// commandsHelpCommand prints help about the commands
func (w *Wyrm) commandsHelpCommand(recursive bool) error {
	fmt.Printf("Available command keys:\n")

	pad := "    "

	// Define a recursive command info printer
	var p func(cmd *Command, indent string)
	p = func(cmd *Command, indent string) {
		states := []state{}
		for k, c := range cmd.Commands {
			states = append(states, state{k, c})
		}

		sort.Sort(stateByOrder(states))

		for _, s := range states {
			// fmt.Println("state:", string(s.key))
			key := string(s.key)
			if info, ok := globalKeyInfo[s.key]; ok {
				key = info[0]
			}
			fmt.Printf(indent+"[%s] %q - %s\n", key, s.cmd.Title, s.cmd.Description)
			if recursive {
				p(s.cmd, indent+pad)
			}
		}
	}

	// Print recursivly from current command
	p(w.state.cmd, pad)

	return nil
}

// detailedHelpCommand prints help about commands and global commands
func (w *Wyrm) detailedHelpCommand() error {
	w.commandsHelpCommand(true)
	fmt.Printf("Global command keys:\n")
	for r := range w.getGlobalCommands() {
		text := globalKeyInfo[r][1]
		if _, exists := w.state.cmd.Commands[r]; exists {
			text = "overridden for current command"
		}
		fmt.Printf("%12s - %s\n", "["+globalKeyInfo[r][0]+"]", text)
	}
	return nil
}

// shellCommand makes it possible to execute shell commands
func (w *Wyrm) shellCommand() error {

	// Get command line
	line, err := InputText(w.InputPrompt("enter shell command"), "")
	if err != nil {
		return err
	}

	parts := strings.Split(line, " ")
	bs, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		return err
	}

	fmt.Printf("%v", string(bs))

	return err
}

// quitCommand is executed to leave program
func (w *Wyrm) quitCommand() error {
	fmt.Printf("bye!\n")

	// This does not work on Darwin
	// But command works in a Darwin terminal so an alias like this:
	//   alias bu='clear && bujogo && reset'
	// can be helpful
	exec.Command("reset").Run()

	os.Exit(0)
	return nil
}

// isSpecialKey returns true if input rune is a member of global key info
func isSpecialKey(r rune) bool {
	for k := range globalKeyInfo {
		if r == k {
			return true
		}
	}
	return false
}

// specialKeys returns the special keys
func specialKeys(c map[rune]*Command) []string {
	keys := []string{}
	for k := range c {
		s := string(k)
		if isSpecialKey(k) {
			s = globalKeyInfo[k][0]
		}
		keys = append(keys, s)
	}

	sort.Strings(keys)
	return keys
}
