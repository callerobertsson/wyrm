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

func main() {
	// Define the Wyrm command hiearchy using a mix of inline and var commands
	cmds := wyrm.Command{
		Title:       "wyrm",
		Description: "wyrm example program",
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
					},
					'n': {
						Title:       "number",
						Description: "input a number",
						Function:    inputNumber,
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
		},
	}

	// Create Wyrm
	w = wyrm.New(&cmds)

	fmt.Println("Wyrm Example")
	fmt.Println("use q to quit and ? for help")

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
