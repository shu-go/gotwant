Package gotwant provides "got: , want: " style test functions.

[![](https://godoc.org/github.com/shu-go/gotwant?status.svg)](https://godoc.org/github.com/shu-go/gotwant)
[![Go Report Card](https://goreportcard.com/badge/github.com/shu-go/gotwant)](https://goreportcard.com/report/github.com/shu-go/gotwant)
![MIT License](https://img.shields.io/badge/License-MIT-blue)

# Usage

## Test

```go
package hoge_test

import "github.com/shu-go/gotwant"

func TestWithGotwant(t *testing.T) {
    v := 1 + 2
    gotwant.Test(t, v, 3) // pass
    gotwant.Test(t, v, "12") // error
    // hoge_test.go:8:
    //     got:  3
    //     want: 12

    v = somefunc()
    gotwant.TestExpr(t, v, v == nil)
    // hoge_test.go:14:
    //     got:  100

    _, e = openFile("hoge_test.go")
    gotwant.Error(t, e, "perm") // pass
    gotwant.Error(t, e, "not found")
    // hoge_test.go:20:
    //     got error:  "permission denied"
    //     want error: "not found"
}
```

## Colorise test output

```
go install github.com/shu-go/gotwant/...
```

```
go test | gotwant
```

The **want** part is colored-diff, showing how the `got` part should be changed. (red should be deleted, green should be inserted)

The got part is as-is.
