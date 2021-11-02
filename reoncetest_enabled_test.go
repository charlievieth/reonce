//go:build reoncetest
// +build reoncetest

package reonce

import "testing"

func TestCompileAlways(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		defer func() {
			if e := recover(); e == nil {
				t.Error("expected New() to panic when built with the 'reoncetest' tag")
			}
		}()
		New("[")
	})

	t.Run("NewPOSIX", func(t *testing.T) {
		defer func() {
			if e := recover(); e == nil {
				t.Error("expected NewPOSIX() to panic when built with the 'reoncetest' tag")
			}
		}()
		NewPOSIX("[")
	})
}
