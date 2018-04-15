package gotwant

import (
	"fmt"
	"reflect"
	"testing"
)

// Case constructs a value-comaration test case.
func Case(got, want interface{}, opts ...Option) *cmpCase {
	c := &cmpCase{
		Got:  got,
		Want: want,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

type cmpCase struct {
	Got  interface{} // what you got.
	Want interface{} // what you expected.

	Fmt  string // used in t.Errorf displaying got and want.  default: FmtDefault
	Desc string // a line description
}

func (c *cmpCase) SetFmt(format string) {
	c.Fmt = format
}

func (c *cmpCase) SetDesc(desc string) {
	c.Desc = desc
}

func (c *cmpCase) Test(t *testing.T) {
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
