/*
This package provides access to basic GNU Readline functions. Currently supported are:

	- getting text from a prompt (via the String() and NewReader() functions).
	- managing the prompt's history (via the AddHistory(), GetHistory(), ClearHistory() and HistorySize() functions).
	- controlling tab completion (via the Completer variable).

Here is a simple example:

	package main

	import (
	    "fmt"
	    "io"
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
*/
package readline

/*

#cgo LDFLAGS: -lreadline -lhistory

#include <stdio.h>
#include <stdlib.h>
#include <readline/readline.h>
#include <readline/history.h>

extern char *_completion_function(char *s, int i);

static char *_completion_function_trans(const char *s, int i) {
	return _completion_function((char *) s, i);
}

static void register_readline() {
	rl_completion_entry_function = _completion_function_trans;
	using_history();
}

*/
import "C"

import (
	"io"
	"unsafe"
	"syscall"
)

// The default prompt used by Reader().
var Prompt = "> "

// The continue prompt used by Reader().
var Continue = ".."

type state byte

const (
	readerStart state = iota
	readerContinue
	readerEnd
)

type reader struct {
	buf []byte
	state state
}

// Begin reading lines. If more than one line is required, the continue prompt
// is used for subsequent lines.
func NewReader() io.Reader {
	return new(reader)
}

func (r *reader) getLine() error {
	prompt := Prompt
	if r.state == readerContinue {
		prompt = Continue
	}
	s, err := String(prompt)
	if err != nil {
		return err
	}
	r.buf = []byte(s)
	return nil
}

func (r *reader) Read(buf []byte) (int, error) {
	if r.state == readerEnd {
		return 0, io.EOF
	}
	if len(r.buf) == 0 {
		err := r.getLine()
		if err == io.EOF {
			r.state = readerEnd
		}
		if err != nil {
			return 0, err
		}
		r.state = readerContinue
	}
	copy(buf, r.buf)
	l := len(buf)
	if len(buf) > len(r.buf) {
		l = len(r.buf)
	}
	r.buf = r.buf[l:]
	return l, nil
}

// Read a line with the given prompt.
func String(prompt string) (string, error) {
	p := C.CString(prompt)
	rp := C.readline(p)
	s := C.GoString(rp)
	C.free(unsafe.Pointer(p))
	if rp != nil {
		C.free(unsafe.Pointer(rp))
		return s, nil
	}
	return s, io.EOF
}

// This function provides entries for the tab completer.
var Completer = func(query, ctx string) []string {
	return nil
}

var entries []*C.char

//export _completion_function
func _completion_function(p *C.char, _i C.int) *C.char {
	C.rl_completion_suppress_append = 1
	i := int(_i)
	if i == 0 {
		es := Completer(C.GoString(p), C.GoString(C.rl_line_buffer))
		entries = make([]*C.char, len(es))
		for i, x := range es {
			entries[i] = C.CString(x)
		}
	}
	if i >= len(entries) {
		return nil
	}
	return entries[i]
}

func SetWordBreaks(cs string) {
	C.rl_completer_word_break_characters = C.CString(cs)
}

// Add an item to the history.
func AddHistory(s string) {
	n := HistorySize()
	if n == 0 || s != GetHistory(n - 1) {
		C.add_history(C.CString(s))
	}
}

// Retrieve a line from the history.
func GetHistory(i int) string {
	e := C.history_get(C.int(i+1))
	if e == nil {
		return ""
	}
	return C.GoString(e.line)
}

// Deletes all the items in the history.
func ClearHistory() {
	C.clear_history()
}

// Returns the number of items in the history.
func HistorySize() int {
	return int(C.history_length)
}

// Load the history from a file.
func LoadHistory(path string) error {
	p := C.CString(path)
	e := C.read_history(p)
	C.free(unsafe.Pointer(p))

	if e == 0 {
		return nil
	}
	return syscall.Errno(e)
}

// Save the history to a file.
func SaveHistory(path string) error {
	p := C.CString(path)
	e := C.write_history(p)
	C.free(unsafe.Pointer(p))

	if e == 0 {
		return nil
	}
	return syscall.Errno(e)
}

// Frees internal memory and restores terminal attributes. This
// function should be called when readline doesn't return and would
// leave the terminal in a corrupted state.
func Cleanup() {
	C.rl_free_line_state()
	C.rl_cleanup_after_signal()
}

// This function should be called when the SIGWINCH signal is not
// handled by readline, to update readline's internal notion of the
// screen size.
func Resize() {
	// Print the 'reset' ANSI escape code, so that the current prompt
	// ANSI codes won't corrupt the start of the refreshed prompt
	fmt.Print("\x1b[0m")

	C.rl_resize_terminal()
}

// Sets if readline should catch the SIGWINCH signal emitted when the
// terminal size changes. If set to false, the Resize() function
// should be called manually.
func CatchResize(catch bool) {
	if catch {
		C.rl_catch_sigwinch = 1
	} else {
		C.rl_catch_sigwinch = 0
	}
}

func init() {
	C.register_readline()
}
