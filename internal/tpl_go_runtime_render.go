package internal

import (
	_ "embed"
	"strings"
)

//go:embed tpl_go_runtime_core.go.txt
var goRuntimeCoreSource string

//go:embed tpl_go_runtime_enum.go.txt
var goRuntimeEnumSource string

//go:embed tpl_go_runtime_struct.go.txt
var goRuntimeStructSource string

//go:embed tpl_go_runtime_bits.go.txt
var goRuntimeBitsSource string

//go:embed tpl_go_runtime_base.go.txt
var goRuntimeBaseSource string

//go:embed tpl_go_runtime_text.go.txt
var goRuntimeTextSource string

//go:embed tpl_go_runtime_bool_list.go.txt
var goRuntimeBoolListSource string

//go:embed tpl_go_runtime_bin.go.txt
var goRuntimeBinSource string

//go:embed tpl_go_runtime_list.go.txt
var goRuntimeListSource string

//go:embed tpl_go_runtime_text_list.go.txt
var goRuntimeTextListSource string

func renderGoRuntimeSource() string {
	parts := []string{
		goRuntimeCoreSource,
		goRuntimeEnumSource,
		goRuntimeStructSource,
		goRuntimeBitsSource,
		goRuntimeBaseSource,
		goRuntimeTextSource,
		goRuntimeBoolListSource,
		goRuntimeBinSource,
		goRuntimeListSource,
		goRuntimeTextListSource,
	}
	return strings.Join(parts, "\n") + "\n"
}
