{{- $resultData := .Api.Result -}}
{{- $innerFuncName := .Api.Name | SnakeCase -}}
{{- $retType := GoLogicType $resultData -}}
{{- $hasRet := ne $resultData.Name "nil" -}}

package {{.Package}}

import (
	"context"
)

func {{$innerFuncName}}(ctx context.Context{{range .Api.Args}}, {{.Name}} {{GoLogicType .Type}}{{end}}) ({{if $hasRet}}result {{$retType}}, {{end}}errCode RpcErrCode) {
	return {{if $hasRet}}{{GoValue .Api.Result.Name}}, {{end}}RpcRespErr
}
