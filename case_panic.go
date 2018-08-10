package gotwant

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type panicCase struct {
	Got  func()      // what you got.
	Want interface{} // an panic message or nil(not panicked)

	Fmt  string // used in t.Errorf displaying got and want.  default: FmtDefault
	Desc string // a line description
}

// Case constructs a value-comaration test case.
func Panic(got func(), want interface{}, opts ...Option) *panicCase {
	c := &panicCase{
		Got:  got,
		Want: want,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *panicCase) SetFmt(format string) {
	c.Fmt = format
}

func (c *panicCase) SetDesc(desc string) {
	c.Desc = desc
}

func (c *panicCase) Test(t *testing.T) {
	t.Helper()

	valfmt := c.Fmt
	if valfmt == "" {
		valfmt = FmtDefault
	}

	var gotErr interface{} = nil
	func() {
		defer func() {
			if err := recover(); err != nil {
				gotErr = err
			}
		}()

		c.Got()
	}()

	if gotErr == nil && c.Want == nil {
		return
	}

	if gotErr == nil {
		errfmt := fmt.Sprintf("%s\ngot NO panic.\nwant error: %s", c.Desc, valfmt)
		t.Errorf(errfmt, c.Want)
	}

	var wantErrMsg *string
	if c.Want != nil {
		if err, ok := c.Want.(error); ok {
			str := err.Error()
			wantErrMsg = &str
		} else if str, ok := c.Want.(string); ok {
			wantErrMsg = &str
		} else if s, ok := c.Want.(fmt.Stringer); ok {
			str := s.String()
			wantErrMsg = &str
		}
	}

	if wantErrMsg != nil {
		// compare message
		if s, ok := gotErr.(string); ok {
			if strings.Contains(strings.ToLower(s), strings.ToLower(*wantErrMsg)) {
				return
			}
		} else if s, ok := gotErr.(fmt.Stringer); ok {
			if strings.Contains(strings.ToLower(s.String()), strings.ToLower(*wantErrMsg)) {
				return
			}
		}

	}

	if reflect.DeepEqual(gotErr, c.Want) {
		return
	}

	errfmt := fmt.Sprintf("%s\ngot error:  %s\nwant error: %s", c.Desc, valfmt, valfmt)
	t.Errorf(errfmt, gotErr, c.Want)
}
