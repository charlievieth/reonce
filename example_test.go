package reonce_test

import (
	"fmt"

	"github.com/charlievieth/reonce"
)

// goExtRe is not initialized until one of its methods is called
// this can make program startup faster
var goExtRe = reonce.New(`\.go$`)

// Math file names against the regex `\.go$`
func ExampleRegexp() {
	for _, name := range []string{"main.go", "nope.py"} {
		fmt.Println(name, goExtRe.MatchString(name))
	}
	// Output:
	// main.go true
	// nope.py false
}
