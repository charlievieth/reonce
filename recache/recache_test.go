package recache

import (
	"fmt"
	"regexp"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestMustCompileOnce(t *testing.T) {
	var patterns sync.Map
	defer func() { compile = regexp.Compile }()
	compile = func(str string) (*regexp.Regexp, error) {
		v, _ := patterns.LoadOrStore(str, new(int64))
		i := v.(*int64)
		atomic.AddInt64(i, 1)
		return regexp.MustCompile(str), nil
	}

	c := New(4096)

	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	ch := make(chan string, numCPU)
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for s := range ch {
				c.MustCompile(s)
			}
		}()
	}

	n := 1024 * 32
	if testing.Short() {
		n = 2048
	}
	for i := 0; i < n; i++ {
		s := fmt.Sprintf(`(?P<foo>.*)(?P<bar>(a)b(%d)?)(?P<foo>.*)a`, i)
		for j := 0; j < numCPU; j++ {
			ch <- s
		}
	}
	close(ch)
	wg.Wait()

	patterns.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(*int64)
		if *v > 1 {
			t.Errorf("%s: %d\n", k, *v)
		}
		return true
	})
}

func TestCompile(t *testing.T) {
	if _, err := Compile(`[a`); err == nil {
		t.Errorf("Compile(`[a`): expected error got: %v", err)
	}
	if _, err := CompilePOSIX(`[a`); err == nil {
		t.Errorf("CompilePOSIX(`[a`): expected error got: %v", err)
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

func TestMustCompilePOSIX(t *testing.T) {
	mustPanic(t, "MustCompilePOSIX: `[a`", func() { MustCompilePOSIX(`[a`) })

	re := MustCompilePOSIX("^abc$")
	if !re.MatchString("abc") {
		t.Fatal("failed to match string")
	}
}

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
}

const benchRe = "[\\pL\\pN][^\\pL\\pN]|(.$)"

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
