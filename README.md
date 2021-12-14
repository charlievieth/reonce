[![github-actions](https://github.com/charlievieth/reonce/actions/workflows/go.yml/badge.svg)](https://github.com/charlievieth/reonce/actions) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/charlievieth/reonce)

# reonce
Lazily initialized Go regexes.

The `reonce` package provides a lazily initialized wrapper around Go's
[`regexp`](https://golang.org/pkg/regexp) package. This package allows for
regexes to be declared globally without incurring the compilation cost on
program startup/initialization as the regexes are not compiled until first use.
The regexes and compilation are thread-safe
([test](https://github.com/charlievieth/reonce/blob/da431544ab5f2be2359ea6a84c5b40f47aba8bd5/reonce_test.go#L237-L264)).

### Usage

The [`New()`](https://pkg.go.dev/github.com/charlievieth/reonce#New) and
[`NewPOSIX()`](https://pkg.go.dev/github.com/charlievieth/reonce#NewPOSIX)
functions return a lazily initialized
[`*Regexp`](https://pkg.go.dev/github.com/charlievieth/reonce#Regexp) that
wraps a [`*regexp.Regexp`](https://golang.org/pkg/regexp/#Regexp). The
underlying regexp will not be compiled until used and will panic if there is a
compilation error.

```go
import "github.com/charlievieth/reonce"

// Define UUID regexes up-front, neither will be compiled until first use.
var (
	uuidRe = reonce.New(
		`\b[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}\b`,
	)
	uuidExactRe = reonce.New(
		`^[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}$`,
	)
)

// ContainsUUID returns if string s contains a UUID.
func ContainsUUID(s string) bool { return uuidRe.MatchString(s) }

// IsUUID returns if string s is a UUID.
func IsUUID(s string) bool { return uuidExactRe.MatchString(s) }
```

### Testing

The
[`reoncetest`](https://github.com/charlievieth/reonce/blob/master/reoncetest_enabled.go)
build tag can be used to make
[`New()`](https://pkg.go.dev/github.com/charlievieth/reonce#New) and
[`NewPOSIX()`](https://pkg.go.dev/github.com/charlievieth/reonce#NewPOSIX)
immediately compile the regexp and panic with any error. This makes it easy to
check for invalid patterns (especially those created on program initialization)
and is meant for testing only (since it defeats the purpose of this package).

```sh
# use the `reoncetest` build tag when testing your code
$ go test -tags reoncetest ./...
```
