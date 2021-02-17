![Go](https://github.com/charlievieth/reonce/workflows/Go/badge.svg?branch=master) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/charlievieth/reonce)

# reonce
Lazily initialized Go regexes.

The `reonce` package provides a lazily initialized wrapper around Go's [`regexp`](https://golang.org/pkg/regexp) package. This package allows for regexes to be declared globally without incurring the compilation cost on program startup/initialization as the regexes are not compiled until first use. The regexes and compilation are thread-safe ([test](https://github.com/charlievieth/reonce/blob/6f5299ea34e785e202421258a42c336bd6a9a02f/reonce_test.go#L177-L204)).

### Usage

The [`New()`](https://pkg.go.dev/github.com/charlievieth/reonce#New) and [`NewPOSIX()`](https://pkg.go.dev/github.com/charlievieth/reonce#NewPOSIX) functions return a lazily initialized [`*Regexp`](https://pkg.go.dev/github.com/charlievieth/reonce#Regexp) that
wraps a [`*regexp.Regexp`](https://golang.org/pkg/regexp/#Regexp). The
underlying regexp will not be compiled until used and will panic if there is a
compilation error.

```go
// New returns a new lazily initialized *Regexp. The underlying *regexp.Regexp
// will be compiled on first use. If pattern expr is invalid it will panic.
func New(expr string) *Regexp {
	return &Regexp{expr: expr}
}

// New returns a new lazily initialized POSIX Regexp.
func NewPOSIX(expr string) *Regexp {
	return &Regexp{expr: expr, posix: true}
}
```

### Testing

The [`reoncetest`](https://github.com/charlievieth/reonce/blob/5faff1a5ae70387f6f4a73320b726063a7834fb8/reoncetest_enabled.go#L1) build tag can be used to make `New()` and `NewPOSIX()` call
[`MustCompile()`](https://golang.org/pkg/regexp/#MustCompile) on the new Regexp. This makes it easy to check for bad patterns
(especially those created on program initialization) and is meant for testing
only (since it defeats the purpose of this package).

```sh
# use the `reoncetest` build tag when testing your code
$ go test -tags reoncetest ./...
```
