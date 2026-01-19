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
	var {{.Name}} {{GoRpcType .Type}}
	{{- end}}

	if !parseRequest(w, r{{range .Args}}, &{{.Name}}{{end}}) { return }

	{{if ne $resData.Name "nil" -}}
	result, status := {{.Name | SnakeCase}}(r.Context()
		{{- range $i, $arg := .Args}}, {{if IsBaseType .Type}}{{GoLogicType .Type}}({{.Name}}){{else if IsEnum .Type}}{{if .Type.IsList}}{{.Name}}{{else}}{{PascalCase .Type.Name}}({{.Name}}){{end}}{{else if IsStruct .Type}}&{{.Name}}{{else}}{{.Name}}{{end}}{{end}})
	if !checkStatus(w, status) { return }
	{{if or (IsStruct $resData) (IsList $resData) -}}
	sendResponse(w, result)
	{{- else if IsEnum $resData -}}
	sendResponse(w, U8(result))
	{{- else -}}
	sendResponse(w, {{PascalCase $resData.Name}}(result))
	{{- end}}
	{{- else -}}
	status := {{.Name | SnakeCase}}(r.Context()
		{{- range $i, $arg := .Args}}, {{if IsBaseType .Type}}{{GoLogicType .Type}}({{.Name}}){{else if IsEnum .Type}}{{if .Type.IsList}}{{.Name}}{{else}}{{PascalCase .Type.Name}}({{.Name}}){{end}}{{else if IsStruct .Type}}&{{.Name}}{{else}}{{.Name}}{{end}}{{end}})
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
func Register{{$module | PascalCase}}(mux *http.ServeMux, mws ...Middleware) {
	mw := composeMiddleware(mws...)
{{- range $pkgApis}}
	mux.HandleFunc("POST /{{.Name}}", mw({{.Name | PascalCase}}Handler))
{{- end}}
}
{{end}}

// --- 内部辅助函数 ---

func checkStatus(w http.ResponseWriter, status RpcErrCode) bool {
	if status == RpcOk { return true }
	w.WriteHeader(int(status)); return false
}

func parseRequest(w http.ResponseWriter, r *http.Request, args ...Deserializable) bool {
	if len(args) == 0 { return true }
	body, err := io.ReadAll(r.Body); if err != nil { w.WriteHeader(http.StatusBadRequest); return false }
	if err := GetAll(bytes.NewBuffer(body), args...); err != nil { w.WriteHeader(http.StatusBadRequest); return false }
	return true
}

func sendResponse(w http.ResponseWriter, result Serializable) {
	var buf bytes.Buffer
	if err := SetAll(&buf, result); err != nil { w.WriteHeader(http.StatusInternalServerError); return }
	w.Write(buf.Bytes())
}
