package recache

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

func TestDefaultCacheValues(t *testing.T) {
	var c Cache
	c.MustCompile("a")

	if c.cache == nil {
		t.Error("nil cache")
	}
	if c.ll == nil {
		t.Error("nil list")
	}
	if c.POSIX() {
		t.Errorf("posix: got: %t want: %t", c.POSIX(), false)
	}
	// TODO(charlie): do we really want the defualt to be unlimited?
	if c.MaxEntries() != 0 {
		t.Errorf("MaxEntries: got: %d want: %d", c.MaxEntries(), 0)
	}
}

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

func TestCacheSize(t *testing.T) {
	c := New(8)
	for i := 0; i < 16; i++ {
		c.MustCompile(strconv.Itoa(i))
	}
	if c.Len() != 8 {
		t.Errorf("Len: got: %d want: %d", c.Len(), 8)
	}
	e := c.ll.Front()
	for i := 15; i >= 8; i-- {
		exp := strconv.Itoa(i)
		if e.re.String() != exp {
			t.Errorf("%d: got: %s want: %s", i, quote(e.re.String()), exp)
		}
		e = e.next
	}

	c.SetMaxEntries(4)
	if c.Len() != 4 {
		t.Errorf("Len: got: %d want: %d", c.Len(), 4)
	}
	e = c.ll.Front()
	for i := 15; i >= 12; i-- {
		exp := strconv.Itoa(i)
		if e.re.String() != exp {
			t.Errorf("%d: got: %s want: %s", i, quote(e.re.String()), exp)
		}
		e = e.next
	}
}

func TestNegativeCacheSize(t *testing.T) {
	mustPanic(t, "New(-1)", func() { New(-1) })
}

func TestUnlimitedCacheSize(t *testing.T) {
	c := New(0)
	for i := 0; i < 1024; i++ {
		c.MustCompile(strconv.Itoa(i))
	}
	if c.Len() != 1024 {
		t.Errorf("Len: got: %d want: %d", c.Len(), 1024)
	}
	c.SetMaxEntries(32)
	if c.Len() != 32 {
		t.Errorf("Len: got: %d want: %d", c.Len(), 32)
	}
}

func TestCompile(t *testing.T) {
	// make sure an error is always returned
	for i := 0; i < 3; i++ {
		if _, err := Compile(`[a`); err == nil {
			t.Errorf("Compile(`[a`): expected error got: %v", err)
		}
	}
}

// make sure we're thread-safe
func TestCompileParallel(t *testing.T) {
	if testing.Short() {
		t.Skip("short test")
	}
	const format = `(?P<foo>.*)(?P<bar>(a)b(%d)?)(?P<foo>.*)a`

	var wg sync.WaitGroup
	c := New(124)
	n := new(int64)
	start := make(chan struct{})
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for i := 0; i < 4096; i++ {
				expr := fmt.Sprintf(format, atomic.AddInt64(n, 1))
				c.MustCompile(expr)
			}
		}()
	}
	close(start)
	wg.Wait()

	// Make sure that roughly the right values are cached. We can't be
	// exact here because the scheduling of goroutines impacts the order
	// of entries.
	e := c.ll.Front()
	misses := 0
	for i := *n; i > *n-int64(c.MaxEntries()); i-- {
		expr := fmt.Sprintf(format, i)
		if _, ok := c.cache[expr]; !ok {
			t.Logf("missing: %s", expr)
			misses++
		}
		e = e.next
	}
	if misses > c.Len()/10 {
		t.Errorf("missing: %d", misses)
	}
}

func mustPanic(t *testing.T, msg string, fn func()) {
	t.Helper()
	defer func() {
		if e := recover(); e == nil {
			t.Fatal("test did not panic:", msg)
		}
	}()
	fn()
}

func TestMustCompile(t *testing.T) {
	mustPanic(t, "MustCompile: `[a`", func() { MustCompile(`[a`) })

	re := MustCompile("^abc$")
	if !re.MatchString("abc") {
		t.Fatal("failed to match string")
	}
}

// func TestMustCompilePOSIX(t *testing.T) {
// 	mustPanic(t, "MustCompilePOSIX: `[a`", func() { MustCompilePOSIX(`[a`) })
//
// 	re := MustCompilePOSIX("^abc$")
// 	if !re.MatchString("abc") {
// 		t.Fatal("failed to match string")
// 	}
// }

func TestSetMaxEntries(t *testing.T) {
	c := New(15)
	for i := 'a'; i < 'a'+20; i++ {
		c.MustCompile(string(i))
	}
	if c.Len() != 15 {
		t.Fatalf("Len: got: %d want: %d", 15, c.Len())
	}
	c.SetMaxEntries(10)
	if c.Len() != 10 {
		t.Fatalf("Len: got: %d want: %d", 10, c.Len())
	}
	if c.MaxEntries() != 10 {
		t.Fatalf("MaxEntries: got: %d want: %d", 10, c.MaxEntries())
	}
	c.SetMaxEntries(0)
	if c.Len() != 10 {
		t.Fatalf("Len: got: %d want: %d", 10, c.Len())
	}
	if c.MaxEntries() != 0 {
		t.Fatalf("MaxEntries: got: %d want: %d", 0, c.MaxEntries())
	}
	mustPanic(t, "SetMaxEntries(-1)", func() {
		c.SetMaxEntries(-1)
	})
}

const benchRe = `[\pL\pN][^\pL\pN]|(.$)`

func TestBenchmarkRegex(t *testing.T) {
	c := New(0)
	c.MustCompile(benchRe)
}

func BenchmarkCompile_Baseline(b *testing.B) {
	b.Skip("skip")
	for i := 0; i < b.N; i++ {
		_ = regexp.MustCompile(benchRe)
	}
}

func BenchmarkCompile(b *testing.B) {
	c := New(100)
	for i := 0; i < b.N; i++ {
		_ = c.MustCompile(benchRe)
	}
}

func BenchmarkCompile_Baseline_Parallel(b *testing.B) {
	b.Skip("skip")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = regexp.MustCompile(benchRe)
		}
	})
}

func BenchmarkCompile_Parallel(b *testing.B) {
	c := New(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = c.MustCompile(benchRe)
		}
	})
}
