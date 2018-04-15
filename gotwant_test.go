package gotwant_test

import (
	"fmt"
	"testing"

	"bitbucket.org/shu/gotwant"
)

func TestHowToWrite(t *testing.T) {
	t.Run("Comparation", func(t *testing.T) {
		gotwant.Test(t, 1, 1)
		gotwant.Test(t, "1", "1")
		gotwant.Test(t, struct{ A string }{A: "aaa"}, struct{ A string }{A: "aaa"})
	})

	t.Run("Expr", func(t *testing.T) {
		i := 1
		gotwant.TestExpr(t, i, i > 0)
		i++
		gotwant.TestExpr(t, i, i == 2)
	})

	t.Run("Error", func(t *testing.T) {
		errfunc1 := func() error {
			return nil
		}
		myErr := fmt.Errorf("this is an error")
		errfunc2 := func() error {
			return myErr
		}

		gotwant.TestError(t, errfunc1(), nil)
		gotwant.TestError(t, errfunc2(), myErr)
		gotwant.TestError(t, errfunc2(), "is an")
	})

	t.Run("Options", func(t *testing.T) {
		gotwant.Test(t, 1, 1, gotwant.Desc("one"))
		gotwant.Test(t, "1", "1", gotwant.Desc("one as string"), gotwant.Format("%q"))
		gotwant.Test(t, struct{ A string }{A: "aaa"}, struct{ A string }{A: "aaa"}, gotwant.Format("%#v"))
	})

	t.Run("Table", func(t *testing.T) {
		table := []gotwant.TestCase{
			gotwant.Case(1, 1),
			gotwant.Case("1", "1"),
			gotwant.Case(struct{ A string }{A: "aaa"}, struct{ A string }{A: "aaa"}),
		}
		gotwant.TestAll(t, table)
	})
}
