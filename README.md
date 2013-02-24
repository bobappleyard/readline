Readline Bindings
=================

This is a set of bindings to the GNU Readline Library.

The existing readline bindings for Go are more limited than this library, if
you can believe it.

Note that the return type of String() has changed.

It was

	func String(prompt string) string
	
and it's now

	func String(prompt string) (string, error)

Installing the library
----------------------

To install the library in order to use it, type:

	go get github.com/bobappleyard/readline

You may need to be root.

For Mac OS X users, you may see errors like `rl_catch_sigwinch undeclared`. If so, you need to install GNU Readline via [Homebrew](http://mxcl.github.com/homebrew/):

	brew install readline

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
		"io"
		"fmt"
		"github.com/bobappleyard/readline"
	)

	func main() {
		for {
			l, err := readline.String("> ")
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("error: ", err)
				break
			}
			fmt.Println(l)
			readline.AddHistory(l)
		}
	}


