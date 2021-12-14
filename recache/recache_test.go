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

func testCacheSize(t *testing.T, fn func(int) *Cache) {
	c := fn(8)
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

func TestCacheSize(t *testing.T) {
	testCacheSize(t, New)
}

func TestCacheSizePOSIX(t *testing.T) {
	testCacheSize(t, NewPOSIX)
}

func mustPanic(t *testing.T, msg string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("test did not panic:", msg)
		}
	}()
	fn()
}

func TestNegativeCacheSize(t *testing.T) {
	mustPanic(t, "New(-1)", func() { New(-1) })
	mustPanic(t, "NewPOSIX(-1)", func() { NewPOSIX(-1) })
}

func testUnlimitedCacheSize(t *testing.T, fn func(int) *Cache) {
	c := fn(0)
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

func TestUnlimitedCacheSize(t *testing.T) {
	testUnlimitedCacheSize(t, New)
}

func TestUnlimitedCacheSizePOSIX(t *testing.T) {
	testUnlimitedCacheSize(t, NewPOSIX)
}

func TestCompile(t *testing.T) {
	// TODO: this uses the global caches - consider changing
	if _, err := Compile(`a`); err != nil {
		t.Fatal(err)
	}
	if _, err := CompilePOSIX(`a`); err != nil {
		t.Fatal(err)
	}

	// make sure an error is always returned
	for i := 0; i < 3; i++ {
		if _, err := Compile(`[a`); err == nil {
			t.Errorf("Compile(`[a`): expected error got: %v", err)
		}
	}
	for i := 0; i < 3; i++ {
		if _, err := CompilePOSIX(`[a`); err == nil {
			t.Errorf("Compile(`[a`): expected error got: %v", err)
		}
	}
}

// make sure we're thread-safe
func testCompileParallel(t *testing.T, format string, fn func(int) *Cache) {
	if testing.Short() {
		t.Skip("short test")
	}

	var wg sync.WaitGroup
	c := fn(124)
	n := new(int64)
	start := make(chan struct{})
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for i := 1; i <= 4096; i++ {
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

// make sure we're thread-safe
func TestCompileParallel(t *testing.T) {
	const format = `(?P<foo>.*)(?P<bar>(a)b(%d)?)(?P<foo>.*)a`
	testCompileParallel(t, format, New)
}

func TestCompileParallelPOSIX(t *testing.T) {
	const format = `.{1}{2}.{3}.%d`
	testCompileParallel(t, format, NewPOSIX)
}

func TestMustCompile(t *testing.T) {
	mustPanic(t, "MustCompile: `[a`", func() { MustCompile(`[a`) })

	re := MustCompile("^abc$")
	if !re.MatchString("abc") {
		t.Fatal("failed to match string")
	}
}

func TestMustCompilePOSIX(t *testing.T) {
	mustPanic(t, "MustCompilePOSIX: `[a`", func() { MustCompilePOSIX(`[a`) })

	re := MustCompilePOSIX("^abc$")
	if !re.MatchString("abc") {
		t.Fatal("failed to match string")
	}
}

func testSetMaxEntries(t *testing.T, fn func(int) *Cache) {
	c := fn(15)
	for i := 'a'; i < 'a'+20; i++ {
		c.MustCompile(string(i))
	}
	if c.Len() != 15 {
		t.Fatalf("Len: got: %d want: %d", 15, c.Len())
	}
	// Make sure the expected entries are there
	for i := 'a' + 5; i < 'a'+20; i++ {
		key := string(i)
		ee, ok := c.cache[key]
		if !ok {
			t.Errorf("Evicted key: %q", key)
		}
		exp := ee.re.String()
		got := c.MustCompile(key).String()
		if exp != got {
			t.Errorf("MustCompile(`%s`) = %q want: %q", key, got, exp)
		}
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

func TestSetMaxEntries(t *testing.T) {
	testSetMaxEntries(t, New)

	// Test global cache
	t.Run("Global", func(t *testing.T) {
		testSetMaxEntries(t, func(n int) *Cache {
			for Len() != 0 {
				std.removeOldest()
			}
			SetMaxEntries(n)
			if MaxEntries() != n {
				t.Fatalf("MaxEntries: got: %d want: %d", MaxEntries(), n)
			}
			return std
		})
	})
}

func TestSetMaxEntriesPOSIX(t *testing.T) {
	testSetMaxEntries(t, NewPOSIX)

	// Test global cache
	t.Run("Global", func(t *testing.T) {
		testSetMaxEntries(t, func(n int) *Cache {
			for LenPOSIX() != 0 {
				posix.removeOldest()
			}
			SetMaxEntriesPOSIX(n)
			if MaxEntriesPOSIX() != n {
				t.Fatalf("MaxEntriesPOSIX: got: %d want: %d", MaxEntriesPOSIX(), n)
			}
			return posix
		})
	})
}

func TestGlobals(t *testing.T) {
	if std.POSIX() {
		t.Error("std.POSIX() should be false")
	}
	if !posix.POSIX() {
		t.Error("posix.POSIX() should be true")
	}
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
