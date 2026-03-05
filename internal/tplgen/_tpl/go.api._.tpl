package {{.Package}}

import (
	"bytes"
	"io"
	"net/http"
)

// --- API Handlers ---

{{range .Apis}}
{{- $resData := .Result -}}
{{- $handlerName := .Name | PascalCase -}}
func {{$handlerName}}Handler(w http.ResponseWriter, r *http.Request) {
	{{- range .Args}}
	var {{.Name}} {{GoLogicType .Type}}
	{{- end}}

	if !parseRequest(w, r{{range .Args}}, func(buf *bytes.Buffer) error {
		val, err := {{if and (IsStruct .Type) (not .Type.IsList)}}Read{{PascalCase .Type.Name}}{{else}}Get{{PascalCase .Type.Name}}{{if .Type.IsList}}List{{end}}{{end}}(buf)
		if err != nil { return err }
		{{.Name}} = val
		return nil
	}{{end}}) { return }
	{{- range .Args}}
	{{- if IsEnum .Type}}
	{{- if .Type.IsList}}
	if !Is{{PascalCase .Type.Name}}List({{.Name}}) { w.WriteHeader(http.StatusBadRequest); return }
	{{- else}}
	if !Is{{PascalCase .Type.Name}}({{.Name}}) { w.WriteHeader(http.StatusBadRequest); return }
	{{- end}}
	{{- else if IsStruct .Type}}
	{{- if .Type.IsList}}
	if err := Validate{{PascalCase .Type.Name}}List({{.Name}}); err != nil { w.WriteHeader(http.StatusBadRequest); return }
	{{- else}}
	if err := Validate{{PascalCase .Type.Name}}({{.Name}}); err != nil { w.WriteHeader(http.StatusBadRequest); return }
	{{- end}}
	{{- end}}
	{{- end}}

	{{if ne $resData.Name "nil" -}}
	result, status := {{.Name | SnakeCase}}(r.Context()
		{{- range $i, $arg := .Args}}, {{.Name}}{{end}})
	if !checkStatus(w, status) { return }
	sendResponse(w, func(buf *bytes.Buffer) error { return Set{{PascalCase $resData.Name}}{{if $resData.IsList}}List{{end}}(buf, {{if and (not (IsBaseType $resData)) $resData.IsList}}({{PascalCase $resData.Name}}List)(result){{else}}result{{end}}) })
	{{- else -}}
	status := {{.Name | SnakeCase}}(r.Context()
		{{- range $i, $arg := .Args}}, {{.Name}}{{end}})
	if !checkStatus(w, status) { return }
	w.WriteHeader(http.StatusOK)
	{{- end}}
}
{{end}}

// --- 路由注册 ---

type Middleware func(http.HandlerFunc) http.HandlerFunc

func composeMiddleware(mws ...Middleware) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		for i := len(mws) - 1; i >= 0; i-- { h = mws[i](h) }
		return h
	}
}

{{range $module, $pkgApis := .Groups}}
{{- $moduleName := $module | PascalCase}}
{{- if eq $moduleName "Api"}}
func RegisterApi(mux *http.ServeMux, mws ...Middleware) {
{{- else}}
func Register{{$moduleName}}Api(mux *http.ServeMux, mws ...Middleware) {
{{- end}}
	mw := composeMiddleware(mws...)
{{- range $pkgApis}}
	mux.HandleFunc("POST /{{.Name}}", mw({{.Name | PascalCase}}Handler))
{{- end}}
}
{{end}}

// --- 内部辅助函数 ---

const defaultMaxReqBytes int64 = 4 * 1024 * 1024

func checkStatus(w http.ResponseWriter, status RpcErrCode) bool {
	if status == RpcOk { return true }
	w.WriteHeader(int(status)); return false
}

func parseRequest(w http.ResponseWriter, r *http.Request, args ...Getter) bool {
	if len(args) == 0 {
		body, err := io.ReadAll(io.LimitReader(r.Body, 1)); if err != nil { w.WriteHeader(http.StatusBadRequest); return false }
		if len(body) != 0 { w.WriteHeader(http.StatusBadRequest); return false }
		return true
	}
	readLimit := defaultMaxReqBytes + 1
	body, err := io.ReadAll(io.LimitReader(r.Body, readLimit)); if err != nil { w.WriteHeader(http.StatusBadRequest); return false }
	if int64(len(body)) > defaultMaxReqBytes { w.WriteHeader(http.StatusRequestEntityTooLarge); return false }
	buf := bytes.NewBuffer(body)
	if err := GetAll(buf, args...); err != nil { w.WriteHeader(http.StatusBadRequest); return false }
	if buf.Len() != 0 { w.WriteHeader(http.StatusBadRequest); return false }
	return true
}

func sendResponse(w http.ResponseWriter, setter Setter) {
	var buf bytes.Buffer
	if err := SetAll(&buf, setter); err != nil { w.WriteHeader(http.StatusInternalServerError); return }
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	if _, err := w.Write(buf.Bytes()); err != nil { w.WriteHeader(http.StatusInternalServerError); return }
}
