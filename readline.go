// Access to basic readline functions. Support for history manipulation and tab
// completion are provided. Not much of the library is exposed though.
package readline

/*

#cgo LDFLAGS: -lreadline -lhistory

#include <stdio.h>
#include <stdlib.h>
#include <readline/readline.h>
#include <readline/history.h>

extern char *_completion_function(const char *s, int i);

static void register_readline() {
	rl_completion_entry_function = _completion_function;
	using_history();
}

*/
import "C"

import (
	"io"
	"os"
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

func (r *reader) getLine(prompt string) {
	*r = []byte(String(prompt) + "\n")
}

func (r *reader) Read(buf []byte) (int, os.Error) {
	l := len(*r)
	m := len(buf)
	if m > l {
		copy(buf, *r)
		r.getLine(Continue)
		copy(buf[l:], *r)
		*r = (*r)[m-l:]
	} else {
		copy(buf, *r)
		*r = (*r)[m:]
	}	
	return m, nil
}

// Read a line with the given prompt.
func String(prompt string) string {
	p := C.CString(prompt)
	rp := C.readline(p)
	s := C.GoString(rp)
	C.free(unsafe.Pointer(p))
	if rp != nil {
		C.free(unsafe.Pointer(rp))
	}
	return s
}

type TabCompleter interface {
	Complete(query string) []string
}

type CompleterFunc func(string) []string

func (f CompleterFunc) Complete(query string) []string {
	return f(query)
}

// This function is called whenever tab completion is invoked. Given a query
// string, it should return all items that match.
var Completer TabCompleter = CompleterFunc(func(query string) []string {
	return nil
})

var entries []*C.char

//export _completion_function
func _completion_function(p *C.char, _i C.int) *C.char {
	i := int(_i)
	if i == 0 {
		es := Completer.Complete(C.GoString(p))
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

// Add an item to the history.
func AddHistory(s string) {
	C.add_history(C.CString(s))
}

// Retrieve a line from the history.
func GetHistory(i int) string {
	e := C.history_get(C.int(i))
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

func init() {
	C.register_readline()
}

