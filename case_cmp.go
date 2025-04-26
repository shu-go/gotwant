package gotwant

import (
	"fmt"
	"reflect"
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

func (c *cmpCase) Test(t T) {
	t.Helper()

	if !reflect.DeepEqual(c.Got, c.Want) {
		valfmt := c.Fmt
		if valfmt == "" {
			valfmt = FmtDefault
			for _, f := range []string{FmtDefault, "%#v", "%T"} {
				fmted1, fmted2 := fmt.Sprintf(f, c.Got), fmt.Sprintf(f, c.Want)
				if fmted1 != fmted2 {
					valfmt = f
					break
				}
			}
		}

		got := indent(fmt.Sprintf(fmt.Sprintf("got:  %s", valfmt), c.Got))
		want := indent(fmt.Sprintf(fmt.Sprintf("want: %s", valfmt), c.Want))
		t.Errorf("%s\n%s\n%s", c.Desc, got, want)
	}
}
