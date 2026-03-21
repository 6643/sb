package internal

import _ "embed"

//go:embed tpl_ts_runtime.ts.txt
var tsRuntimeSource string

func renderTsRuntimeSource() string {
	return tsRuntimeSource
}
