// Package gotwant provides "got: , want: " style test functions.
package gotwant

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

var (
	FmtDefault string = "%v"
)

// Case is a test case.
type Case struct {
	Got  interface{} // what you got.
	Want interface{} // what you expected.
	Fmt  string      // used in t.Errorf displaying got and want.  default: FmtDefault
}

// C constructs a Case.
// You can construct a Case by hand.
func C(got, want interface{}, optFmt ...string) Case {
	c := Case{
		Got:  got,
		Want: want,
	}
	if len(optFmt) > 0 {
		c.Fmt = optFmt[0]
	}
	return c
}

// Test if for a single try.
func Test(t *testing.T, got, want interface{}, optFmt ...string) {
	c := C(got, want, optFmt...)

	if !reflect.DeepEqual(c.Got, c.Want) {
		valfmt := c.Fmt
		if valfmt == "" {
			valfmt = FmtDefault
		}
		_, file, line, _ := runtime.Caller(1)
		errfmt := fmt.Sprintf("%s:%d\ngot:  %s\nwant: %s", filepath.Base(file), line, valfmt, valfmt)
		t.Errorf(errfmt, c.Got, c.Want)
	}
}

// TestAll is for a series of tries(Cases).
func TestAll(t *testing.T, cases ...Case) {
	for _, c := range cases {
		if !reflect.DeepEqual(c.Got, c.Want) {
			valfmt := c.Fmt
			if valfmt == "" {
				valfmt = FmtDefault
			}
			_, file, line, _ := runtime.Caller(1)
			errfmt := fmt.Sprintf("%s:%d\ngot:  %s\nwant: %s", filepath.Base(file), line, valfmt, valfmt)
			t.Errorf(errfmt, c.Got, c.Want)
		}
	}
}
