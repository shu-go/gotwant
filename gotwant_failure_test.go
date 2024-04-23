//go:build failure

package gotwant_test

import (
	"testing"

	"github.com/shu-go/gotwant"
)

func TestFailure1(t *testing.T) {
	gotwant.Test(t, "got", "want")
}

func TestFailure2(t *testing.T) {
	gotwant.Test(t, `aaa
bbb
ccc`, `aaa
bab
ddd`)
}

func TestFailure3(t *testing.T) {
	gotwant.Test(t, "got", "want")
	gotwant.Test(t, `aaa
bbb
ccc`, `aaa
bab
ddd`)
	gotwant.Test(t, `aaa
bbb
ccc
ddd`, `aaaa
bb
ddd`)
}

func TestFailureSub(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		gotwant.Test(t, "got", "want")
	})

	t.Run("2", func(t *testing.T) {
		gotwant.Test(t, `aaa
bbb
ccc`, `aaa
bab
ddd`)
	})

	t.Run("3", func(t *testing.T) {
		gotwant.Test(t, "got", "want")
		gotwant.Test(t, `aaa
bbb
ccc`, `aaa
bab
ddd`)
		gotwant.Test(t, `aaa
bbb
ccc
ddd`, `aaaa
bb
ddd`)
	})
}
