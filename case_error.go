package gotwant

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type errCase struct {
	Got  error       // what you got.
	Want interface{} // an error-type error or a message

	Fmt  string // used in t.Errorf displaying got and want.  default: FmtDefault
	Desc string // a line description
}

// Case constructs a value-comaration test case.
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

func (c *errCase) Test(t *testing.T) {
	t.Helper()

	var wantErrMsg string
	var ok bool

	if c.Want != nil {
		if _, ok = c.Want.(error); !ok {
			wantErrMsg, ok = c.Want.(string)
			if !ok {
				t.Fatal("wanted non-error type")
			}
		}
	}

	if c.Got == nil && c.Want == nil {
		return
	}

	if wantErrMsg != "" {
		// compare message
		if strings.Contains(strings.ToLower(c.Got.Error()), strings.ToLower(wantErrMsg)) {
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
