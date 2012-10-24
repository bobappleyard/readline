package readline_test

import (
	"fmt"
	"os"
	"os/signal"
	"github.com/bobappleyard/readline"
	"syscall"
)

func ExampleCleanup() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	
	readline.CatchSigint = false

	var line string
	var err error

	done := make(chan struct{})

	go func() {
		line, err = readline.String("> ")
		close(done)
	}()

	select {
	case <-sigint:
		fmt.Println("\nInterrupted")

		// Restore terminal attributes
		readline.Cleanup()
		// Note that we still have a goroutine reading from Stdin that
		// will terminate when we exit.
		os.Exit(1)
	case <-done:
		fmt.Printf("Read line %s, error %v\n", line, err)
	}
}

func ExampleEscapePrompt() {
	// ANSI escape sequences
    bright := "\x1b[1m"
    reset := "\x1b[0m"
	
	prompt := readline.EscapePrompt("Command: " + bright)
	line, err := readline.String(prompt)
	fmt.Print(reset) // Revert terminal to non-bright text
	fmt.Printf("Read line %s, error %v\n", line, err)
}
