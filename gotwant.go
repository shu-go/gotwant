// Package gotwant provides "got: , want: " style test functions.
package gotwant

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var (
	FmtDefault string = "%v"
)

// Case is a test case.
type Case struct {
	Got  interface{} // what you got.
	Want interface{} // what you expected.

	Fmt  string // used in t.Errorf displaying got and want.  default: FmtDefault
	Desc string // a line description
}

// Option sets optional values.
type Option func(*Case) error

// Format sets(changes) print format. default: FmtDefault
func Format(fmt string) Option {
	return func(c *Case) error {
		c.Fmt = fmt
		return nil
	}
}

// Desc sets a line description.
func Desc(desc string) Option {
	return func(c *Case) error {
		c.Desc = desc
		return nil
	}
}

// C constructs a Case.
// You can construct a Case by hand.
func C(got, want interface{}, opts ...Option) Case {
	c := Case{
		Got:  got,
		Want: want,
	}
	for _, o := range opts {
		o(&c)
	}
	return c
}

// Test if for a single try.
func Test(t *testing.T, got, want interface{}, opts ...Option) {
	t.Helper()

	test(t, C(got, want, opts...))
}

// TestAll is for a series of tries(Cases).
func TestAll(t *testing.T, cases ...Case) {
	t.Helper()

	for _, c := range cases {
		test(t, c)
	}
}

func test(t *testing.T, c Case) {
	t.Helper()

	if !reflect.DeepEqual(c.Got, c.Want) {
		valfmt := c.Fmt
		if valfmt == "" {
			valfmt = FmtDefault
		}
		errfmt := fmt.Sprintf("%s\ngot:  %s\nwant: %s", c.Desc, valfmt, valfmt)
		t.Errorf(errfmt, c.Got, c.Want)
	}
}

// Error tests given error (got) is (1) exactly the error you wanted or (2) its message matches your pattern.
// If want is a string, this func tests with strings.Contains(got, want)
// else, this func tests with reflect.DeepEqual(got, want)
func Error(t *testing.T, got, want interface{}, opts ...Option) {
	t.Helper()

	errorLike(t, C(got, want, opts...))
}

func errorLike(t *testing.T, c Case) {
	t.Helper()

	var gotErr error
	var wantErrMsg string
	var ok bool

	if c.Got != nil {
		gotErr, ok = c.Got.(error)
		if !ok {
			t.Fatal("got non-error type")
		}
	}

	if c.Want != nil {
		if _, ok = c.Want.(error); !ok {
			wantErrMsg, ok = c.Want.(string)
			if !ok {
				t.Fatal("wanted non-error type")
			}
		}
	}

	if wantErrMsg != "" {
		// compare message
		if strings.Contains(strings.ToLower(gotErr.Error()), strings.ToLower(wantErrMsg)) {
			return
		}

	} else {
		if reflect.DeepEqual(c.Got, c.Want) {
			return
		}
	}

	valfmt := c.Fmt
	if valfmt == "" {
		valfmt = FmtDefault
	}
	errfmt := fmt.Sprintf("%s\ngot error:  %s\nwant error: %s", c.Desc, valfmt, valfmt)
	t.Errorf(errfmt, c.Got, c.Want)
}
