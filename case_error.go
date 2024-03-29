package gotwant

import (
	"fmt"
	"reflect"
	"strings"
)

type errCase struct {
	Got  error       // what you got.
	Want interface{} // an error-type error or a message

	Fmt  string // used in t.Errorf displaying got and want.  default: FmtDefault
	Desc string // a line description
}

// Error constructs a error-comaration(nil, string) test case.
func Error(got error, want interface{}, opts ...Option) *errCase {
	c := &errCase{
		Got:  got,
		Want: want,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *errCase) SetFmt(format string) {
	c.Fmt = format
}

func (c *errCase) SetDesc(desc string) {
	c.Desc = desc
}

func (c *errCase) Test(t T) {
	t.Helper()

	valfmt := c.Fmt
	if valfmt == "" {
		valfmt = FmtDefault
	}

	if c.Got == nil && c.Want == nil {
		return
	}

	if c.Got == nil {
		errfmt := fmt.Sprintf("%s\ngot NO error.\nwant error: %s", c.Desc, valfmt)
		t.Errorf(errfmt, c.Want)
		return
	}

	wantErrMsg := stringify(c.Want)
	if wantErrMsg != nil {
		// compare message
		if strings.Contains(strings.ToLower(c.Got.Error()), strings.ToLower(*wantErrMsg)) {
			return
		}

	} else {
		if reflect.DeepEqual(c.Got, c.Want) {
			return
		}
	}

	errfmt := fmt.Sprintf("%s\ngot error:  %s\nwant error: %s", c.Desc, valfmt, valfmt)
	t.Errorf(errfmt, c.Got, c.Want)
}
