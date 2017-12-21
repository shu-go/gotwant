// Package gotwant provides "got: , want: " style test functions.
package gotwant

import (
	"fmt"
	"reflect"
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
