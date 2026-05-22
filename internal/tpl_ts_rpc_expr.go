package internal

import (
	"fmt"
	"strings"
)

func (g *TsGenerator) rpcArgList(args []TplApiArg) string {
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, fmt.Sprintf("%s: %s", CamelCase(arg.Name), g.getTsLogicType(arg.Type)))
	}
	return strings.Join(parts, ", ")
}

func (g *TsGenerator) rpcPromiseType(result TplType) string {
	if result.Name == "nil" {
		return "RpcErrCode"
	}
	return fmt.Sprintf("[%s, RpcErrCode]", g.getTsLogicType(result))
}

func (g *TsGenerator) rpcDefaultValue(result TplType) string {
	if result.IsList {
		return "[]"
	}
	switch result.Kind {
	case TplKindBase:
		return g.getTsValue(result.Name)
	case TplKindEnum:
		return fmt.Sprintf("Default%s()", PascalCase(result.Name))
	case TplKindStruct:
		return fmt.Sprintf("new%s()", PascalCase(result.Name))
	default:
		return "null"
	}
}

func (g *TsGenerator) rpcReqErrReturn(result TplType, defaultVal string) string {
	if result.Name == "nil" {
		return "RpcErrCode.ReqErr"
	}
	return fmt.Sprintf("[%s, RpcErrCode.ReqErr]", defaultVal)
}

func (g *TsGenerator) rpcStatusReturn(result TplType, defaultVal string) string {
	if result.Name == "nil" {
		return "status"
	}
	return fmt.Sprintf("[%s, status]", defaultVal)
}

func (g *TsGenerator) apiReqStruct(api TplApi) (TplStruct, bool) {
	if len(api.Args) == 0 {
		return TplStruct{}, false
	}
	fields := make([]TplStructField, 0, len(api.Args))
	for _, arg := range api.Args {
		fields = append(fields, TplStructField{Name: arg.Name, Type: arg.Type})
	}
	return TplStruct{Name: "api_" + SnakeCase(api.Name) + "_req", Fields: fields}, true
}

func (g *TsGenerator) apiRespStruct(api TplApi) (TplStruct, bool) {
	if api.Result.Name == "nil" {
		return TplStruct{}, false
	}
	return TplStruct{
		Name:   "api_" + SnakeCase(api.Name) + "_resp",
		Fields: []TplStructField{{Name: "result", Type: api.Result}},
	}, true
}

func (g *TsGenerator) apiReqTypeName(api TplApi) (string, bool) {
	req, ok := g.apiReqStruct(api)
	if !ok {
		return "", false
	}
	return PascalCase(req.Name), true
}

func (g *TsGenerator) apiRespTypeName(api TplApi) (string, bool) {
	resp, ok := g.apiRespStruct(api)
	if !ok {
		return "", false
	}
	return PascalCase(resp.Name), true
}
