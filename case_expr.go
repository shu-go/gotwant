package gotwant

import (
	"fmt"
)

// ExprCase constructs a test case of given expr.
func ExprCase(got interface{}, expr bool, opts ...Option) *exprCase {
	c := &exprCase{
		Got:  got,
		Expr: expr,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

type exprCase struct {
	Got  interface{} // what you got.
	Expr bool        // what you expected with Got.

	Desc string // a line description
}

func (c *exprCase) SetFmt(format string) {
}

func (c *exprCase) SetDesc(desc string) {
	c.Desc = desc
}

func (c *exprCase) Test(t T) {
	t.Helper()

	if !c.Expr {
		valfmt := FmtDefault
		errfmt := fmt.Sprintf("%s\ngot:  %s", c.Desc, valfmt)
		t.Errorf(errfmt, c.Got)
	}
}
