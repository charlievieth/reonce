// +build reoncetest

package reonce

// mustCompile forces compilation of the regex in New() and NewPOSIX()
// and is useful for testing Regexp's declared on initialization.
const mustCompile = true
