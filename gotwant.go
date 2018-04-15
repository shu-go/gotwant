// Package gotwant provides "got: , want: " style test functions.
package gotwant

import (
	"testing"
)

var (
	FmtDefault string = "%v"
)

type TestCase interface {
	Test(t *testing.T)

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
func Test(t *testing.T, got, want interface{}, opts ...Option) {
	t.Helper()

	Case(got, want, opts...).Test(t)
}

func TestExpr(t *testing.T, got interface{}, expr bool, opts ...Option) {
	t.Helper()

	ExprCase(got, expr, opts...).Test(t)
}

// TestError tests given error (got) is (1) exactly the error you wanted or (2) its message matches your pattern.
// If want is a string, this func tests with strings.Contains(got, want)
// else, this func tests with reflect.DeepEqual(got, want)
func TestError(t *testing.T, got error, want interface{}, opts ...Option) {
	t.Helper()

	Error(got, want, opts...).Test(t)
}

// TestAll is for a series of tries(Cases).
func TestAll(t *testing.T, cases []TestCase) {
	t.Helper()

	for _, c := range cases {
		c.Test(t)
	}
}
