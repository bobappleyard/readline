// Access to basic readline functions. Support for history manipulation and tab
// completion are provided. Not much of the library is exposed though.
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
)

// The default prompt used by Reader().
var Prompt = "> "

// The continue prompt used by Reader().
var Continue = ".."

type reader []byte

// Begin reading lines. If more than one line is required, the continue prompt
// is used for subsequent lines.
func Reader() io.Reader {
	return new(reader).init()
}

func (r *reader) init() *reader {
	r.getLine(Prompt)
	return r
}

func (r *reader) getLine(prompt string) error {
	s, err := String(prompt)
	if err != nil {
		return err
	}
	*r = []byte(s)
	return nil
}

func (r *reader) Read(buf []byte) (int, error) {
	if len(*r) == 0 {
		err := r.getLine(Continue)
		if err != nil {
			return 0, err
		}
	}
	copy(buf, *r)
	l := len(buf)
	if len(buf) > len(*r) {
		l = len(*r)
	}
	*r = (*r)[l:]
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

// Load the history from a file. Returns whether or not this was successful.
func LoadHistory(path string) bool {
	p := C.CString(path)
	e := C.read_history(p)
	C.free(unsafe.Pointer(p))
	return e == 0
}

// Save the history to a file. Returns whether or not this was successful.
func SaveHistory(path string) bool {
	p := C.CString(path)
	e := C.write_history(p)
	C.free(unsafe.Pointer(p))
	return e == 0
}

func init() {
	C.register_readline()
}
