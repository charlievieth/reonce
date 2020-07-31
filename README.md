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
