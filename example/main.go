// Package main implement a simple Wyrm quick command executioner example
package main

import (
	"fmt"

	"github.com/callerobertsson/wyrm"
)

var w *wyrm.Wyrm

// helloCmd prints hello world!
var helloCmd = wyrm.Command{
	Title:       "hello",
	Description: "print hello world",
	Sort:        1, // want this to be first
	Function:    func() error { fmt.Println("hello world!"); return nil },
	Pre:         func() error { fmt.Println("(pre command printed this)"); return nil },
	Post:        func() error { fmt.Println("(post command printed this)"); return nil },
}

// abortCmd always returns ErrAbout
var abortCmd = wyrm.Command{
	Title:       "abort",
	Description: "prints message and returns ErrAbout",
	Function:    func() error { fmt.Println("aborting"); return wyrm.ErrAbort },
}

// inputTimeCmd prompts the user for a time and prints it
var inputTimeCmd = wyrm.Command{
	Title:       "time",
	Description: "input a time as HH:MM or HHMM",
	Function:    inputTime,
}

// errorCmd shows Pre and Command errors
var errorCmd = wyrm.Command{
	Title:       "errors",
	Description: "select what error to show",
	Commands: map[rune]*wyrm.Command{
		'<': {
			Title:       "pre",
			Description: "Show a Pre function error",
			Pre:         func() error { return fmt.Errorf("Planned Pre Error") },
			Sort:        1,
			Function:    func() error { fmt.Println("This should not be shown"); return nil },
		},
		'c': {
			Title:       "command",
			Description: "Show a command error",
			Sort:        2,
			Function:    func() error { return fmt.Errorf("Planned Command Error") },
		},
		'>': {
			Title:       "post",
			Description: "Show a Post function error",
			Post:        func() error { return fmt.Errorf("Planned Post Error") },
			Sort:        3,
			Function:    func() error { fmt.Println("Correct output"); return nil },
		},
	},
}

func main() {
	// Define the Wyrm command hiearchy using a mix of inline and var commands
	cmds := wyrm.Command{
		Title:       "wyrm",
		Description: "wyrm example program",
		Pre:         func() error { fmt.Println("Root Pre (could clear screen)"); return nil },
		Commands: map[rune]*wyrm.Command{
			'i': {
				Title:       "input",
				Description: "input different stuff using sub commands",
				// no Sort value -> place it last somewhere
				Commands: map[rune]*wyrm.Command{
					's': {
						Title:       "string",
						Description: "input a string",
						Function:    inputText,
						Post:        func() error { fmt.Println("text was inputted"); return nil },
					},
					'n': {
						Title:       "number",
						Description: "input a number",
						Function:    inputNumber,
						Pre:         func() error { fmt.Println("pre number selection"); return nil },
					},
					't': &inputTimeCmd,
					wyrm.RuneSpace: {
						Title:       "extra",
						Description: "test another level",
						Commands: map[rune]*wyrm.Command{
							'a': &abortCmd,
						},
					},
				},
			},
			's': {
				Title:       "select",
				Description: "select by index",
				Sort:        2, // want this to be second
				Function:    selectIndex,
			},
			'h': &helloCmd, // Sort: 1
			'e': &errorCmd, // Sort: 0 -> put at the end somewhere
		},
	}

	// Create Wyrm
	w = wyrm.New(&cmds)

	fmt.Println("Wyrm Example")
	fmt.Println("use q to quit and ? for help")

	// Manually run Pre command for root
	w.GetCurrentCommand().Pre()

	// Run Wyrm
	w.Run()
}

func inputText() error {
	input, err := wyrm.InputText(w.InputPrompt("enter text"), "")
	if err != nil {
		return err
	}

	fmt.Printf("Your entered: %q\n", input)

	return nil
}

func inputNumber() error {
	input, err := wyrm.InputInt(w.InputPrompt("enter number less than 10"), "", 10)
	if err != nil {
		return err
	}

	fmt.Printf("Your entered: %v\n", input)

	return nil
}

func inputTime() error {
	input, err := wyrm.InputTime(w.InputPrompt("HH:MM"), "12:34")
	if err != nil {
		return err
	}

	fmt.Printf("Your entered: %v\n", input)

	return nil
}

func selectIndex() error {
	options := map[rune]string{
		'1': "one option",
		'a': "another option",
	}

	fmt.Println("Select:")
	for k, v := range options {
		fmt.Printf("  %v: %v\n", string(k), v)
	}

	r, err := wyrm.InputRune(w.InputPrompt("select index"))
	if err != nil {
		return err
	}

	v, ok := options[r]
	if !ok {
		return fmt.Errorf("no match for %q", string(r))
	}

	fmt.Printf("You selected %q\n", v)

	return nil
}
