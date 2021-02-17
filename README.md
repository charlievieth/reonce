![Go](https://github.com/charlievieth/reonce/workflows/Go/badge.svg?branch=master)

# reonce
Lazily initialized Go regexes

### Usage

The `New` and `NewPOSIX` functions return a lazily initialized `*Regexp` that
wraps a [`*regexp.Regexp`](https://golang.org/pkg/regexp/#Regexp). The
underlying regexp will not be compiled until used and will panic if there is a
compilation error.

```go
// New returns a new lazily initialized Regexp. The underlying *regexp.Regexp
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

The `reoncetest` build tag can be used to make `New()` and `NewPOSIX()` call
MustCompile() on the new Regexp. This makes it easy to check for bad patterns
(especially those created on program initialization) and is meant for testing
only (since it defeats the purpose of this package).

```sh
# use the `reoncetest` build tag when testing your code
go test -tags reoncetest
```
