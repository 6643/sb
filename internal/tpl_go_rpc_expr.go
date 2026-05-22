package internal

import "fmt"

func (g *GoGenerator) logicType(t TplType) string {
	return g.getGoLogicType(t)
}

func (g *GoGenerator) defaultReturn(t TplType) string {
	if t.Kind == TplKindStruct || t.IsList {
		return "nil"
	}
	return g.getGoValue(t.Name)
}

func (g *GoGenerator) stubArgs(args []TplApiArg) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, fmt.Sprintf("%s %s", arg.Name, g.logicType(arg.Type)))
	}
	return ", " + joinWithComma(parts)
}

func (g *GoGenerator) rpcArgs(args []TplApiArg) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, fmt.Sprintf("%s %s", CamelCase(arg.Name), g.logicType(arg.Type)))
	}
	return ", " + joinWithComma(parts)
}

func (g *GoGenerator) rpcReturn(result TplType) string {
	if result.Name == "nil" {
		return "errCode RpcErrCode"
	}
	return fmt.Sprintf("result %s, errCode RpcErrCode", g.logicType(result))
}

func (g *GoGenerator) rpcReqErrReturn(result TplType) string {
	if result.Name == "nil" {
		return "RpcReqErr"
	}
	return "result, RpcReqErr"
}

func (g *GoGenerator) rpcStatusReturn(result TplType) string {
	if result.Name == "nil" {
		return "status"
	}
	return "result, status"
}

func (g *GoGenerator) apiReqStruct(api TplApi) (TplStruct, bool) {
	if len(api.Args) == 0 {
		return TplStruct{}, false
	}
	fields := make([]TplStructField, 0, len(api.Args))
	for _, arg := range api.Args {
		fields = append(fields, TplStructField{
			Name: arg.Name,
			Type: arg.Type,
		})
	}
	return TplStruct{
		Name:   "api_" + SnakeCase(api.Name) + "_req",
		Fields: fields,
	}, true
}

func (g *GoGenerator) apiRespStruct(api TplApi) (TplStruct, bool) {
	if api.Result.Name == "nil" {
		return TplStruct{}, false
	}
	return TplStruct{
		Name: "api_" + SnakeCase(api.Name) + "_resp",
		Fields: []TplStructField{
			{Name: "result", Type: api.Result},
		},
	}, true
}

func (g *GoGenerator) apiReqTypeName(api TplApi) (string, bool) {
	req, ok := g.apiReqStruct(api)
	if !ok {
		return "", false
	}
	return PascalCase(req.Name), true
}

func (g *GoGenerator) apiRespTypeName(api TplApi) (string, bool) {
	resp, ok := g.apiRespStruct(api)
	if !ok {
		return "", false
	}
	return PascalCase(resp.Name), true
}

func (g *GoGenerator) apiReqCodecInfo(api TplApi) (typeName string, readName string, setName string, ok bool) {
	req, ok := g.apiReqStruct(api)
	if !ok {
		return "", "", "", false
	}
	return PascalCase(req.Name), g.structReadName(req.Name), g.structSetName(req.Name), true
}

func (g *GoGenerator) apiRespCodecInfo(api TplApi) (typeName string, readName string, setName string, ok bool) {
	resp, ok := g.apiRespStruct(api)
	if !ok {
		return "", "", "", false
	}
	return PascalCase(resp.Name), g.structReadName(resp.Name), g.structSetName(resp.Name), true
}

func (g *GoGenerator) handlerCallArgs(args []TplApiArg, reqVar string) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, reqVar+"."+PascalCase(arg.Name))
	}
	return ", " + joinWithComma(parts)
}

func (g *GoGenerator) callArgNames(args []TplApiArg) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, CamelCase(arg.Name))
	}
	return ", " + joinWithComma(parts)
}
