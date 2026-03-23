package internal

import _ "embed"

//go:embed tpl_ts_runtime.ts.txt
var tsRuntimeSource string

//go:embed tpl_ts_runtime_core.ts.txt
var tsRuntimeCoreSource string

//go:embed tpl_ts_runtime_base.ts.txt
var tsRuntimeBaseSource string

//go:embed tpl_ts_runtime_header.ts.txt
var tsRuntimeHeaderSource string

//go:embed tpl_ts_runtime_text.ts.txt
var tsRuntimeTextSource string

//go:embed tpl_ts_runtime_bin.ts.txt
var tsRuntimeBinSource string

//go:embed tpl_ts_runtime_list.ts.txt
var tsRuntimeListSource string

//go:embed tpl_ts_runtime_meta.ts.txt
var tsRuntimeMetaSource string

//go:embed tpl_ts_runtime_enum.ts.txt
var tsRuntimeEnumSource string

//go:embed tpl_ts_runtime_struct.ts.txt
var tsRuntimeStructSource string

func renderTsRuntimeSource() string {
	return tsRuntimeSource
}

func renderTsRuntimeCoreSource() string {
	return tsRuntimeCoreSource
}

func renderTsRuntimeBaseSource() string {
	return tsRuntimeBaseSource
}

func renderTsRuntimeHeaderSource() string {
	return tsRuntimeHeaderSource
}

func renderTsRuntimeTextSource() string {
	return tsRuntimeTextSource
}

func renderTsRuntimeBinSource() string {
	return tsRuntimeBinSource
}

func renderTsRuntimeListSource() string {
	return tsRuntimeListSource
}

func renderTsRuntimeMetaSource() string {
	return tsRuntimeMetaSource
}

func renderTsRuntimeEnumSource() string {
	return tsRuntimeEnumSource
}

func renderTsRuntimeStructSource() string {
	return tsRuntimeStructSource
}
