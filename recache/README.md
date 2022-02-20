[![github-actions](https://github.com/charlievieth/reonce/actions/workflows/go.yml/badge.svg)](https://github.com/charlievieth/reonce/actions) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/charlievieth/reonce/recache)

# recache

Package `recache` implements a thread-safe LRU cache of
[`reonce.Regexp`](https://pkg.go.dev/github.com/charlievieth/reonce#Regexp)
regular expressions.

### Usage

The `reonce` package provides two global regexp caches (POSIX non-POSIX) that
can be used with the top-level
[`Compile`](https://pkg.go.dev/github.com/charlievieth/reonce/recache#Compile),
[`MustCompile`](https://pkg.go.dev/github.com/charlievieth/reonce/recache#MustCompile),
[`MaxEntries`](https://pkg.go.dev/github.com/charlievieth/reonce/recache#MaxEntries),
[`SetMaxEntries`](https://pkg.go.dev/github.com/charlievieth/reonce/recache#SetMaxEntries),
and [`Len`](https://pkg.go.dev/github.com/charlievieth/reonce/recache#Len)
functions (and their `*POSIX` counterparts). Caches can also be created with
[`New`](https://pkg.go.dev/github.com/charlievieth/reonce/recache#New) and
[`NewPOSIX`](https://pkg.go.dev/github.com/charlievieth/reonce/recache#NewPOSIX).

