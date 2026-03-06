package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestTsRPCUsesSetterClosuresForAllArgKinds(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewTsGenerator(Config{TsDir: tmpDir})

	schema := &TplSchema{
		Structs: []TplStruct{
			{Name: "sim_info"},
		},
		Apis: []TplApi{
			{
				Name: "foo.test",
				Args: []TplApiArg{
					{Name: "nums", Type: TplType{Name: "u16", Kind: TplKindBase, IsList: true}},
					{Name: "id", Type: TplType{Name: "u64", Kind: TplKindBase}},
					{Name: "score", Type: TplType{Name: "i64", Kind: TplKindBase}},
					{Name: "infos", Type: TplType{Name: "sim_info", Kind: TplKindStruct, IsList: true}},
					{Name: "ok", Type: TplType{Name: "bool", Kind: TplKindBase}},
				},
				Result: TplType{Name: "nil", Kind: TplKindNil},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate ts failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc.ts"))
	if err != nil {
		t.Fatalf("read rpc.ts failed: %v", err)
	}
	text := string(content)

	assertContains(t, text, "_.setU16List(buf, nums)")
	assertContains(t, text, "_.setU64(buf, id)")
	assertContains(t, text, "_.setI64(buf, score)")
	assertContains(t, text, "_.setSimInfoList(buf, infos as any)")
	assertContains(t, text, "_.setBool(buf, ok)")
	assertContains(t, text, "_.setAll(buf, (buf: _.Buffer) => _.setU16List(buf, nums), (buf: _.Buffer) => _.setU64(buf, id), (buf: _.Buffer) => _.setI64(buf, score), (buf: _.Buffer) => _.setSimInfoList(buf, infos as any), (buf: _.Buffer) => _.setBool(buf, ok))")
	assertContains(t, text, "const defaultMaxRespBytes = 4 * 1024 * 1024;")
	assertContains(t, text, "const maxSafeRespBytes = Number.MAX_SAFE_INTEGER;")
	assertContains(t, text, "const maxTimeoutMs = 2147483647;")
	assertContains(t, text, "maxRespBytes?: number;")
	assertContains(t, text, "private maxRespBytes: number;")
	assertContains(t, text, "this.config = { ...config, host: config.host.replace(/\\/+$/, \"\") };")
	assertContains(t, text, "const cfgTimeout = config.timeout;")
	assertContains(t, text, "if (cfgTimeout !== undefined && Number.isFinite(cfgTimeout) && cfgTimeout >= 0) {")
	assertContains(t, text, "this.timeout = Math.min(Math.floor(cfgTimeout), maxTimeoutMs);")
	assertNotContains(t, text, "this.timeout = Math.floor(cfgTimeout);")
	assertContains(t, text, "this.timeout = 5000;")
	assertNotContains(t, text, "this.timeout = config.timeout || 5000;")
	assertContains(t, text, "const cfgRetries = config.retries;")
	assertContains(t, text, "if (cfgRetries !== undefined && Number.isFinite(cfgRetries) && cfgRetries >= 0) {")
	assertContains(t, text, "this.retries = Math.floor(cfgRetries);")
	assertContains(t, text, "} else {")
	assertContains(t, text, "this.retries = 3;")
	assertContains(t, text, "const cfgMaxRespBytes = config.maxRespBytes;")
	assertContains(t, text, "if (cfgMaxRespBytes !== undefined && Number.isFinite(cfgMaxRespBytes) && cfgMaxRespBytes > 0) {")
	assertContains(t, text, "this.maxRespBytes = Math.min(Math.floor(cfgMaxRespBytes), maxSafeRespBytes);")
	assertContains(t, text, "} else {")
	assertContains(t, text, "this.maxRespBytes = defaultMaxRespBytes;")
	assertNotContains(t, text, "if (this.retries < 0) this.retries = 0;")
	assertContains(t, text, "let timeoutId: ReturnType<typeof setTimeout> | null = null;")
	assertContains(t, text, "if (this.timeout > 0) timeoutId = setTimeout(() => controller.abort(), this.timeout);")
	assertContains(t, text, "if (timeoutId !== null) clearTimeout(timeoutId);")
	assertContains(t, text, "const contentLength = res.headers.get(\"content-length\");")
	assertContains(t, text, "if (bytes.byteLength > this.maxRespBytes) return [null, RpcErrCode.RespErr];")
	assertContains(t, text, "if (bytes.byteLength !== 0) return RpcErrCode.RespErr;")
	assertContains(t, text, "import * as _ from \"./_\"")
	assertNotContains(t, text, "_.setU16List(nums)")
	assertNotContains(t, text, "from \"./_.ts\"")

	indexContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "_.ts"))
	if err != nil {
		t.Fatalf("read _.ts failed: %v", err)
	}
	indexText := string(indexContent)
	assertContains(t, indexText, "export * from \"./type\"")
	assertNotContains(t, indexText, "export * from \"./type.ts\"")

	smokeContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc_smoke.test.ts"))
	if err != nil {
		t.Fatalf("read rpc_smoke.test.ts failed: %v", err)
	}
	smokeText := string(smokeContent)
	assertContains(t, smokeText, "const baseUrl = process.env.SB_BASE_URL || \"http://127.0.0.1:18080\";")
	assertContains(t, smokeText, "test(\"client construction\", () => {")
	assertContains(t, smokeText, "expect(typeof (client as any).fooTest).toBe(\"function\");")
}

func TestRuntimeLengthPrefixesUpdated(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := Config{
		GoDir: filepath.Join(tmpDir, "go"),
		TsDir: filepath.Join(tmpDir, "ts"),
	}

	goGen := NewGoGenerator(cfg)
	if err := goGen.Generate(&TplSchema{}); err != nil {
		t.Fatalf("generate go failed: %v", err)
	}

	tsGen := NewTsGenerator(cfg)
	if err := tsGen.Generate(&TplSchema{}); err != nil {
		t.Fatalf("generate ts failed: %v", err)
	}

	goType, err := os.ReadFile(filepath.Join(cfg.GoDir, "sb", "type.go"))
	if err != nil {
		t.Fatalf("read go type.go failed: %v", err)
	}
	goText := string(goType)
	assertContains(t, goText, "count, err := GetU16(buf)")
	assertContains(t, goText, "if len(list) > 65535 { return fmt.Errorf(\"list length exceeds uint16 max\") }")
	assertContains(t, goText, "if err := SetU16(buf, uint16(len(list))); err != nil { return err }")
	assertContains(t, goText, "l, err := GetU32(buf)")
	assertContains(t, goText, "func GetBinView(buf *bytes.Buffer) ([]byte, error) {")
	assertContains(t, goText, "return buf.Next(int(l)), nil")
	assertContains(t, goText, "res, err := GetBinView(buf)")
	assertContains(t, goText, "return bytes.Clone(res), nil")
	assertContains(t, goText, "if uint64(len(v)) > uint64(^uint32(0)) { return fmt.Errorf(\"bin length exceeds uint32 max\") }")
	assertContains(t, goText, "if err := SetU32(buf, uint32(len(v))); err != nil { return err }")
	assertContains(t, goText, "if len(data) > 65535 { return fmt.Errorf(\"text length exceeds uint16 max\") }")
	assertContains(t, goText, "if err := SetU16(buf, uint16(len(data))); err != nil { return err }")

	tsType, err := os.ReadFile(filepath.Join(cfg.TsDir, "sb", "type.ts"))
	if err != nil {
		t.Fatalf("read ts type.ts failed: %v", err)
	}
	tsText := string(tsType)
	assertContains(t, tsText, "const [count, err] = getU16(buf);")
	assertContains(t, tsText, "if (list.length > 65535) return new Error(`list length ${list.length} exceeds u16 max`);")
	assertContains(t, tsText, "const err = setU16(buf, list.length);")
	assertContains(t, tsText, "const [len, err] = getU32(buf);")
	assertContains(t, tsText, "if (value.length > 4294967295) return new Error(`bin length ${value.length} exceeds u32 max`);")
	assertContains(t, tsText, "const err = setU32(buf, value.length);")
	assertContains(t, tsText, "if (data.length > 65535) return new Error(`text length ${data.length} exceeds u16 max`);")
	assertContains(t, tsText, "const err = setU16(buf, data.length);")
}

func TestTsEnumValidationGeneratedAndWired(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewTsGenerator(Config{TsDir: tmpDir})

	schema := &TplSchema{
		Enums: []TplEnum{
			{
				Name: "sim_operator",
				Children: []TplEnumChild{
					{ID: 2, Name: "zz"},
					{ID: 3, Name: "lt"},
				},
			},
		},
		Structs: []TplStruct{
			{
				Name: "sim",
				Fields: []TplStructField{
					{Name: "operator", Type: TplType{Name: "sim_operator", Kind: TplKindEnum}},
					{Name: "operators", Type: TplType{Name: "sim_operator", Kind: TplKindEnum, IsList: true}},
				},
			},
		},
		Apis: []TplApi{
			{
				Name: "sim.set",
				Args: []TplApiArg{
					{Name: "operator", Type: TplType{Name: "sim_operator", Kind: TplKindEnum}},
					{Name: "operators", Type: TplType{Name: "sim_operator", Kind: TplKindEnum, IsList: true}},
				},
				Result: TplType{Name: "sim_operator", Kind: TplKindEnum},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate ts failed: %v", err)
	}

	enumContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "enum.ts"))
	if err != nil {
		t.Fatalf("read enum.ts failed: %v", err)
	}
	enumText := string(enumContent)
	assertContains(t, enumText, "export const IsSimOperator = (v: SimOperator): boolean => {")
	assertContains(t, enumText, "export const IsSimOperatorList = (v: SimOperator[]): boolean => {")

	structContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "struct_sim.ts"))
	if err != nil {
		t.Fatalf("read struct_sim.ts failed: %v", err)
	}
	structText := string(structContent)
	assertContains(t, structText, "if (!_.IsSimOperatorList(s.operators as any)) return [s, new Error(\"get Sim Operators: invalid enum value\")];")
	assertContains(t, structText, "if (!_.IsSimOperator(s.operator as any)) return [s, new Error(\"get Sim Operator: invalid enum value\")];")
	assertNotContains(t, structText, "if ((s.operator as any) !== 0 && !_.IsSimOperator(s.operator as any)) return [s, new Error(\"get Sim Operator: invalid enum value\")];")
	assertContains(t, structText, "if (!_.IsSimOperatorList(s.operators as any)) return new Error(\"set Sim Operators: invalid enum value\");")
	assertContains(t, structText, "if (!_.IsSimOperator(s.operator as any)) return new Error(\"set Sim Operator: invalid enum value\");")

	rpcContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc.ts"))
	if err != nil {
		t.Fatalf("read rpc.ts failed: %v", err)
	}
	rpcText := string(rpcContent)
	assertContains(t, rpcText, "if (!_.IsSimOperator(operator as any))")
	assertNotContains(t, rpcText, "if ((operator as any) !== 0 && !_.IsSimOperator(operator as any))")
	assertContains(t, rpcText, "if (!_.IsSimOperatorList(operators as any))")
	assertContains(t, rpcText, "const respBuf = new _.Buffer(bytes);")
	assertContains(t, rpcText, "const [result, err] = _.getU8(respBuf);")
	assertNotContains(t, rpcText, "const respBuf = new _.Buffer(bytes);const [result, err] = _.getU8(respBuf);")
	assertContains(t, rpcText, "if (err === null && !_.IsSimOperator(result as any))")
	assertContains(t, rpcText, "if (respBuf.len !== 0) return [0 as _.SimOperator, RpcErrCode.RespErr];")
	assertNotContains(t, rpcText, "const [result, err] = _.getU8(new _.Buffer(bytes));")
	assertNotContains(t, rpcText, "if (err === null && (result as any) !== 0 && !_.IsSimOperator(result as any))")
}

func TestTsStructNestedFieldUsesNullablePresenceSemantics(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewTsGenerator(Config{TsDir: tmpDir})

	schema := &TplSchema{
		Structs: []TplStruct{
			{
				Name: "child",
				Fields: []TplStructField{
					{Name: "id", Type: TplType{Name: "u8", Kind: TplKindBase}},
				},
			},
			{
				Name: "parent",
				Fields: []TplStructField{
					{Name: "child", Type: TplType{Name: "child", Kind: TplKindStruct}},
					{Name: "children", Type: TplType{Name: "child", Kind: TplKindStruct, IsList: true}},
				},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate ts failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "sb", "struct_parent.ts"))
	if err != nil {
		t.Fatalf("read struct_parent.ts failed: %v", err)
	}
	text := string(content)

	assertContains(t, text, "child: _.Child | null;")
	assertContains(t, text, "children: _.Child[];")
	assertContains(t, text, "child: null,")
	assertNotContains(t, text, "child: _.newChild()")
	assertContains(t, text, "if ((a.child === null) !== (b.child === null)) return false;")
	assertContains(t, text, "if (a.child !== null && b.child !== null && !_.eqChild(a.child, b.child)) return false;")
	assertNotContains(t, text, "if (!_.eqChild(a.child, b.child)) return false;")
	assertContains(t, text, "if (s.child !== null && s.child !== undefined) {")
	assertContains(t, text, "const err = _.setChild(body, s.child);")
}

func TestGoRPCHandlesReadAllErrorAndUsesClientTimeout(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGoGenerator(Config{GoDir: tmpDir})

	schema := &TplSchema{
		Apis: []TplApi{
			{
				Name:   "foo.ping",
				Result: TplType{Name: "nil", Kind: TplKindNil},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate go failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc.go"))
	if err != nil {
		t.Fatalf("read rpc.go failed: %v", err)
	}
	text := string(content)

	assertContains(t, text, "context.WithTimeout(ctx, c.Timeout)")
	assertContains(t, text, "\"strings\"")
	assertContains(t, text, "baseURL = strings.TrimRight(baseURL, \"/\")")
	assertContains(t, text, "defaultMaxRespBytes int64 = 4 * 1024 * 1024")
	assertContains(t, text, "MaxRespBytes int64")
	assertContains(t, text, "headersMu sync.RWMutex")
	assertContains(t, text, "c.headersMu.Lock()")
	assertContains(t, text, "c.headersMu.RLock()")
	assertContains(t, text, "headers := make(map[string]string, len(c.headers))")
	assertContains(t, text, "if req.Header.Get(\"Content-Type\") == \"\" {")
	assertContains(t, text, "req.Header.Set(\"Content-Type\", \"application/octet-stream\")")
	assertContains(t, text, "maxRetries := c.Retries")
	assertContains(t, text, "if maxRetries < 0 { maxRetries = 0 }")
	assertContains(t, text, "maxRespBytes := c.MaxRespBytes")
	assertContains(t, text, "if maxRespBytes <= 0 { maxRespBytes = defaultMaxRespBytes }")
	assertContains(t, text, "if ctx == nil { ctx = context.Background() }")
	assertContains(t, text, "httpClient := c.HTTP")
	assertContains(t, text, "if httpClient == nil { httpClient = &http.Client{} }")
	assertNotContains(t, text, "if c.HTTP == nil { c.HTTP = &http.Client{} }")
	assertContains(t, text, "resp, err = httpClient.Do(req)")
	assertContains(t, text, "for i := 0; i <= maxRetries; i++")
	assertContains(t, text, "if resp == nil { return nil, RpcNoConn }")
	assertContains(t, text, "b, readErr := io.ReadAll(io.LimitReader(resp.Body, readLimit))")
	assertContains(t, text, "if readLimit > maxRespBytes && int64(len(b)) > maxRespBytes { return nil, RpcRespErr }")
	assertContains(t, text, "if readErr != nil")
	assertContains(t, text, "func CallFooPing(c *Client, ctx context.Context)")
	assertContains(t, text, "body, status := doClient(c, ctx, \"/foo.ping\", buf.Bytes())")
	assertNotContains(t, text, "_, status := doClient(c, ctx, \"/foo.ping\", buf.Bytes())")
	assertContains(t, text, "if len(body) != 0 { return RpcRespErr }")
	assertContains(t, text, "func doClient(c *Client, ctx context.Context, path string, body []byte)")
	assertNotContains(t, text, "func (c *Client)")
}

func TestGoRPCValidatesEnumResults(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGoGenerator(Config{GoDir: tmpDir})

	schema := &TplSchema{
		Enums: []TplEnum{
			{
				Name: "sim_operator",
				Children: []TplEnumChild{
					{ID: 2, Name: "zz"},
					{ID: 3, Name: "lt"},
				},
			},
		},
		Apis: []TplApi{
			{
				Name:   "sim.get",
				Result: TplType{Name: "sim_operator", Kind: TplKindEnum},
			},
			{
				Name:   "sim.list",
				Result: TplType{Name: "sim_operator", Kind: TplKindEnum, IsList: true},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate go failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc.go"))
	if err != nil {
		t.Fatalf("read rpc.go failed: %v", err)
	}
	text := string(content)

	assertContains(t, text, "respBuf := bytes.NewBuffer(body)")
	assertContains(t, text, "if respBuf.Len() != 0 { return result, RpcRespErr }")
	assertContains(t, text, "if !IsSimOperator(result) { return result, RpcRespErr }")
	assertContains(t, text, "if !IsSimOperatorList((SimOperatorList)(result)) { return result, RpcRespErr }")
}

func TestGoRPCValidatesEnumArgsBeforeRequest(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGoGenerator(Config{GoDir: tmpDir})

	schema := &TplSchema{
		Enums: []TplEnum{
			{
				Name: "sim_operator",
				Children: []TplEnumChild{
					{ID: 2, Name: "zz"},
					{ID: 3, Name: "lt"},
				},
			},
		},
		Apis: []TplApi{
			{
				Name: "sim.set",
				Args: []TplApiArg{
					{Name: "operator", Type: TplType{Name: "sim_operator", Kind: TplKindEnum}},
					{Name: "operators", Type: TplType{Name: "sim_operator", Kind: TplKindEnum, IsList: true}},
				},
				Result: TplType{Name: "nil", Kind: TplKindNil},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate go failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc.go"))
	if err != nil {
		t.Fatalf("read rpc.go failed: %v", err)
	}
	text := string(content)

	assertContains(t, text, "func CallSimSet(c *Client, ctx context.Context, operator SimOperator, operators SimOperatorList) (errCode RpcErrCode)")
	assertContains(t, text, "if !IsSimOperator(operator) {")
	assertContains(t, text, "if !IsSimOperatorList((SimOperatorList)(operators)) {")
	assertContains(t, text, "return RpcReqErr")
}

func TestGoRPCRejectsNilStructArgsBeforeRequest(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGoGenerator(Config{GoDir: tmpDir})

	schema := &TplSchema{
		Structs: []TplStruct{
			{
				Name: "sim_info",
				Fields: []TplStructField{
					{Name: "id", Type: TplType{Name: "u32", Kind: TplKindBase}},
				},
			},
		},
		Apis: []TplApi{
			{
				Name: "sim.push",
				Args: []TplApiArg{
					{Name: "info", Type: TplType{Name: "sim_info", Kind: TplKindStruct}},
				},
				Result: TplType{Name: "nil", Kind: TplKindNil},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate go failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc.go"))
	if err != nil {
		t.Fatalf("read rpc.go failed: %v", err)
	}
	text := string(content)

	assertContains(t, text, "func CallSimPush(c *Client, ctx context.Context, info *SimInfo) (errCode RpcErrCode)")
	assertContains(t, text, "if info == nil {")
	assertContains(t, text, "return RpcReqErr")
}

func TestGoEnumTemplateOmitsImportsWhenNoEnums(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGoGenerator(Config{GoDir: tmpDir})

	schema := &TplSchema{}
	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate go failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "sb", "enum.go"))
	if err != nil {
		t.Fatalf("read enum.go failed: %v", err)
	}
	text := string(content)

	assertContains(t, text, "package sb")
	assertNotContains(t, text, "import (")
	assertNotContains(t, text, "\"unsafe\"")
	assertNotContains(t, text, "\"slices\"")
	assertNotContains(t, text, "\"bytes\"")
}

func TestGoApiHandlersValidateEnumAndStructArgs(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGoGenerator(Config{GoDir: tmpDir})

	schema := &TplSchema{
		Enums: []TplEnum{
			{
				Name: "sim_operator",
				Children: []TplEnumChild{
					{ID: 2, Name: "zz"},
					{ID: 3, Name: "lt"},
				},
			},
		},
		Structs: []TplStruct{
			{
				Name: "sim_info",
				Fields: []TplStructField{
					{Name: "id", Type: TplType{Name: "u32", Kind: TplKindBase}},
				},
			},
			{
				Name: "sim",
				Fields: []TplStructField{
					{Name: "operator", Type: TplType{Name: "sim_operator", Kind: TplKindEnum}},
					{Name: "operators", Type: TplType{Name: "sim_operator", Kind: TplKindEnum, IsList: true}},
					{Name: "info", Type: TplType{Name: "sim_info", Kind: TplKindStruct}},
					{Name: "infos", Type: TplType{Name: "sim_info", Kind: TplKindStruct, IsList: true}},
				},
			},
		},
		Apis: []TplApi{
			{
				Name: "sim.set",
				Args: []TplApiArg{
					{Name: "operator", Type: TplType{Name: "sim_operator", Kind: TplKindEnum}},
					{Name: "operators", Type: TplType{Name: "sim_operator", Kind: TplKindEnum, IsList: true}},
					{Name: "sim", Type: TplType{Name: "sim", Kind: TplKindStruct}},
				},
				Result: TplType{Name: "nil", Kind: TplKindNil},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate go failed: %v", err)
	}

	enumContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "enum.go"))
	if err != nil {
		t.Fatalf("read enum.go failed: %v", err)
	}
	enumText := string(enumContent)
	assertContains(t, enumText, "func IsSimOperator(v SimOperator) bool")
	assertContains(t, enumText, "func IsSimOperatorList(v SimOperatorList) bool")
	assertContains(t, enumText, "func GetSimOperator(buf *bytes.Buffer) (SimOperator, error)")
	assertContains(t, enumText, "func SetSimOperator(buf *bytes.Buffer, v SimOperator) error")
	assertNotContains(t, enumText, "func (v SimOperatorList)")

	structContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "struct_sim.go"))
	if err != nil {
		t.Fatalf("read struct_sim.go failed: %v", err)
	}
	structText := string(structContent)
	assertContains(t, structText, "func ValidateSim(s *Sim) error")
	assertContains(t, structText, "if s.Operator != 0 && !IsSimOperator(s.Operator)")
	assertContains(t, structText, "if !IsSimOperator(s.Operator) { return fmt.Errorf(\"GetSim Operator: 非法枚举值: %d\", s.Operator) }")
	assertContains(t, structText, "if err := ValidateSim(s); err != nil")
	assertContains(t, structText, "func GetSim(buf *bytes.Buffer, s *Sim) error")
	assertContains(t, structText, "func ReadSim(buf *bytes.Buffer) (*Sim, error)")
	assertContains(t, structText, "func SetSim(buf *bytes.Buffer, s *Sim) error")
	assertContains(t, structText, "if s == nil { return fmt.Errorf(\"SetSim: nil value\") }")
	assertNotContains(t, structText, "func SetSim(buf *bytes.Buffer, s *Sim) error {\n\tif s == nil { return nil }")
	assertContains(t, structText, "func GetSimList(buf *bytes.Buffer) (SimList, error)")
	assertContains(t, structText, "func SetSimList(buf *bytes.Buffer, v SimList) error")
	assertContains(t, structText, "func EqSim(a, b *Sim) bool")
	assertContains(t, structText, "func EqSimList(a, b SimList) bool")
	assertContains(t, structText, "if !EqSimOperatorList(a.Operators, b.Operators) { return false }")
	assertContains(t, structText, "if !EqSimInfo(a.Info, b.Info) { return false }")
	assertContains(t, structText, "if !EqSimInfoList((SimInfoList)(a.Infos), (SimInfoList)(b.Infos)) { return false }")
	assertContains(t, structText, "if item == nil { return fmt.Errorf(\"SimList[%d]: nil item\", i) }")
	assertNotContains(t, structText, "if item == nil { continue }")
	assertNotContains(t, structText, "if buf.Len() == 0 { return nil }")
	assertNotContains(t, structText, "func (s *Sim) Get(")
	assertNotContains(t, structText, "func (s *Sim) Set(")
	assertNotContains(t, structText, "func (s *Sim) Eq(")
	assertNotContains(t, structText, ".Eq(")
	assertNotContains(t, structText, ".Get(")
	assertNotContains(t, structText, ".Set(")

	apiContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "api._.go"))
	if err != nil {
		t.Fatalf("read api._.go failed: %v", err)
	}
	apiText := string(apiContent)
	assertContains(t, apiText, "if !IsSimOperator(operator)")
	assertContains(t, apiText, "if !IsSimOperatorList(operators)")
	assertContains(t, apiText, "if err := ValidateSim(sim); err != nil")
	assertContains(t, apiText, "const defaultMaxReqBytes int64 = 4 * 1024 * 1024")
	assertContains(t, apiText, "func parseRequest(w http.ResponseWriter, r *http.Request, args ...Getter) bool")
	assertContains(t, apiText, "body, err := io.ReadAll(io.LimitReader(r.Body, 1))")
	assertContains(t, apiText, "if len(body) != 0 { w.WriteHeader(http.StatusBadRequest); return false }")
	assertContains(t, apiText, "body, err := io.ReadAll(io.LimitReader(r.Body, readLimit))")
	assertContains(t, apiText, "buf := bytes.NewBuffer(body)")
	assertContains(t, apiText, "if buf.Len() != 0 { w.WriteHeader(http.StatusBadRequest); return false }")
	assertContains(t, apiText, "if int64(len(body)) > defaultMaxReqBytes { w.WriteHeader(http.StatusRequestEntityTooLarge); return false }")
	assertContains(t, apiText, "if w.Header().Get(\"Content-Type\") == \"\" {")
	assertContains(t, apiText, "w.Header().Set(\"Content-Type\", \"application/octet-stream\")")
	assertContains(t, apiText, "if _, err := w.Write(buf.Bytes()); err != nil { w.WriteHeader(http.StatusInternalServerError); return }")
	assertContains(t, apiText, "func sendResponse(w http.ResponseWriter, setter Setter)")

	typeContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "type.go"))
	if err != nil {
		t.Fatalf("read type.go failed: %v", err)
	}
	typeText := string(typeContent)
	assertContains(t, typeText, "type Setter func(*bytes.Buffer) error")
	assertContains(t, typeText, "type Getter func(*bytes.Buffer) error")
	assertContains(t, typeText, "if uint64(buf.Len()) < uint64(l) { return nil, fmt.Errorf(\"not enough data\") }")
	assertNotContains(t, typeText, "if uint16(buf.Len()) < l { return nil, fmt.Errorf(\"not enough data\") }")
	assertNotContains(t, typeText, "type Serializable interface")
	assertNotContains(t, typeText, "type Deserializable interface")
	assertNotContains(t, typeText, "func (v ")
}

func TestGoApiDefaultErrReturnMatchesResultType(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGoGenerator(Config{GoDir: tmpDir})

	schema := &TplSchema{
		Structs: []TplStruct{
			{
				Name: "envelope",
				Fields: []TplStructField{
					{Name: "id", Type: TplType{Name: "u32", Kind: TplKindBase}},
				},
			},
		},
		Apis: []TplApi{
			{
				Name: "echo",
				Args: []TplApiArg{
					{Name: "env", Type: TplType{Name: "envelope", Kind: TplKindStruct}},
				},
				Result: TplType{Name: "envelope", Kind: TplKindStruct},
			},
			{
				Name:   "echo_list",
				Result: TplType{Name: "envelope", Kind: TplKindStruct, IsList: true},
			},
			{
				Name:   "get_count",
				Result: TplType{Name: "u32", Kind: TplKindBase},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate go failed: %v", err)
	}

	echoContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "api.echo.go"))
	if err != nil {
		t.Fatalf("read api.echo.go failed: %v", err)
	}
	echoText := string(echoContent)
	assertContains(t, echoText, "func echo(ctx context.Context, env *Envelope) (result *Envelope, errCode RpcErrCode) {")
	assertContains(t, echoText, "return nil, RpcRespErr")
	assertNotContains(t, echoText, "return 0, RpcRespErr")

	echoListContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "api.echo_list.go"))
	if err != nil {
		t.Fatalf("read api.echo_list.go failed: %v", err)
	}
	echoListText := string(echoListContent)
	assertContains(t, echoListText, "func echo_list(ctx context.Context) (result []*Envelope, errCode RpcErrCode) {")
	assertContains(t, echoListText, "return nil, RpcRespErr")

	countContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "api.get_count.go"))
	if err != nil {
		t.Fatalf("read api.get_count.go failed: %v", err)
	}
	countText := string(countContent)
	assertContains(t, countText, "func get_count(ctx context.Context) (result uint32, errCode RpcErrCode) {")
	assertContains(t, countText, "return 0, RpcRespErr")
}

func TestTemplateGeneratorRendersMultilineNotesSafely(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := Config{
		GoDir: filepath.Join(tmpDir, "go"),
		TsDir: filepath.Join(tmpDir, "ts"),
	}

	schema := &TplSchema{
		Structs: []TplStruct{
			{
				Name: "sim_info",
				Fields: []TplStructField{
					{
						Name: "id",
						Type: TplType{Name: "u32", Kind: TplKindBase},
						Note: "主键\n数据库ID",
					},
				},
			},
		},
		Enums: []TplEnum{
			{
				Name: "order_status",
				Note: "订单状态\nPending      待处理\nClosed       已关闭",
				Children: []TplEnumChild{
					{ID: 0, Name: "pending"},
					{ID: 1, Name: "closed"},
				},
			},
		},
		Apis: []TplApi{
			{
				Name: "user.set_sim_info",
				Args: []TplApiArg{
					{Name: "info", Type: TplType{Name: "sim_info", Kind: TplKindStruct}},
				},
				Result: TplType{Name: "nil", Kind: TplKindNil},
				Note:   "设置sim信息\n无返回值",
			},
		},
	}

	goGen := NewGoGenerator(cfg)
	if err := goGen.Generate(schema); err != nil {
		t.Fatalf("generate go failed: %v", err)
	}
	tsGen := NewTsGenerator(cfg)
	if err := tsGen.Generate(schema); err != nil {
		t.Fatalf("generate ts failed: %v", err)
	}
	if err := generateDoc(schema, cfg); err != nil {
		t.Fatalf("generate doc failed: %v", err)
	}

	goEnum, err := os.ReadFile(filepath.Join(cfg.GoDir, "sb", "enum.go"))
	if err != nil {
		t.Fatalf("read go enum failed: %v", err)
	}
	goEnumText := string(goEnum)
	assertContains(t, goEnumText, "// OrderStatus 订单状态\n// Pending      待处理\n// Closed       已关闭")
	assertNotContains(t, goEnumText, "\nPending      待处理\n")

	goStruct, err := os.ReadFile(filepath.Join(cfg.GoDir, "sb", "struct_sim_info.go"))
	if err != nil {
		t.Fatalf("read go struct failed: %v", err)
	}
	goStructText := string(goStruct)
	assertContains(t, goStructText, "\t// 主键\n\t// 数据库ID\n\tId uint32")

	goRPC, err := os.ReadFile(filepath.Join(cfg.GoDir, "sb", "rpc.go"))
	if err != nil {
		t.Fatalf("read go rpc failed: %v", err)
	}
	goRPCText := string(goRPC)
	assertContains(t, goRPCText, "// CallUserSetSimInfo 设置sim信息\n// 无返回值\nfunc CallUserSetSimInfo")
	assertNotContains(t, goRPCText, "\n无返回值\nfunc CallUserSetSimInfo")

	tsEnum, err := os.ReadFile(filepath.Join(cfg.TsDir, "sb", "enum.ts"))
	if err != nil {
		t.Fatalf("read ts enum failed: %v", err)
	}
	tsEnumText := string(tsEnum)
	assertContains(t, tsEnumText, "// 订单状态\n// Pending      待处理\n// Closed       已关闭")
	assertContains(t, tsEnumText, "export enum OrderStatus {")
	assertNotContains(t, tsEnumText, "Settled = 6,\n\n}")
	assertNotContains(t, tsEnumText, "\nPending      待处理\nexport enum OrderStatus")

	tsRPC, err := os.ReadFile(filepath.Join(cfg.TsDir, "sb", "rpc.ts"))
	if err != nil {
		t.Fatalf("read ts rpc failed: %v", err)
	}
	tsRPCText := string(tsRPC)
	assertContains(t, tsRPCText, "    // 设置sim信息\n    // 无返回值\n    public userSetSimInfo = async")
	assertNotContains(t, tsRPCText, "        // 设置sim信息")
	assertContains(t, tsRPCText, "_.setAll(buf, (buf: _.Buffer) => _.setSimInfo(buf, info as any))")
	assertNotContains(t, tsRPCText, "/** 设置sim信息")

	doc, err := os.ReadFile(filepath.Join(cfg.GoDir, "sb", "DOC.md"))
	if err != nil {
		t.Fatalf("read doc failed: %v", err)
	}
	docText := string(doc)
	assertContains(t, docText, "| user_set_sim_info | info sim_info<br> | Void | 设置sim信息<br>无返回值 |")
	assertContains(t, docText, "> 订单状态\n> Pending      待处理\n> Closed       已关闭")
}

func TestTemplateGeneratorSkipsRewritingUnchangedFiles(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := Config{
		GoDir: filepath.Join(tmpDir, "go"),
		TsDir: filepath.Join(tmpDir, "ts"),
	}

	schema := &TplSchema{
		Structs: []TplStruct{
			{
				Name: "sim_info",
				Fields: []TplStructField{
					{Name: "id", Type: TplType{Name: "u32", Kind: TplKindBase}},
				},
			},
		},
		Enums: []TplEnum{
			{
				Name: "sim_operator",
				Children: []TplEnumChild{
					{ID: 1, Name: "cmcc"},
				},
			},
		},
		Apis: []TplApi{
			{
				Name: "sim.get",
				Args: []TplApiArg{
					{Name: "id", Type: TplType{Name: "u32", Kind: TplKindBase}},
				},
				Result: TplType{Name: "sim_info", Kind: TplKindStruct},
			},
		},
	}

	goGen := NewGoGenerator(cfg)
	if err := goGen.Generate(schema); err != nil {
		t.Fatalf("generate go failed: %v", err)
	}
	tsGen := NewTsGenerator(cfg)
	if err := tsGen.Generate(schema); err != nil {
		t.Fatalf("generate ts failed: %v", err)
	}
	if err := generateDoc(schema, cfg); err != nil {
		t.Fatalf("generate doc failed: %v", err)
	}

	paths := []string{
		filepath.Join(cfg.GoDir, "sb", "type.go"),
		filepath.Join(cfg.GoDir, "sb", "enum.go"),
		filepath.Join(cfg.GoDir, "sb", "struct_sim_info.go"),
		filepath.Join(cfg.GoDir, "sb", "api._.go"),
		filepath.Join(cfg.GoDir, "sb", "rpc.go"),
		filepath.Join(cfg.GoDir, "sb", "DOC.md"),
		filepath.Join(cfg.TsDir, "sb", "type.ts"),
		filepath.Join(cfg.TsDir, "sb", "enum.ts"),
		filepath.Join(cfg.TsDir, "sb", "struct_sim_info.ts"),
		filepath.Join(cfg.TsDir, "sb", "_.ts"),
		filepath.Join(cfg.TsDir, "sb", "rpc.ts"),
		filepath.Join(cfg.TsDir, "sb", "rpc_smoke.test.ts"),
		filepath.Join(cfg.TsDir, "sb", "DOC.md"),
	}

	stableTime := time.Date(2000, time.January, 2, 3, 4, 5, 0, time.UTC)
	wantTimes := make(map[string]time.Time, len(paths))
	for _, path := range paths {
		if err := os.Chtimes(path, stableTime, stableTime); err != nil {
			t.Fatalf("set file time for %s failed: %v", path, err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat %s failed: %v", path, err)
		}
		wantTimes[path] = info.ModTime()
	}

	if err := goGen.Generate(schema); err != nil {
		t.Fatalf("regenerate go failed: %v", err)
	}
	if err := tsGen.Generate(schema); err != nil {
		t.Fatalf("regenerate ts failed: %v", err)
	}
	if err := generateDoc(schema, cfg); err != nil {
		t.Fatalf("regenerate doc failed: %v", err)
	}

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat %s after regenerate failed: %v", path, err)
		}
		if !info.ModTime().Equal(wantTimes[path]) {
			t.Fatalf("expected %s mod time to stay %s, got %s", path, wantTimes[path], info.ModTime())
		}
	}
}

func assertContains(t *testing.T, text, want string) {
	t.Helper()
	if strings.Contains(text, want) {
		return
	}
	t.Fatalf("missing text %q", want)
}

func assertNotContains(t *testing.T, text, want string) {
	t.Helper()
	if !strings.Contains(text, want) {
		return
	}
	t.Fatalf("unexpected text %q", want)
}
