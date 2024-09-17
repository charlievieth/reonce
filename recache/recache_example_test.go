package recache_test

import (
	"fmt"

	"github.com/charlievieth/reonce/recache"
)

func ExampleMustCompile() {
	const text = "The quick brown fox jumps over the lazy dog."
	// Here search could be some user input that is likely to
	// be repeated and thus worth caching.
	for _, search := range []string{
		"fox",
		"dog",
		"fox", // reused
	} {
		// This uses the global cache
		re := recache.MustCompile(search)
		fmt.Println(re.ReplaceAllString(text, "cat"))
	}
	// Output:
	// The quick brown cat jumps over the lazy dog.
	// The quick brown fox jumps over the lazy cat.
	// The quick brown cat jumps over the lazy dog.
}

func ExampleNew() {
	cache := recache.New(8)
	for _, expr := range []string{"foo", "bar", "foo"} {
		fmt.Printf("%q: %t\n", expr, cache.MustCompile(expr).MatchString("foo"))
	}
	fmt.Println(cache.Len())
	fmt.Println(cache.POSIX())
	// Output:
	// "foo": true
	// "bar": false
	// "foo": true
	// 2
	// false
}
