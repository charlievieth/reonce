package reonce_test

import (
	"fmt"

	"github.com/charlievieth/reonce"
)

// ExtRe matches the extension of a file with a non-empty name.
// Typically, reonce regexes should be declared globally via a
// var or an init() function.
//
// NB: this is just an example and filepath.Ext() should really
// be used instead.
var ExtRe = reonce.New(`(?m)(?:\w+)(\.\w+)$`)

func ExampleRegexp_FindStringSubmatch() {
	fmt.Printf("%q\n", ExtRe.FindStringSubmatch("main.go"))
	fmt.Printf("%q\n", ExtRe.FindStringSubmatch("nope.py"))
	fmt.Printf("%q\n", ExtRe.FindStringSubmatch("noext"))
	fmt.Printf("%q\n", ExtRe.FindStringSubmatch(".extonly"))
	// Output:
	// ["main.go" ".go"]
	// ["nope.py" ".py"]
	// []
	// []
}
