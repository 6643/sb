package internal

import (
	"fmt"
	"strings"
)

func (g *GoGenerator) renderGoEnumFile(enums []TplEnum) string {
	var w sourceWriter
	w.Line("package sb")
	if len(enums) == 0 {
		return w.String()
	}
	w.Blank()
	w.Line("import (")
	w.Line("\t\"bytes\"")
	w.Line("\t\"slices\"")
	w.Line("\t\"unsafe\"")
	w.Line(")")
	for _, enum := range enums {
		enumName := PascalCase(enum.Name)
		w.Blank()
		w.WriteLineCommentWithHead("// ", enumName+" ", enum.Note)
		w.Linef("type %s uint8", enumName)
		w.Blank()
		w.Line("const (")
		for _, child := range enum.Children {
			w.WriteLineComment("\t// ", child.Note)
			w.Linef("\t%s%s %s = %d", enumName, PascalCase(child.Name), enumName, child.ID)
		}
		w.Line(")")
		w.Blank()
		w.Linef("func Is%s(v %s) bool {", enumName, enumName)
		w.Line("\tswitch v {")
		cases := make([]string, 0, len(enum.Children))
		for _, child := range enum.Children {
			cases = append(cases, enumName+PascalCase(child.Name))
		}
		w.Linef("\tcase %s:", joinWithComma(cases))
		w.Line("\t\treturn true")
		w.Line("\tdefault:")
		w.Line("\t\treturn false")
		w.Line("\t}")
		w.Line("}")
		w.Blank()
		w.Linef("func Get%s(buf *bytes.Buffer) (%s, error) {", enumName, enumName)
		w.Line("\tval, err := GetU8(buf)")
		w.Linef("\treturn %s(val), err", enumName)
		w.Line("}")
		w.Blank()
		w.Linef("func Set%s(buf *bytes.Buffer, v %s) error { return SetU8(buf, uint8(v)) }", enumName, enumName)
		w.Blank()
		w.Linef("type %sList []%s", enumName, enumName)
		w.Linef("func Get%sList(buf *bytes.Buffer) (%sList, error) {", enumName, enumName)
		w.Line("\tval, err := GetU8List(buf)")
		w.Line("\tif err != nil { return nil, err }")
		w.Linef("\treturn *(*%sList)(unsafe.Pointer(&val)), nil", enumName)
		w.Line("}")
		w.Linef("func Set%sList(buf *bytes.Buffer, v %sList) error { return SetU8List(buf, *(*[]uint8)(unsafe.Pointer(&v))) }", enumName, enumName)
		w.Linef("func Is%sList(v %sList) bool {", enumName, enumName)
		w.Line("\tfor _, item := range v {")
		w.Linef("\t\tif !Is%s(item) { return false }", enumName)
		w.Line("\t}")
		w.Line("\treturn true")
		w.Line("}")
		w.Linef("func Eq%sList(a, b %sList) bool { return slices.Equal(a, b) }", enumName, enumName)
	}
	return w.String()
}

func (g *GoGenerator) renderGoStructFile(st TplStruct) string {
	structName := PascalCase(st.Name)
	bitSize := goBitmaskSize(len(st.Fields))
	var w sourceWriter
	w.Line("package sb")
	w.Blank()
	w.Line("import (")
	w.Line("\t\"bytes\"")
	w.Line("\t\"fmt\"")
	w.Line("\t\"slices\"")
	w.Line(")")
	w.Blank()
	w.Linef("type %s struct {", structName)
	for _, field := range st.Fields {
		w.WriteLineComment("\t// ", field.Note)
		tag := g.getGoTag(field)
		if tag == "" {
			w.Linef("\t%s %s", PascalCase(field.Name), g.getGoLogicType(field.Type))
			continue
		}
		w.Linef("\t%s %s %s", PascalCase(field.Name), g.getGoLogicType(field.Type), tag)
	}
	w.Line("}")
	w.Blank()
	w.Linef("func size%sBody(s *%s, bits []byte) (int, error) {", structName, structName)
	w.Linef("\tif s == nil { return 0, fmt.Errorf(\"size%s: nil value\") }", structName)
	w.Line("\tbodySize := 0")
	for i, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		fieldRef := "s." + fieldName
		if field.Type.Name == "bool" {
			w.Linef("\tif bits != nil { SetBit(bits, uint8(%d), %s) }", i, fieldRef)
			continue
		}
		checkExpr := g.goPresenceCheck(field.Type, fieldRef)
		w.Linef("\tif %s {", checkExpr)
		g.writeGoSizeAddLines(&w, "size"+structName, fieldName, field.Type, fieldRef, "bodySize", "\t\t")
		w.Linef("\t\tif bits != nil { SetBit(bits, uint8(%d), true) }", i)
		w.Line("\t}")
	}
	w.Line("\treturn bodySize, nil")
	w.Line("}")
	w.Blank()
	w.Linef("func size%s(s *%s) (int, error) {", structName, structName)
	w.Linef("\tbodySize, err := size%sBody(s, nil)", structName)
	w.Line("\tif err != nil { return 0, err }")
	w.Linef("\treturn %d + bodySize, nil", bitSize)
	w.Line("}")
	w.Blank()
	w.Linef("func size%sList(v %sList) (int, error) {", structName, structName)
	w.Line("\tif len(v) > 65535 { return 0, fmt.Errorf(\"list length exceeds uint16 max\") }")
	w.Line("\ttotalSize := 2")
	w.Line("\tfor i, item := range v {")
	w.Linef("\t\titemSize, err := size%s(item)", structName)
	w.Linef("\t\tif err != nil { return 0, fmt.Errorf(\"size%sList[%%d]: %%w\", i, err) }", structName)
	w.Line("\t\ttotalSize += itemSize")
	w.Line("\t}")
	w.Line("\treturn totalSize, nil")
	w.Line("}")
	w.Blank()
	w.Linef("func Get%s(buf *bytes.Buffer, s *%s) error {", structName, structName)
	w.Line("\tif s == nil { return nil }")
	w.Linef("\tconst bitSize = %d", bitSize)
	w.Linef("\tif buf.Len() < bitSize { return fmt.Errorf(\"Get%s bitmask: %%d - %%d\", buf.Len(), bitSize) }", structName)
	w.Line("\tbits := buf.Next(bitSize)")
	for i, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		if field.Type.Name == "bool" {
			w.Linef("\ts.%s = GetBit(bits, uint8(%d))", fieldName, i)
			continue
		}
		w.Linef("\tif GetBit(bits, uint8(%d)) {", i)
		readExpr := g.goReadExpr(field.Type, "buf")
		w.Linef("\t\tval, err := %s", readExpr)
		w.Linef("\t\tif err != nil { return fmt.Errorf(\"Get%s %s: %%w\", err) }", structName, fieldName)
		w.Linef("\t\ts.%s = val", fieldName)
		if field.Type.Kind == TplKindEnum && !field.Type.IsList {
			w.Linef("\t\tif !Is%s(s.%s) { return fmt.Errorf(\"Get%s %s: 非法枚举值: %%d\", s.%s) }", PascalCase(field.Type.Name), fieldName, structName, fieldName, fieldName)
		}
		w.Line("\t}")
	}
	w.Linef("\tif err := Validate%s(s); err != nil { return fmt.Errorf(\"Validate%s: %%w\", err) }", structName, structName)
	w.Line("\treturn nil")
	w.Line("}")
	w.Blank()
	w.Linef("func Set%s(buf *bytes.Buffer, s *%s) error {", structName, structName)
	w.Linef("\tif s == nil { return fmt.Errorf(\"Set%s: nil value\") }", structName)
	w.Linef("\tif err := Validate%s(s); err != nil { return fmt.Errorf(\"Validate%s: %%w\", err) }", structName, structName)
	w.Linef("\tconst bitSize = %d", bitSize)
	w.Linef("\tvar bits [%d]byte", bitSize)
	w.Linef("\tbodySize, err := size%sBody(s, bits[:])", structName)
	w.Line("\tif err != nil { return err }")
	w.Line("\tbuf.Grow(bitSize + bodySize)")
	w.Linef("\tif _, err := buf.Write(bits[:]); err != nil { return fmt.Errorf(\"Set%s write bitmask: %%w\", err) }", structName)
	for _, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		fieldRef := "s." + fieldName
		if field.Type.Name == "bool" {
			continue
		}
		checkExpr := g.goPresenceCheck(field.Type, fieldRef)
		w.Linef("\tif %s {", checkExpr)
		w.Linef("\t\tif err := %s; err != nil { return fmt.Errorf(\"Set%s %s: %%w\", err) }", g.goSetCall(field.Type, "buf", fieldRef), structName, fieldName)
		w.Line("\t}")
	}
	w.Line("\treturn nil")
	w.Line("}")
	w.Blank()
	w.Linef("func Validate%s(s *%s) error {", structName, structName)
	w.Line("\tif s == nil { return nil }")
	for _, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		fieldRef := "s." + fieldName
		switch field.Type.Kind {
		case TplKindEnum:
			if field.Type.IsList {
				w.Linef("\tfor i, item := range %s {", fieldRef)
				w.Linef("\t\tif !Is%s(item) { return fmt.Errorf(\"%s[%%d] 非法枚举值: %%d\", i, item) }", PascalCase(field.Type.Name), fieldName)
				w.Line("\t}")
				continue
			}
			w.Linef("\tif %s != 0 && !Is%s(%s) { return fmt.Errorf(\"%s 非法枚举值: %%d\", %s) }", fieldRef, PascalCase(field.Type.Name), fieldRef, fieldName, fieldRef)
		case TplKindStruct:
			if field.Type.IsList {
				w.Linef("\tif err := Validate%sList(%s); err != nil { return fmt.Errorf(\"%s: %%w\", err) }", PascalCase(field.Type.Name), fieldRef, fieldName)
				continue
			}
			w.Linef("\tif %s != nil {", fieldRef)
			w.Linef("\t\tif err := Validate%s(%s); err != nil { return fmt.Errorf(\"%s: %%w\", err) }", PascalCase(field.Type.Name), fieldRef, fieldName)
			w.Line("\t}")
		}
	}
	w.Line("\treturn nil")
	w.Line("}")
	w.Blank()
	w.Linef("func Eq%s(a, b *%s) bool {", structName, structName)
	w.Line("\tif a == b { return true }")
	w.Line("\tif a == nil || b == nil { return false }")
	for _, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		left := "a." + fieldName
		right := "b." + fieldName
		w.Linef("\tif !%s { return false }", g.goEqExpr(field.Type, left, right))
	}
	w.Line("\treturn true")
	w.Line("}")
	w.Blank()
	w.Line("// Standalone functions")
	w.Linef("func Read%s(buf *bytes.Buffer) (*%s, error) {", structName, structName)
	w.Linef("\ts := new(%s)", structName)
	w.Linef("\treturn s, Get%s(buf, s)", structName)
	w.Line("}")
	w.Blank()
	w.Linef("type %sList []*%s", structName, structName)
	w.Linef("func Get%sList(buf *bytes.Buffer) (%sList, error) {", structName, structName)
	w.Line("\tcount, err := GetU16(buf)")
	w.Line("\tif err != nil { return nil, err }")
	w.Linef("\tlist := make([]*%s, count)", structName)
	w.Line("\tfor i := range list {")
	w.Linef("\t\titem, err := Read%s(buf)", structName)
	w.Line("\t\tif err != nil { return nil, err }")
	w.Line("\t\tlist[i] = item")
	w.Line("\t}")
	w.Linef("\treturn %sList(list), nil", structName)
	w.Line("}")
	w.Linef("func Set%sList(buf *bytes.Buffer, v %sList) error {", structName, structName)
	w.Linef("\ttotalSize, err := size%sList(v)", structName)
	w.Line("\tif err != nil { return err }")
	w.Line("\tbuf.Grow(totalSize)")
	w.Line("\tif err := SetU16(buf, uint16(len(v))); err != nil { return err }")
	w.Line("\tfor _, item := range v {")
	w.Linef("\t\tif err := Set%s(buf, item); err != nil { return err }", structName)
	w.Line("\t}")
	w.Line("\treturn nil")
	w.Line("}")
	w.Linef("func Validate%sList(v %sList) error {", structName, structName)
	w.Line("\tfor i, item := range v {")
	w.Linef("\t\tif item == nil { return fmt.Errorf(\"%sList[%%d]: nil item\", i) }", structName)
	w.Linef("\t\tif err := Validate%s(item); err != nil { return fmt.Errorf(\"%sList[%%d]: %%w\", i, err) }", structName, structName)
	w.Line("\t}")
	w.Line("\treturn nil")
	w.Line("}")
	w.Linef("func Eq%sList(a, b %sList) bool { return slices.EqualFunc(a, b, Eq%s) }", structName, structName, structName)
	return w.String()
}

func (g *GoGenerator) renderGoAPIStub(api TplApi) []byte {
	apiName := SnakeCase(api.Name)
	var w sourceWriter
	w.Line("package sb")
	w.Blank()
	w.Line("import (")
	w.Line("\t\"context\"")
	w.Line(")")
	w.Blank()
	if api.Result.Name == "nil" {
		w.Linef("func %s(ctx context.Context%s) (errCode RpcErrCode) {", apiName, g.goStubArgs(api.Args))
		w.Line("\treturn RpcRespErr")
		w.Line("}")
		return []byte(w.String())
	}
	retType := g.goStubResultType(api.Result)
	w.Linef("func %s(ctx context.Context%s) (result %s, errCode RpcErrCode) {", apiName, g.goStubArgs(api.Args), retType)
	w.Linef("\treturn %s, RpcRespErr", g.goDefaultReturn(api.Result))
	w.Line("}")
	return []byte(w.String())
}

func (g *GoGenerator) renderGoAPIHandlers(apis []TplApi) string {
	groups := groupAPIs(apis)
	var w sourceWriter
	w.Line("package sb")
	w.Blank()
	w.Line("import (")
	w.Line("\t\"bytes\"")
	w.Line("\t\"context\"")
	w.Line("\t\"io\"")
	w.Line("\t\"net/http\"")
	w.Line("\t\"time\"")
	w.Line(")")
	w.Blank()
	w.Line("// --- API Handlers ---")
	for _, api := range apis {
		callName := PascalCase(api.Name)
		w.Blank()
		w.Linef("func %sHandler(w http.ResponseWriter, r *http.Request) {", callName)
		for _, arg := range api.Args {
			w.Linef("\tvar %s %s", arg.Name, g.goHandlerArgType(arg.Type))
		}
		w.Blank()
		if len(api.Args) == 0 {
			w.Line("\tif !parseEmptyRequest(w, r) { return }")
		} else {
			w.Line("\tbuf, ok := parseRequest(w, r)")
			w.Line("\tif !ok { return }")
			for _, arg := range api.Args {
				w.Line("\t{")
				w.Linef("\t\tval, err := %s", g.goReadExpr(arg.Type, "buf"))
				w.Line("\t\tif err != nil { w.WriteHeader(http.StatusBadRequest); return }")
				w.Linef("\t\t%s = val", arg.Name)
				w.Line("\t}")
			}
			w.Line("\tif buf.Len() != 0 { w.WriteHeader(http.StatusBadRequest); return }")
		}
		for _, arg := range api.Args {
			name := arg.Name
			typeName := PascalCase(arg.Type.Name)
			switch arg.Type.Kind {
			case TplKindEnum:
				if arg.Type.IsList {
					w.Linef("\tif !Is%sList(%s) { w.WriteHeader(http.StatusBadRequest); return }", typeName, name)
				} else {
					w.Linef("\tif !Is%s(%s) { w.WriteHeader(http.StatusBadRequest); return }", typeName, name)
				}
			case TplKindStruct:
				if arg.Type.IsList {
					w.Linef("\tif err := Validate%sList(%s); err != nil { w.WriteHeader(http.StatusBadRequest); return }", typeName, name)
				} else {
					w.Linef("\tif err := Validate%s(%s); err != nil { w.WriteHeader(http.StatusBadRequest); return }", typeName, name)
				}
			}
		}
		w.Blank()
		callArgs := g.goHandlerCallArgs(api.Args)
		if api.Result.Name == "nil" {
			w.Linef("\tstatus := %s(r.Context()%s)", SnakeCase(api.Name), callArgs)
			w.Line("\tstatus = normalizeHandlerStatus(r.Context(), status)")
			w.Line("\tif !checkStatus(w, status) { return }")
			w.Line("\tw.WriteHeader(http.StatusOK)")
			w.Line("}")
			continue
		}
		w.Linef("\tresult, status := %s(r.Context()%s)", SnakeCase(api.Name), callArgs)
		w.Line("\tstatus = normalizeHandlerStatus(r.Context(), status)")
		w.Line("\tif !checkStatus(w, status) { return }")
		w.Line("\tvar resp bytes.Buffer")
		w.Line("\trespSize := 0")
		g.writeGoSizeAccumulateLines(&w, api.Result, "result", "respSize", "\t", "w.WriteHeader(http.StatusInternalServerError); return")
		w.Line("\tresp.Grow(respSize)")
		w.Linef("\tif err := %s; err != nil { w.WriteHeader(http.StatusInternalServerError); return }", g.goSetCall(api.Result, "&resp", "result"))
		w.Line("\tsendResponse(w, resp.Bytes())")
		w.Line("}")
	}
	w.Blank()
	w.Line("// --- 路由注册 ---")
	w.Blank()
	w.Line("type Middleware func(http.HandlerFunc) http.HandlerFunc")
	w.Blank()
	w.Line("func composeMiddleware(mws ...Middleware) func(http.HandlerFunc) http.HandlerFunc {")
	w.Line("\treturn func(h http.HandlerFunc) http.HandlerFunc {")
	w.Line("\t\tfor i := len(mws) - 1; i >= 0; i-- { h = mws[i](h) }")
	w.Line("\t\treturn h")
	w.Line("\t}")
	w.Line("}")
	w.Blank()
	w.Line("func TimeoutMiddleware(timeout time.Duration) Middleware {")
	w.Line("\treturn func(h http.HandlerFunc) http.HandlerFunc {")
	w.Line("\t\tif timeout <= 0 { return h }")
	w.Line("\t\treturn func(w http.ResponseWriter, r *http.Request) {")
	w.Line("\t\t\tctx, cancel := context.WithTimeout(r.Context(), timeout)")
	w.Line("\t\t\tdefer cancel()")
	w.Line("\t\t\th(w, r.WithContext(ctx))")
	w.Line("\t\t}")
	w.Line("\t}")
	w.Line("}")
	modules := orderedGroupKeys(groups)
	for _, module := range modules {
		funcName := "Register" + PascalCase(module) + "Api"
		if PascalCase(module) == "Api" {
			funcName = "RegisterApi"
		}
		w.Blank()
		w.Linef("func %s(mux *http.ServeMux, mws ...Middleware) {", funcName)
		w.Line("\tmw := composeMiddleware(mws...)")
		for _, api := range groups[module] {
			w.Linef("\tmux.HandleFunc(\"POST /%s\", mw(%sHandler))", api.Name, PascalCase(api.Name))
		}
		w.Line("}")
	}
	w.Blank()
	w.Line("// --- 内部辅助函数 ---")
	w.Blank()
	w.Line("const defaultMaxReqBytes int64 = 4 * 1024 * 1024")
	w.Blank()
	w.Line("func normalizeHandlerStatus(ctx context.Context, status RpcErrCode) RpcErrCode {")
	w.Line("\tif status != RpcOk || ctx == nil { return status }")
	w.Line("\tif ctx.Err() == context.DeadlineExceeded { return RpcTimeout }")
	w.Line("\treturn status")
	w.Line("}")
	w.Blank()
	w.Line("func checkStatus(w http.ResponseWriter, status RpcErrCode) bool {")
	w.Line("\tif status == RpcOk { return true }")
	w.Line("\tw.WriteHeader(int(status)); return false")
	w.Line("}")
	w.Blank()
	w.Line("func parseEmptyRequest(w http.ResponseWriter, r *http.Request) bool {")
	w.Line("\tbody, err := io.ReadAll(io.LimitReader(r.Body, 1)); if err != nil { w.WriteHeader(http.StatusBadRequest); return false }")
	w.Line("\tif len(body) != 0 { w.WriteHeader(http.StatusBadRequest); return false }")
	w.Line("\treturn true")
	w.Line("}")
	w.Blank()
	w.Line("func parseRequest(w http.ResponseWriter, r *http.Request) (*bytes.Buffer, bool) {")
	w.Line("\treadLimit := defaultMaxReqBytes + 1")
	w.Line("\tbody, err := io.ReadAll(io.LimitReader(r.Body, readLimit)); if err != nil { w.WriteHeader(http.StatusBadRequest); return nil, false }")
	w.Line("\tif int64(len(body)) > defaultMaxReqBytes { w.WriteHeader(http.StatusRequestEntityTooLarge); return nil, false }")
	w.Line("\treturn bytes.NewBuffer(body), true")
	w.Line("}")
	w.Blank()
	w.Line("func sendResponse(w http.ResponseWriter, body []byte) {")
	w.Line("\tif w.Header().Get(\"Content-Type\") == \"\" {")
	w.Line("\t\tw.Header().Set(\"Content-Type\", \"application/octet-stream\")")
	w.Line("\t}")
	w.Line("\tif _, err := w.Write(body); err != nil { w.WriteHeader(http.StatusInternalServerError); return }")
	w.Line("}")
	return w.String()
}

func (g *GoGenerator) renderGoRPCFile(apis []TplApi) string {
	var w sourceWriter
	w.Line("package sb")
	w.Blank()
	w.Line("import (")
	w.Line("\t\"bytes\"")
	w.Line("\t\"context\"")
	w.Line("\t\"io\"")
	w.Line("\t\"net/http\"")
	w.Line("\t\"strings\"")
	w.Line("\t\"sync\"")
	w.Line("\t\"time\"")
	w.Line(")")
	w.Blank()
	w.Line("type RpcErrCode int")
	w.Blank()
	w.Line("const (")
	w.Line("\tRpcOk        RpcErrCode = 200")
	w.Line("\tRpcNoConn    RpcErrCode = 0")
	w.Line("\tRpcTimeout   RpcErrCode = 408")
	w.Line("\tRpcReqErr    RpcErrCode = 400")
	w.Line("\tRpcRespErr   RpcErrCode = 500")
	w.Line("\tRpcNotAuth   RpcErrCode = 401")
	w.Line("\tRpcNotExist  RpcErrCode = 404")
	w.Line("\tdefaultMaxRespBytes int64 = 4 * 1024 * 1024")
	w.Line(")")
	w.Blank()
	w.Line("type Client struct {")
	w.Line("\tBaseURL string")
	w.Line("\tHTTP    *http.Client")
	w.Line("\tTimeout time.Duration")
	w.Line("\tRetries int")
	w.Line("\tMaxRespBytes int64")
	w.Line("\theaders map[string]string")
	w.Line("\theadersMu sync.RWMutex")
	w.Line("}")
	w.Blank()
	w.Line("func NewClient(baseURL string) *Client {")
	w.Line("\tbaseURL = strings.TrimRight(baseURL, \"/\")")
	w.Line("\treturn &Client{")
	w.Line("\t\tBaseURL: baseURL,")
	w.Line("\t\tHTTP:    &http.Client{},")
	w.Line("\t\tTimeout: 5 * time.Second,")
	w.Line("\t\tRetries: 3,")
	w.Line("\t\tMaxRespBytes: defaultMaxRespBytes,")
	w.Line("\t\theaders: make(map[string]string),")
	w.Line("\t}")
	w.Line("}")
	w.Blank()
	w.Line("func SetClientHeader(c *Client, key, value string) {")
	w.Line("\tif c == nil { return }")
	w.Line("\tc.headersMu.Lock()")
	w.Line("\tdefer c.headersMu.Unlock()")
	w.Line("\tif c.headers == nil { c.headers = make(map[string]string) }")
	w.Line("\tc.headers[key] = value")
	w.Line("}")
	w.Line("func GetClientHeader(c *Client, key string) string {")
	w.Line("\tif c == nil { return \"\" }")
	w.Line("\tc.headersMu.RLock()")
	w.Line("\tdefer c.headersMu.RUnlock()")
	w.Line("\treturn c.headers[key]")
	w.Line("}")
	w.Line("func RemoveClientHeader(c *Client, key string) {")
	w.Line("\tif c == nil { return }")
	w.Line("\tc.headersMu.Lock()")
	w.Line("\tdefer c.headersMu.Unlock()")
	w.Line("\tdelete(c.headers, key)")
	w.Line("}")
	w.Blank()
	w.Line("func SetClientAuthorization(c *Client, token string) { SetClientHeader(c, \"Authorization\", \"Bearer \"+token) }")
	w.Line("func GetClientAuthorization(c *Client) string        { return GetClientHeader(c, \"Authorization\") }")
	w.Line("func RemoveClientAuthorization(c *Client)            { RemoveClientHeader(c, \"Authorization\") }")
	w.Line("func IsClientAuthorized(c *Client) bool              { return GetClientAuthorization(c) != \"\" }")
	w.Blank()
	w.Line("func isTimeout(err error) bool {")
	w.Line("\tif err == nil { return false }")
	w.Line("\tif err == context.DeadlineExceeded { return true }")
	w.Line("\tif netErr, ok := err.(interface{ Timeout() bool }); ok && netErr.Timeout() { return true }")
	w.Line("\treturn false")
	w.Line("}")
	w.Blank()
	w.Line("func doClient(c *Client, ctx context.Context, path string, body []byte) ([]byte, RpcErrCode) {")
	w.Line("\tif c == nil { return nil, RpcNoConn }")
	w.Line("\tif ctx == nil { ctx = context.Background() }")
	w.Line("\thttpClient := c.HTTP")
	w.Line("\tif httpClient == nil { httpClient = &http.Client{} }")
	w.Line("\tmaxRetries := c.Retries")
	w.Line("\tif maxRetries < 0 { maxRetries = 0 }")
	w.Line("\tmaxRespBytes := c.MaxRespBytes")
	w.Line("\tif maxRespBytes <= 0 { maxRespBytes = defaultMaxRespBytes }")
	w.Line("\tvar resp *http.Response")
	w.Line("\tvar err error")
	w.Blank()
	w.Line("\tfor i := 0; i <= maxRetries; i++ {")
	w.Line("\t\tif i > 0 {")
	w.Line("\t\t\ttimer := time.NewTimer(time.Duration(i) * time.Second)")
	w.Line("\t\tselect {")
	w.Line("\t\tcase <-ctx.Done():")
	w.Line("\t\t\ttimer.Stop()")
	w.Line("\t\t\treturn nil, RpcTimeout")
	w.Line("\t\tcase <-timer.C:")
	w.Line("\t\t}")
	w.Line("\t}")
	w.Blank()
	w.Line("\t\treqCtx := ctx")
	w.Line("\t\tcancel := func() {}")
	w.Line("\t\tif c.Timeout > 0 {")
	w.Line("\t\t\treqCtx, cancel = context.WithTimeout(ctx, c.Timeout)")
	w.Line("\t\t}")
	w.Line("\t\treq, reqErr := http.NewRequestWithContext(reqCtx, \"POST\", c.BaseURL+path, bytes.NewReader(body))")
	w.Line("\t\tif reqErr != nil {")
	w.Line("\t\t\tcancel()")
	w.Line("\t\t\treturn nil, RpcNoConn")
	w.Line("\t\t}")
	w.Blank()
	w.Line("\t\tc.headersMu.RLock()")
	w.Line("\t\tfor k, v := range c.headers {")
	w.Line("\t\t\treq.Header.Set(k, v)")
	w.Line("\t\t}")
	w.Line("\t\tc.headersMu.RUnlock()")
	w.Line("\t\tif req.Header.Get(\"Content-Type\") == \"\" {")
	w.Line("\t\t\treq.Header.Set(\"Content-Type\", \"application/octet-stream\")")
	w.Line("\t\t}")
	w.Blank()
	w.Line("\t\tresp, err = httpClient.Do(req)")
	w.Line("\t\tcancel()")
	w.Line("\t\tif err != nil {")
	w.Line("\t\t\tif isTimeout(err) && i < maxRetries {")
	w.Line("\t\t\t\tcontinue")
	w.Line("\t\t\t}")
	w.Line("\t\t\tif isTimeout(err) {")
	w.Line("\t\t\t\treturn nil, RpcTimeout")
	w.Line("\t\t\t}")
	w.Line("\t\t\treturn nil, RpcNoConn")
	w.Line("\t\t}")
	w.Blank()
	w.Line("\t\tif resp.StatusCode == http.StatusRequestTimeout && i < maxRetries {")
	w.Line("\t\t\tresp.Body.Close()")
	w.Line("\t\t\tcontinue")
	w.Line("\t\t}")
	w.Line("\t\tbreak")
	w.Line("\t}")
	w.Line("\tif resp == nil { return nil, RpcNoConn }")
	w.Line("\tdefer resp.Body.Close()")
	w.Blank()
	w.Line("\tif resp.StatusCode != http.StatusOK {")
	w.Line("\t\treturn nil, RpcErrCode(resp.StatusCode)")
	w.Line("\t}")
	w.Blank()
	w.Line("\treadLimit := maxRespBytes")
	w.Line("\tif readLimit < (1<<63 - 1) { readLimit++ }")
	w.Line("\tb, readErr := io.ReadAll(io.LimitReader(resp.Body, readLimit))")
	w.Line("\tif readErr != nil {")
	w.Line("\t\treturn nil, RpcRespErr")
	w.Line("\t}")
	w.Line("\tif readLimit > maxRespBytes && int64(len(b)) > maxRespBytes { return nil, RpcRespErr }")
	w.Line("\treturn b, RpcOk")
	w.Line("}")
	for _, api := range apis {
		callName := PascalCase(api.Name)
		w.Blank()
		w.WriteLineCommentWithHead("// ", "Call"+callName+" ", api.Note)
		w.Linef("func Call%s(c *Client, ctx context.Context%s) (%s) {", callName, g.goRPCArgs(api.Args), g.goRPCReturn(api.Result))
		w.Line("\tvar buf bytes.Buffer")
		if len(api.Args) > 0 {
			for _, arg := range api.Args {
				argName := CamelCase(arg.Name)
				if arg.Type.Kind == TplKindStruct && !arg.Type.IsList {
					w.Linef("\tif %s == nil {", argName)
					w.Linef("\t\treturn %s", g.goRPCReqErrReturn(api.Result))
					w.Line("\t}")
					w.Linef("\tif err := Validate%s(%s); err != nil {", PascalCase(arg.Type.Name), argName)
					w.Linef("\t\treturn %s", g.goRPCReqErrReturn(api.Result))
					w.Line("\t}")
				}
				if arg.Type.Kind == TplKindStruct && arg.Type.IsList {
					w.Linef("\tif err := Validate%sList((%sList)(%s)); err != nil {", PascalCase(arg.Type.Name), PascalCase(arg.Type.Name), argName)
					w.Linef("\t\treturn %s", g.goRPCReqErrReturn(api.Result))
					w.Line("\t}")
				}
				if arg.Type.Kind == TplKindEnum {
					typeName := PascalCase(arg.Type.Name)
					if arg.Type.IsList {
						w.Linef("\tif !Is%sList((%sList)(%s)) {", typeName, typeName, argName)
						w.Linef("\t\treturn %s", g.goRPCReqErrReturn(api.Result))
						w.Line("\t}")
					} else {
						w.Linef("\tif !Is%s(%s) {", typeName, argName)
						w.Linef("\t\treturn %s", g.goRPCReqErrReturn(api.Result))
						w.Line("\t}")
					}
				}
			}
			w.Line("\tbodySize := 0")
			for _, arg := range api.Args {
				argName := CamelCase(arg.Name)
				g.writeGoSizeAccumulateLines(&w, arg.Type, argName, "bodySize", "\t", fmt.Sprintf("return %s", g.goRPCReqErrReturn(api.Result)))
			}
			w.Line("\tbuf.Grow(bodySize)")
			for _, arg := range api.Args {
				argName := CamelCase(arg.Name)
				w.Linef("\tif err := %s; err != nil {", g.goRPCSetterCall(arg.Type, "&buf", argName))
				w.Linef("\t\treturn %s", g.goRPCReqErrReturn(api.Result))
				w.Line("\t}")
			}
		}
		w.Blank()
		w.Linef("\tbody, status := doClient(c, ctx, \"/%s\", buf.Bytes())", api.Name)
		w.Line("\tif status != RpcOk {")
		w.Linef("\t\treturn %s", g.goRPCStatusReturn(api.Result))
		w.Line("\t}")
		w.Blank()
		if api.Result.Name == "nil" {
			w.Line("\tif len(body) != 0 { return RpcRespErr }")
			w.Line("\treturn status")
			w.Line("}")
			continue
			}
			w.Line("\trespBuf := bytes.NewBuffer(body)")
			w.Linef("\tval, err := %s", g.goReadExpr(api.Result, "respBuf"))
			w.Line("\tif err != nil {")
			w.Line("\t\treturn result, RpcRespErr")
			w.Line("\t}")
			w.Line("\tresult = val")
			w.Line("\tif respBuf.Len() != 0 { return result, RpcRespErr }")
		if api.Result.Kind == TplKindEnum && api.Result.IsList {
			typeName := PascalCase(api.Result.Name)
			w.Linef("\tif !Is%sList((%sList)(result)) { return result, RpcRespErr }", typeName, typeName)
		} else if api.Result.Kind == TplKindEnum {
			w.Linef("\tif !Is%s(result) { return result, RpcRespErr }", PascalCase(api.Result.Name))
		}
		w.Line("\treturn result, status")
		w.Line("}")
	}
	return w.String()
}

func (g *GoGenerator) goReadExpr(t TplType, bufVar string) string {
	name := PascalCase(t.Name)
	if t.Kind == TplKindStruct && !t.IsList {
		return fmt.Sprintf("Read%s(%s)", name, bufVar)
	}
	if t.IsList {
		return fmt.Sprintf("Get%sList(%s)", name, bufVar)
	}
	return fmt.Sprintf("Get%s(%s)", name, bufVar)
}

func (g *GoGenerator) goPresenceCheck(t TplType, ref string) string {
	if t.IsList {
		return fmt.Sprintf("len(%s) > 0", ref)
	}
	if t.Kind == TplKindStruct {
		return fmt.Sprintf("%s != nil", ref)
	}
	return fmt.Sprintf("%s != %s", ref, g.getGoValue(t.Name))
}

func (g *GoGenerator) goSetCall(t TplType, bufVar, ref string) string {
	name := PascalCase(t.Name)
	if t.IsList && t.Kind != TplKindBase {
		return fmt.Sprintf("Set%sList(%s, (%sList)(%s))", name, bufVar, name, ref)
	}
	if t.IsList {
		return fmt.Sprintf("Set%sList(%s, %s)", name, bufVar, ref)
	}
	return fmt.Sprintf("Set%s(%s, %s)", name, bufVar, ref)
}

func (g *GoGenerator) goEqExpr(t TplType, left, right string) string {
	name := PascalCase(t.Name)
	switch t.Kind {
	case TplKindBase:
		if t.IsList {
			return fmt.Sprintf("Eq%sList(%s, %s)", name, left, right)
		}
		return fmt.Sprintf("Eq%s(%s, %s)", name, left, right)
	case TplKindEnum:
		if t.IsList {
			return fmt.Sprintf("Eq%sList(%s, %s)", name, left, right)
		}
		return fmt.Sprintf("(%s == %s)", left, right)
	case TplKindStruct:
		if t.IsList {
			return fmt.Sprintf("Eq%sList((%sList)(%s), (%sList)(%s))", name, name, left, name, right)
		}
		return fmt.Sprintf("Eq%s(%s, %s)", name, left, right)
	}
	return fmt.Sprintf("%s == %s", left, right)
}

func (g *GoGenerator) goStubArgs(args []TplApiArg) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, fmt.Sprintf(", %s %s", arg.Name, g.goHandlerArgType(arg.Type)))
	}
	return strings.Join(parts, "")
}

func (g *GoGenerator) goStubResultType(t TplType) string {
	if t.IsList && t.Kind == TplKindEnum {
		return PascalCase(t.Name) + "List"
	}
	return g.getGoLogicType(t)
}

func (g *GoGenerator) goDefaultReturn(t TplType) string {
	if t.Kind == TplKindStruct || t.IsList {
		return "nil"
	}
	return g.getGoValue(t.Name)
}

func (g *GoGenerator) goHandlerArgType(t TplType) string {
	if t.IsList && t.Kind == TplKindEnum {
		return PascalCase(t.Name) + "List"
	}
	return g.getGoLogicType(t)
}

func (g *GoGenerator) goHandlerCallArgs(args []TplApiArg) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, arg.Name)
	}
	return ", " + strings.Join(parts, ", ")
}

func (g *GoGenerator) goRPCArgs(args []TplApiArg) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		argType := g.getGoLogicType(arg.Type)
		if arg.Type.Kind == TplKindEnum && arg.Type.IsList {
			argType = PascalCase(arg.Type.Name) + "List"
		}
		parts = append(parts, fmt.Sprintf(", %s %s", CamelCase(arg.Name), argType))
	}
	return strings.Join(parts, "")
}

func (g *GoGenerator) goRPCReturn(result TplType) string {
	if result.Name == "nil" {
		return "errCode RpcErrCode"
	}
	return fmt.Sprintf("result %s, errCode RpcErrCode", g.getGoLogicType(result))
}

func (g *GoGenerator) goRPCReqErrReturn(result TplType) string {
	if result.Name == "nil" {
		return "RpcReqErr"
	}
	return "result, RpcReqErr"
}

func (g *GoGenerator) goRPCStatusReturn(result TplType) string {
	if result.Name == "nil" {
		return "status"
	}
	return "result, status"
}

func (g *GoGenerator) goRPCSetters(args []TplApiArg) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		name := CamelCase(arg.Name)
		setCall := ""
		switch arg.Type.Kind {
		case TplKindBase:
			setCall = fmt.Sprintf("Set%s%s(buf, %s)", PascalCase(arg.Type.Name), goListSuffix(arg.Type), name)
		case TplKindEnum:
			if arg.Type.IsList {
				setCall = fmt.Sprintf("Set%sList(buf, (%sList)(%s))", PascalCase(arg.Type.Name), PascalCase(arg.Type.Name), name)
			} else {
				setCall = fmt.Sprintf("Set%s(buf, %s)", PascalCase(arg.Type.Name), name)
			}
		case TplKindStruct:
			if arg.Type.IsList {
				setCall = fmt.Sprintf("Set%sList(buf, (%sList)(%s))", PascalCase(arg.Type.Name), PascalCase(arg.Type.Name), name)
			} else {
				setCall = fmt.Sprintf("Set%s(buf, %s)", PascalCase(arg.Type.Name), name)
			}
		}
		parts = append(parts, fmt.Sprintf(", func(buf *bytes.Buffer) error {\n\t\treturn %s\n\t}", setCall))
	}
	return strings.Join(parts, "")
}

func goBitmaskSize(fieldCount int) int {
	if fieldCount <= 0 {
		return 0
	}
	return (fieldCount + 7) / 8
}

func goBaseEncodedWidth(name string) (int, bool) {
	switch name {
	case "bool", "i8", "u8":
		return 1, true
	case "i16", "u16":
		return 2, true
	case "i32", "u32", "f32":
		return 4, true
	case "i64", "u64", "f64":
		return 8, true
	default:
		return 0, false
	}
}

func (g *GoGenerator) goEncodedSizeExpr(t TplType, ref string) (int, string) {
	if width, ok := goBaseEncodedWidth(t.Name); ok && !t.IsList && t.Kind == TplKindBase {
		return width, ""
	}
	if t.Kind == TplKindEnum && !t.IsList {
		return 1, ""
	}
	name := PascalCase(t.Name)
	switch {
	case t.Kind == TplKindBase && t.IsList:
		switch t.Name {
		case "bool":
			return 0, fmt.Sprintf("sizeBoolList(%s)", ref)
		case "text":
			return 0, fmt.Sprintf("sizeTextList(%s)", ref)
		case "bin":
			return 0, fmt.Sprintf("sizeBinList(%s)", ref)
		default:
			width, _ := goBaseEncodedWidth(t.Name)
			return 0, fmt.Sprintf("sizeFixedList(len(%s), %d)", ref, width)
		}
	case t.Kind == TplKindBase:
		switch t.Name {
		case "text":
			return 0, fmt.Sprintf("sizeText(%s)", ref)
		case "bin":
			return 0, fmt.Sprintf("sizeBin(%s)", ref)
		}
	case t.Kind == TplKindEnum && t.IsList:
		return 0, fmt.Sprintf("sizeFixedList(len(%s), 1)", ref)
	case t.Kind == TplKindStruct && t.IsList:
		return 0, fmt.Sprintf("size%sList((%sList)(%s))", name, name, ref)
	case t.Kind == TplKindStruct:
		return 0, fmt.Sprintf("size%s(%s)", name, ref)
	}
	return 0, ""
}

func (g *GoGenerator) writeGoSizeAddLines(w *sourceWriter, ownerName, fieldName string, t TplType, ref, targetVar, indent string) {
	if width, expr := g.goEncodedSizeExpr(t, ref); width > 0 {
		w.Linef("%s%s += %d", indent, targetVar, width)
		return
	} else if expr != "" {
		w.Linef("%sfieldSize, err := %s", indent, expr)
		w.Linef("%sif err != nil { return 0, fmt.Errorf(\"%s %s: %%w\", err) }", indent, ownerName, fieldName)
		w.Linef("%s%s += fieldSize", indent, targetVar)
	}
}

func (g *GoGenerator) writeGoSizeAccumulateLines(w *sourceWriter, t TplType, ref, targetVar, indent, onError string) {
	if width, expr := g.goEncodedSizeExpr(t, ref); width > 0 {
		w.Linef("%s%s += %d", indent, targetVar, width)
		return
	} else if expr != "" {
		w.Linef("%s{", indent)
		w.Linef("%s\tfieldSize, err := %s", indent, expr)
		w.Linef("%s\tif err != nil { %s }", indent, onError)
		w.Linef("%s\t%s += fieldSize", indent, targetVar)
		w.Linef("%s}", indent)
	}
}

func (g *GoGenerator) goRPCSetterCall(t TplType, bufVar, ref string) string {
	return g.goSetCall(t, bufVar, ref)
}

func goListSuffix(t TplType) string {
	if t.IsList {
		return "List"
	}
	return ""
}

func orderedGroupKeys(groups map[string][]TplApi) []string {
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[j] < keys[i] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}
