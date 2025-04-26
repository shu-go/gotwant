// Package gotwant provides "got: , want: " style test functions.
package gotwant

import (
	"fmt"
	"strings"
)

var (
	// FmtDefault is a default value of displaying contents of got/want.
	FmtDefault = "%v"
)

// T has a few part of testing.T, to test gotwant itself.
type T interface {
	Helper()
	Errorf(format string, args ...interface{})
}

// TestCase is made by calling Case, Error, ...
type TestCase interface {
	Test(t T)

	SetFmt(string)
	SetDesc(string)
}

// Option sets optional values.
type Option func(TestCase)

// Format sets(changes) print format. default: FmtDefault
func Format(fmt string) Option {
	return func(c TestCase) {
		c.SetFmt(fmt)
	}
}

// Desc sets a line description.
func Desc(desc string) Option {
	return func(c TestCase) {
		c.SetDesc(desc)
	}
}

// Test if for a single try.
func Test(t T, got, want interface{}, opts ...Option) {
	t.Helper()

	Case(got, want, opts...).Test(t)
}

// TestExpr tests got == expr (boolean comparison)
func TestExpr(t T, got interface{}, expr bool, opts ...Option) {
	t.Helper()

	ExprCase(got, expr, opts...).Test(t)
}

// TestError tests given error (got) is (1) exactly the error you wanted or (2) its message matches your pattern.
// If want is a string, this func tests with strings.Contains(got, want)
// else, this func tests with reflect.DeepEqual(got, want)
func TestError(t T, got error, want interface{}, opts ...Option) {
	t.Helper()

	Error(got, want, opts...).Test(t)
}

// TestPanic tests function got panic(want) or not.
// Pass nil to `want` if `got` is not expected panicked.
func TestPanic(t T, got func(), want interface{}, opts ...Option) {
	t.Helper()

	Panic(got, want, opts...).Test(t)
}

// TestAll is for a series of tries(Cases).
func TestAll(t T, cases []TestCase) {
	t.Helper()

	for _, c := range cases {
		c.Test(t)
	}
}

func stringify(s interface{}) *string {
	if s == nil {
		return nil
	}

	if err, ok := s.(error); ok {
		ss := err.Error()
		return &ss
	} else if str, ok := s.(string); ok {
		return &str
	} else if ss, ok := s.(fmt.Stringer); ok {
		sss := ss.String()
		return &sss
	}
	return nil
}

func indent(s string) string {
	return strings.ReplaceAll(s, "\n", "\n      ")
}
