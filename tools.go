//go:build tools
// +build tools

package tools

import (
	_ "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment"
	_ "golang.org/x/tools/go/analysis/passes/findcall/cmd/findcall"
	_ "golang.org/x/tools/go/analysis/passes/ifaceassert/cmd/ifaceassert"
	_ "golang.org/x/tools/go/analysis/passes/lostcancel/cmd/lostcancel"
	_ "golang.org/x/tools/go/analysis/passes/nilness/cmd/nilness"
	_ "golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow"
	_ "golang.org/x/tools/go/analysis/passes/stringintconv/cmd/stringintconv"
	_ "golang.org/x/tools/go/analysis/passes/unmarshal/cmd/unmarshal"
)
