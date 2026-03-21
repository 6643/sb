package internal

import _ "embed"

//go:embed tpl_go_runtime.go.txt
var goRuntimeSource string

func renderGoRuntimeSource() string {
	return goRuntimeSource
}
