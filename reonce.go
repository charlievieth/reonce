package reonce

import (
	"io"
	"regexp"
	"sync"
)

type Regexp struct {
	once  sync.Once
	re    *regexp.Regexp
	expr  string
	posix bool
}

// New returns a new lazily initialized Regexp. The underlying *regexp.Regexp
// will be compiled on first use. If pattern expr is invalid it will panic.
func New(expr string) *Regexp {
	return &Regexp{expr: expr}
}

// New returns a new lazily initialized POSIX Regexp.
func NewPOSIX(expr string) *Regexp {
	return &Regexp{expr: expr, posix: true}
}

// Compile manually compiles the Regexp and returns the error, this is a no-op
// if the Regexp was already lazily compiled by a call to any of it's methods.
func (re *Regexp) Compile() (err error) {
	re.once.Do(func() {
		if re.posix {
			re.re, err = regexp.CompilePOSIX(re.expr)
		} else {
			re.re, err = regexp.Compile(re.expr)
		}
	})
	return err
}

func (re *Regexp) mustCompile() {
	if re.posix {
		re.re = regexp.MustCompilePOSIX(re.expr)
	} else {
		re.re = regexp.MustCompile(re.expr)
	}
}

// MustCompile compiles the Regexp and panics if there is an error, this is a
// no-op if the Regexp was already lazily compiled by a call to any of  it's
// methods.
func (re *Regexp) MustCompile() { re.once.Do(re.mustCompile) }

// Regexp returns the underlying *regexp.Regexp.
func (re *Regexp) Regexp() *regexp.Regexp {
	re.once.Do(re.mustCompile)
	return re.re
}

func (re *Regexp) Copy() *regexp.Regexp {
	re.once.Do(re.mustCompile)
	return re.re.Copy()
}

func (re *Regexp) Expand(dst []byte, template []byte, src []byte, match []int) []byte {
	re.once.Do(re.mustCompile)
	return re.re.Expand(dst, template, src, match)
}

func (re *Regexp) ExpandString(dst []byte, template string, src string, match []int) []byte {
	re.once.Do(re.mustCompile)
	return re.re.ExpandString(dst, template, src, match)
}

func (re *Regexp) Find(b []byte) []byte {
	re.once.Do(re.mustCompile)
	return re.re.Find(b)
}

func (re *Regexp) FindAll(b []byte, n int) [][]byte {
	re.once.Do(re.mustCompile)
	return re.re.FindAll(b, n)
}

func (re *Regexp) FindAllIndex(b []byte, n int) [][]int {
	re.once.Do(re.mustCompile)
	return re.re.FindAllIndex(b, n)
}

func (re *Regexp) FindAllString(s string, n int) []string {
	re.once.Do(re.mustCompile)
	return re.re.FindAllString(s, n)
}

func (re *Regexp) FindAllStringIndex(s string, n int) [][]int {
	re.once.Do(re.mustCompile)
	return re.re.FindAllStringIndex(s, n)
}

func (re *Regexp) FindAllStringSubmatch(s string, n int) [][]string {
	re.once.Do(re.mustCompile)
	return re.re.FindAllStringSubmatch(s, n)
}

func (re *Regexp) FindAllStringSubmatchIndex(s string, n int) [][]int {
	re.once.Do(re.mustCompile)
	return re.re.FindAllStringSubmatchIndex(s, n)
}

func (re *Regexp) FindAllSubmatch(b []byte, n int) [][][]byte {
	re.once.Do(re.mustCompile)
	return re.re.FindAllSubmatch(b, n)
}

func (re *Regexp) FindAllSubmatchIndex(b []byte, n int) [][]int {
	re.once.Do(re.mustCompile)
	return re.re.FindAllSubmatchIndex(b, n)
}

func (re *Regexp) FindIndex(b []byte) (loc []int) {
	re.once.Do(re.mustCompile)
	return re.re.FindIndex(b)
}

func (re *Regexp) FindReaderIndex(r io.RuneReader) (loc []int) {
	re.once.Do(re.mustCompile)
	return re.re.FindReaderIndex(r)
}

func (re *Regexp) FindReaderSubmatchIndex(r io.RuneReader) []int {
	re.once.Do(re.mustCompile)
	return re.re.FindReaderSubmatchIndex(r)
}

func (re *Regexp) FindString(s string) string {
	re.once.Do(re.mustCompile)
	return re.re.FindString(s)
}

func (re *Regexp) FindStringIndex(s string) (loc []int) {
	re.once.Do(re.mustCompile)
	return re.re.FindStringIndex(s)
}

func (re *Regexp) FindStringSubmatch(s string) []string {
	re.once.Do(re.mustCompile)
	return re.re.FindStringSubmatch(s)
}

func (re *Regexp) FindStringSubmatchIndex(s string) []int {
	re.once.Do(re.mustCompile)
	return re.re.FindStringSubmatchIndex(s)
}

func (re *Regexp) FindSubmatch(b []byte) [][]byte {
	re.once.Do(re.mustCompile)
	return re.re.FindSubmatch(b)
}

func (re *Regexp) FindSubmatchIndex(b []byte) []int {
	re.once.Do(re.mustCompile)
	return re.re.FindSubmatchIndex(b)
}

func (re *Regexp) LiteralPrefix() (prefix string, complete bool) {
	re.once.Do(re.mustCompile)
	return re.re.LiteralPrefix()
}

func (re *Regexp) Longest() {
	re.once.Do(re.mustCompile)
	re.re.Longest()
}

func (re *Regexp) Match(b []byte) bool {
	re.once.Do(re.mustCompile)
	return re.re.Match(b)
}

func (re *Regexp) MatchReader(r io.RuneReader) bool {
	re.once.Do(re.mustCompile)
	return re.re.MatchReader(r)
}

func (re *Regexp) MatchString(s string) bool {
	re.once.Do(re.mustCompile)
	return re.re.MatchString(s)
}

func (re *Regexp) NumSubexp() int {
	re.once.Do(re.mustCompile)
	return re.re.NumSubexp()
}

func (re *Regexp) ReplaceAll(src, repl []byte) []byte {
	re.once.Do(re.mustCompile)
	return re.re.ReplaceAll(src, repl)
}

func (re *Regexp) ReplaceAllFunc(src []byte, repl func([]byte) []byte) []byte {
	re.once.Do(re.mustCompile)
	return re.re.ReplaceAllFunc(src, repl)
}

func (re *Regexp) ReplaceAllLiteral(src, repl []byte) []byte {
	re.once.Do(re.mustCompile)
	return re.re.ReplaceAllLiteral(src, repl)
}

func (re *Regexp) ReplaceAllLiteralString(src, repl string) string {
	re.once.Do(re.mustCompile)
	return re.re.ReplaceAllLiteralString(src, repl)
}

func (re *Regexp) ReplaceAllString(src, repl string) string {
	re.once.Do(re.mustCompile)
	return re.re.ReplaceAllString(src, repl)
}

func (re *Regexp) ReplaceAllStringFunc(src string, repl func(string) string) string {
	re.once.Do(re.mustCompile)
	return re.re.ReplaceAllStringFunc(src, repl)
}

func (re *Regexp) Split(s string, n int) []string {
	re.once.Do(re.mustCompile)
	return re.re.Split(s, n)
}

func (re *Regexp) String() string {
	re.once.Do(re.mustCompile)
	return re.re.String()
}

func (re *Regexp) SubexpNames() []string {
	re.once.Do(re.mustCompile)
	return re.re.SubexpNames()
}
