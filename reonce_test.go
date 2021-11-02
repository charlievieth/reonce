//go:build !reoncetest
// +build !reoncetest

package reonce

import (
	"bytes"
	"io"
	"math/rand"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"testing/quick"
	"time"
)

func TestBadPattern(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected a panic")
		}
	}()
	New("*").MustCompile()
}

func TestCompile(t *testing.T) {
	re := New("a")
	err := re.Compile()
	if err != nil {
		t.Fatal(err)
	}
	// make sure we used the sync.Once
	re.once.Do(func() {
		t.Fatal("Compile() did not use the sync.Once")
	})
}

func TestCompileError(t *testing.T) {
	const BadPattern = "*"
	const ErrMessage = "missing argument to repetition operator"
	re := New(BadPattern)
	err := re.Compile()
	if err == nil || !strings.Contains(err.Error(), ErrMessage) {
		t.Errorf("Compile: expected error to contain: %q got: %v", ErrMessage, err)
	}

	// make sure we still return the error
	err2 := re.Compile()
	if err2 != err {
		t.Error("Second call to Compile() did not return the same error")
	}
}

func TestMustCompileError(t *testing.T) {
	mustPanic := func(t *testing.T, fn func()) (msg string) {
		defer func() {
			e := recover()
			if e == nil {
				t.Error("no panic")
			}
			msg = e.(string)
		}()
		fn()
		return msg
	}
	matchPanic := func(t *testing.T, expr string, posix bool) {
		t.Helper()
		// make sure we panic on subsequent calls
		for i := 0; i < 3; i++ {
			var got, exp string
			if posix {
				exp = mustPanic(t, func() { regexp.MustCompilePOSIX(expr) })
				got = mustPanic(t, func() { NewPOSIX(expr).MustCompile() })
			}
			if got != exp {
				t.Errorf("%d: %q: expected panic: %s got: %s", i, expr, exp, got)
			}
		}
	}

	matchPanic(t, "*", false)
	matchPanic(t, "*", true)
}

func TestString(t *testing.T) {
	exprs := []string{
		"aaa",
		"[",
		"*",
		"",
	}
	for _, expr := range exprs {
		got := New(expr).String()
		if got != expr {
			t.Errorf("String: want: %q got: %q", expr, got)
		}
		if re, _ := regexp.Compile(expr); re != nil {
			got := New(expr).String()
			if got != re.String() {
				t.Errorf("String: want: %q got: %q", re.String(), got)
			}
		}
	}
}

func buildMethodArgs(t *testing.T, method reflect.Value) []reflect.Value {
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))

	args := make([]reflect.Value, method.Type().NumIn())

	for i := 0; i < method.Type().NumIn(); i++ {
		typ := method.Type().In(i)
		switch typ.Kind() {
		case reflect.Func:
			// We need to use real funcs here in case they are called
			if typ.In(0).Kind() == reflect.String {
				fn := func(s string) string {
					return s + s
				}
				args[i] = reflect.ValueOf(fn)
			} else {
				fn := func(s []byte) []byte {
					return append(s, s...)
				}
				args[i] = reflect.ValueOf(fn)
			}
		case reflect.Interface:
			// We need to use a real io.RuneReader here in case it's called
			args[i] = reflect.ValueOf((io.RuneReader)(&bytes.Buffer{}))
		case reflect.String:
			args[i] = reflect.ValueOf("aaa")
		case reflect.Slice:
			switch typ.Elem().Kind() {
			case reflect.Uint8:
				args[i] = reflect.ValueOf([]byte("aaa"))
			case reflect.Int:
				args[i] = reflect.ValueOf([]int{1, 2})
			default:
				t.Errorf("Invalid slice type: %s", typ.String())
			}
		default:
			var ok bool
			args[i], ok = quick.Value(typ, rr)
			if !ok {
				t.Fatalf("Failed to create value for Type: %s", typ)
			}
		}
	}

	return args
}

func TestLazyCompile(t *testing.T) {
	const GoodPattern = ".*"

	onceCalled := func(re *Regexp) bool {
		called := true
		re.once.Do(func() { called = false })
		return called
	}

	testMethod := func(t *testing.T, newRe func(string) *Regexp, methodName string) {
		re := newRe(GoodPattern)

		m := reflect.ValueOf(re).MethodByName(methodName)

		m.Call(buildMethodArgs(t, m))
		if re.rx == nil {
			t.Error("Failed to initialize re.re: nil")
		}
		if !onceCalled(re) {
			t.Error("Failed to initialize re.re: once never called")
		}
		if re.String() != GoodPattern {
			t.Errorf("Want expr: %q got: %q", GoodPattern, re.rx.String())
		}
	}

	typ := reflect.TypeOf(&Regexp{})
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if m.Name == "String" {
			continue
		}
		t.Run(m.Name, func(t *testing.T) {
			testMethod(t, New, m.Name)
		})
		t.Run(m.Name+"_POSIX", func(t *testing.T) {
			testMethod(t, NewPOSIX, m.Name)
		})
	}
}

func TestLazyCompilePanic(t *testing.T) {
	const BadPattern = "*"
	const ErrMessage = "missing argument to repetition operator"

	testPanic := func(t *testing.T, methodName string) {
		e := recover()
		switch v := e.(type) {
		case nil:
			t.Errorf("%s: should have panicked", methodName)
		case string:
			if !strings.Contains(v, ErrMessage) {
				t.Errorf("%s: expected error message to contain: %q got: %q",
					methodName, ErrMessage, v)
			}
		default:
			t.Errorf("%s: unexpect panic type: %T", methodName, e)
		}
	}

	testMethod := func(t *testing.T, methodName string) {
		re := New(BadPattern)

		m := reflect.ValueOf(re).MethodByName(methodName)

		defer testPanic(t, methodName)
		m.Call(buildMethodArgs(t, m))
	}

	typ := reflect.TypeOf(&Regexp{})
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if m.Name == "Compile" || m.Name == "String" {
			continue
		}
		t.Run(m.Name, func(t *testing.T) {
			testMethod(t, m.Name)
		})
	}
}

func TestLazyCompileParallel(t *testing.T) {
	const expr = `(|(((((((((x{1}){1,7}){1,2}){1,2}){2}){2}){2}){2}){2})?` +
		`(((((((((x{1}){2,3}){2}){2}){2}){2}){2}){2}){2}){2})*`
	re := New(expr)
	start := make(chan struct{})
	wg := new(sync.WaitGroup)
	ready := new(sync.WaitGroup)

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		ready.Add(1)
		go func(i int) {
			defer wg.Done()
			ready.Done()
			<-start
			for i := 0; i < 100; i++ {
				if !re.MatchString("xxxxxxx") {
					t.Errorf("%d: failed to match string", i)
					return
				}
			}
		}(i)
	}

	ready.Wait()
	close(start)
	wg.Wait()
}

func BenchmarkInitOverhead(b *testing.B) {
	re := New("a")
	for i := 0; i < b.N; i++ {
		re.Longest() // cheapest method
	}
}

func BenchmarkInitOverhead_Baseline(b *testing.B) {
	re := regexp.MustCompile("a")
	for i := 0; i < b.N; i++ {
		re.Longest() // cheapest method
	}
}

func BenchmarkInitOverhead_Parallel(b *testing.B) {
	re := New("a")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			re.Longest()
		}
	})
	for i := 0; i < b.N; i++ {
		re.Longest() // cheapest method
	}
}

func BenchmarkInitOverhead_Parallel_Baseline(b *testing.B) {
	re := regexp.MustCompile("a")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			re.Longest()
		}
	})
}
