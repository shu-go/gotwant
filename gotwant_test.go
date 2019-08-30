package gotwant_test

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/shu-go/gotwant"
)

type testerT struct {
	buf bytes.Buffer
}

func (t *testerT) Helper() {
	// nop
}

func (t *testerT) Reset() {
	t.buf.Reset()
}

func (t *testerT) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(&t.buf, format, args...)
}

func TestCase(t *testing.T) {
	tt := &testerT{buf: bytes.Buffer{}}

	c := gotwant.Case("got", "want")
	c.Test(tt)
	r := tt.buf.String()
	if !regexp.MustCompile(`\s*got:  got\s*want: want`).MatchString(r) {
		t.Error(r)
	}

	tt.Reset()
	c = gotwant.Case("want", "want")
	c.Test(tt)
	r = tt.buf.String()
	if r != "" {
		t.Error(r)
	}

	tt.Reset()
	c = gotwant.Case("want", errors.New("want"))
	c.Test(tt)
	r = tt.buf.String()
	if !regexp.MustCompile(`\s*got:  "want"\s*want: &errors.errorString{s:"want"}`).MatchString(r) {
		t.Error(r)
	}
}

func TestExprCase(t *testing.T) {
	tt := &testerT{buf: bytes.Buffer{}}

	c := gotwant.ExprCase("abc", "abc" == "adc")
	c.Test(tt)
	r := tt.buf.String()
	if !regexp.MustCompile(`\s*got:  abc`).MatchString(r) {
		t.Error(r)
	}

	tt.Reset()
	c = gotwant.ExprCase("abc", "abc" == "abc")
	c.Test(tt)
	r = tt.buf.String()
	if r != "" {
		t.Error(r)
	}
}

func TestError(t *testing.T) {
	tt := &testerT{buf: bytes.Buffer{}}

	c := gotwant.Error(errors.New("my error"), "unmatch message")
	c.Test(tt)
	r := tt.buf.String()
	if !regexp.MustCompile(`\s*got error:  my error\s*want error: unmatch message`).MatchString(r) {
		t.Error(r)
	}

	tt.Reset()
	c = gotwant.Error(errors.New("my error"), "my")
	c.Test(tt)
	r = tt.buf.String()
	if r != "" {
		t.Error(r)
	}

	tt.Reset()
	c = gotwant.Error(errors.New("my error"), nil)
	c.Test(tt)
	r = tt.buf.String()
	if !regexp.MustCompile(`\s*got error:  my error`).MatchString(r) {
		t.Error(r)
	}

	tt.Reset()
	c = gotwant.Error(nil, nil)
	c.Test(tt)
	r = tt.buf.String()
	if r != "" {
		t.Error(r)
	}

	tt.Reset()
	c = gotwant.Error(nil, "my error")
	c.Test(tt)
	r = tt.buf.String()
	if !regexp.MustCompile(`\s*want error: my error`).MatchString(r) {
		t.Error(r)
	}
}

func TestPanic(t *testing.T) {
	tt := &testerT{buf: bytes.Buffer{}}

	c := gotwant.Panic(func() { panic("awawawawa...") }, "unmatch message")
	c.Test(tt)
	r := tt.buf.String()
	if !regexp.MustCompile(`\s*got error:  awawawawa...\s*want error: unmatch message`).MatchString(r) {
		t.Error(r)
	}

	tt.Reset()
	c = gotwant.Panic(func() { panic("awawawawa...") }, "awawawawa")
	c.Test(tt)
	r = tt.buf.String()
	if r != "" {
		t.Error(r)
	}

	tt.Reset()
	c = gotwant.Panic(func() { panic("awawawawa...") }, nil)
	c.Test(tt)
	r = tt.buf.String()
	if !regexp.MustCompile(`\s*got error:  awawawawa...`).MatchString(r) {
		t.Error(r)
	}

	tt.Reset()
	c = gotwant.Panic(func() { /*panic("awawawawa...")*/ }, nil)
	c.Test(tt)
	r = tt.buf.String()
	if r != "" {
		t.Error(r)
	}

	tt.Reset()
	c = gotwant.Panic(func() { /*panic("awawawawa...")*/ }, "awawawawa...")
	c.Test(tt)
	r = tt.buf.String()
	if !regexp.MustCompile(`\s*want error: awawawawa...`).MatchString(r) {
		t.Error(r)
	}
}

func TestOption(t *testing.T) {
	tt := &testerT{buf: bytes.Buffer{}}

	c := gotwant.Case("got", "want")
	c.Test(tt)
	r := tt.buf.String()
	if !regexp.MustCompile(`\s*got:  got\s*want: want`).MatchString(r) {
		t.Error(r)
	}

	tt.Reset()
	c = gotwant.Case("got", "want", gotwant.Format("%q"))
	c.Test(tt)
	r = tt.buf.String()
	if !regexp.MustCompile(`\s*got:  "got"\s*want: "want"`).MatchString(r) {
		t.Error(r)
	}

	tt.Reset()
	c = gotwant.Case("got", "want", gotwant.Desc("GOTWANT"))
	c.Test(tt)
	r = tt.buf.String()
	if !regexp.MustCompile(`GOTWANT\s*got:  got\s*want: want`).MatchString(r) {
		t.Error(r)
	}
}

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
		gotwant.TestError(t, errfunc2(), "")
		gotwant.TestError(t, errfunc2(), myErr)
		gotwant.TestError(t, errfunc2(), "is an")
	})

	t.Run("Panic", func(t *testing.T) {
		gotwant.TestPanic(t, func() {
			panic("hoge")
		}, "og")

		gotwant.TestPanic(t, func() {
			panic("hoge")
		}, "")

		gotwant.TestPanic(t, func() {
		}, nil)

		gotwant.TestPanic(t, func() {
			panic(123)
		}, 123)
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
