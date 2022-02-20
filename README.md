[![github-actions](https://github.com/charlievieth/reonce/actions/workflows/go.yml/badge.svg)](https://github.com/charlievieth/reonce/actions) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/charlievieth/reonce)

# reonce

Package reonce is a thin wrapper over [`regexp`](https://golang.org/pkg/regexp),
allowing the use of global regexp variables without having to compile them at
program initialization.

Lazy compilation is thread-safe and will panic if there is an error.
This matches the behavior of [`regexp.MustCompile()`](https://pkg.go.dev/regexp#MustCompile)
and [`regexp.MustCompilePOSIX()`](https://pkg.go.dev/regexp#MustCompilePOSIX).
The panic messages are the same as those used by the standard library.

### Usage

The [`New()`](https://pkg.go.dev/github.com/charlievieth/reonce#New) and
[`NewPOSIX()`](https://pkg.go.dev/github.com/charlievieth/reonce#NewPOSIX)
functions return a lazily initialized
[`*Regexp`](https://pkg.go.dev/github.com/charlievieth/reonce#Regexp) that
thinly wraps [`*regexp.Regexp`](https://golang.org/pkg/regexp/#Regexp).
The underlying regexp will not be compiled until used and will panic if
there is a compilation error.

The [`Regexp.Compile()`](https://pkg.go.dev/github.com/charlievieth/reonce#Regexp.Compile)
method can be used to manually force compilation or get the compilation error
of a previously compiled `Regexp`, if any.

```go
import "github.com/charlievieth/reonce"

// Define UUID regexes up-front, neither will be compiled until first use.

var uuidRe = reonce.New(
	`\b[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}\b`,
)

var uuidExactRe = reonce.New(
	`^[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}$`,
)

// ContainsUUID returns if string s contains a UUID.
func ContainsUUID(s string) bool {
	return uuidRe.MatchString(s)
}

// IsUUID returns if string s is a UUID.
func IsUUID(s string) bool {
	return uuidExactRe.MatchString(s)
}
```

### Testing

The `reoncetest` build tag forces `reonce` to immediately compile regexes and
panic on error. This useful for testing regexes that are initialized on program
creation or any other `*Regexp` that is created, but not immediately used.

```sh
# use the `reoncetest` build tag when testing your code
$ go test -tags reoncetest ./...
```

This build tah should not be used in production/releases as it disables lazy
compilation, which is the purpose of this package.

### Overhead

Once compiled, the overhead of lazy compilation is a call to
[`sync.Once.Do()`](https://pkg.go.dev/sync#Once) which should be around \~2ns.

```
goos: darwin
goarch: arm64
pkg: github.com/charlievieth/reonce
BenchmarkInitOverhead-10                      	557014276	         2.173 ns/op
BenchmarkInitOverhead_Parallel-10             	1000000000	         0.2770 ns/op
```

```
goos: linux
goarch: amd64
pkg: github.com/charlievieth/reonce
cpu: Intel(R) Core(TM) i9-9900K CPU @ 3.60GHz
BenchmarkInitOverhead-16                      	710212690	         1.624 ns/op
BenchmarkInitOverhead_Parallel-16             	1000000000	         0.3778 ns/op
```
