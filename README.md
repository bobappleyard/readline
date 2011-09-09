Readline Bindings
=================

This is a set of bindings to the GNU Readline Library.

The existing readline bindings for Go are more limited than this library, if
you can believe it.

Installing the library
----------------------

To install the library in order to use it, type

	goinstall github.com/bobappleyard/readline

To install the library in order to hack on it, type

	git clone git://github.com/bobappleyard/readline.git

Using the library
-----------------

	import "github.com/bobappleyard/readline"

These bindings provide access to three basic features of Readline:

- Getting text from a prompt (via the String() and Reader() functions).
- Managing the prompt's history (via the AddHistory(), GetHistory(), 
  ClearHistory() and HistorySize() functions).
- Controlling tab completion (via the Completer variable).

An example of the library's use:

	package main

	import (
		"fmt"
		"github.com/bobappleyard/readline"
	)

	func main() {
		for {
			l := readline.String("> ")
			if l == "exit" {
				break
			}
			fmt.Println(l)
			readline.AddHistory(l)
		}
	}


