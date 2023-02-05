// Package wyrm default prompt functions and Prompter interface
package wyrm

import (
	"fmt"
	"strings"
)

// Prompter defines an interface for prompts
type Prompter interface {
	CommandPrompt() string
	InputPrompt(string) string
	RunePrompt(string) string
}

// CommandPrompt returns the prompt with state information
func (w *Wyrm) CommandPrompt() string {
	if w.prompter != nil {
		return w.prompter.CommandPrompt()
	}

	d := w.state.cmd.Title
	k := "[" + strings.Join(w.GetCurrentKeyStrings(), "") + "] "
	return fmt.Sprintf("%s %s$ ", d, k)
}

// InputPrompt is used when entering strings
func (w *Wyrm) InputPrompt(p string) string {
	if w.prompter != nil {
		return w.prompter.InputPrompt(p)
	}

	return fmt.Sprintf("%s [%s] > ", w.state.cmd.Title, p)
}

// RunePrompt is used when entering a single rune
func (w *Wyrm) RunePrompt(p string) string {
	if w.prompter != nil {
		return w.prompter.RunePrompt(p)
	}

	return fmt.Sprintf("%s [%s] # ", w.state.cmd.Title, p)
}
