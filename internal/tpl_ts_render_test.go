package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTsGeneratorWritesSchemaFile(t *testing.T) {
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
				Name: "sim_info",
				Fields: []TplStructField{
					{Name: "title", Type: TplType{Name: "text", Kind: TplKindBase}},
				},
			},
			{
				Name: "sim",
				Fields: []TplStructField{
					{Name: "id", Type: TplType{Name: "u32", Kind: TplKindBase}},
					{Name: "name", Type: TplType{Name: "text", Kind: TplKindBase}},
					{Name: "operator", Type: TplType{Name: "sim_operator", Kind: TplKindEnum}},
					{Name: "ok", Type: TplType{Name: "bool", Kind: TplKindBase}},
					{Name: "tags", Type: TplType{Name: "u32", Kind: TplKindBase, IsList: true}},
					{Name: "titles", Type: TplType{Name: "text", Kind: TplKindBase, IsList: true}},
					{Name: "infos", Type: TplType{Name: "sim_info", Kind: TplKindStruct, IsList: true}},
				},
			},
		},
		Apis: []TplApi{
			{
				Name: "user.get_count",
				Args: []TplApiArg{
					{Name: "page", Type: TplType{Name: "u8", Kind: TplKindBase}},
				},
				Result: TplType{Name: "u8", Kind: TplKindBase},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate ts failed: %v", err)
	}

	typeContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "type.ts"))
	if err != nil {
		t.Fatalf("read type.ts failed: %v", err)
	}
	typeText := string(typeContent)
	assertContains(t, typeText, "export interface BitWriter")
	assertContains(t, typeText, "export const BitWriter = _BitWriterCtor as unknown as BitWriterCtor;")
	assertContains(t, typeText, "export interface BitReader")
	assertContains(t, typeText, "export const BitReader = _BitReaderCtor as unknown as BitReaderCtor;")
	assertContains(t, typeText, "export interface Buffer")
	assertContains(t, typeText, "export const Buffer = _BufferCtor as unknown as BufferCtor;")
	assertContains(t, typeText, "export const readHeader = (data: Uint8Array, widths: readonly number[], kind: string): [number[], Err] => {")
	assertContains(t, typeText, "export const writeHeader = (widths: readonly number[], values: readonly number[]): [Uint8Array, Err] => {")
	assertContains(t, typeText, "export const getDefaultList")
	assertContains(t, typeText, "export const getZeroList = <T>(")
	assertContains(t, typeText, "export const setZeroList = <T>(")
	assertContains(t, typeText, "export const textState")

	enumContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "enum.ts"))
	if err != nil {
		t.Fatalf("read enum.ts failed: %v", err)
	}
	enumText := string(enumContent)
	assertContains(t, enumText, "export const DefaultSimOperator")
	assertContains(t, enumText, "export const NormalizeSimOperator")
	assertContains(t, enumText, "export const getSimOperatorListBody")

	structContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "struct_sim.ts"))
	if err != nil {
		t.Fatalf("read struct_sim.ts failed: %v", err)
	}
	structText := string(structContent)
	assertContains(t, structText, "const simHeaderWidths = [1, 2, 1, 1, 2, 2, 2] as const;")
	assertContains(t, structText, "const [headerStates, errHeaderState] = rt.readHeader(header, simHeaderWidths, \"Sim header\");")
	assertContains(t, structText, "const [header, errHeader] = rt.writeHeader(simHeaderWidths, headerStates);")
	assertNotContains(t, structText, "new rt.BitWriter(")
	assertNotContains(t, structText, "new rt.BitReader(")
	assertContains(t, structText, "rt.getZeroList<number>(buf, tagsState, 0, rt.getU32)")
	assertContains(t, structText, "_.getSimInfoListBody(buf, infosState)")
	assertContains(t, structText, "rt.setTextList(buf, titlesState, s.titles)")
	assertContains(t, structText, "rt.setZeroList<number>(buf, tagsState, s.tags, 0, rt.setU32)")

	if _, err := os.Stat(filepath.Join(tmpDir, "sb", "struct_api_user_get_count_req.ts")); !os.IsNotExist(err) {
		t.Fatalf("expected legacy api wrapper file removed, got err=%v", err)
	}

	rpcContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc.ts"))
	if err != nil {
		t.Fatalf("read rpc.ts failed: %v", err)
	}
	rpcText := string(rpcContent)
	assertContains(t, rpcText, "export type RpcStatus = RpcErrCode | number;")
	assertContains(t, rpcText, "async function readResponseBytes(res: Response, maxRespBytes: number): Promise<[Uint8Array | null, RpcErrCode | null]> {")
	assertContains(t, rpcText, "const reader = res.body.getReader();")
	assertContains(t, rpcText, "await reader.cancel();")
	assertContains(t, rpcText, "_.setU8(buf, page")
	assertContains(t, rpcText, "const [value, err] = _.getU8(respBuf);")
	assertContains(t, rpcText, "return [null, res.status];")
	assertNotContains(t, rpcText, "newApiUserGetCountReq")
	assertNotContains(t, rpcText, "readApiUserGetCountResp")

	smokeContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc_smoke.test.ts"))
	if err != nil {
		t.Fatalf("read rpc_smoke.test.ts failed: %v", err)
	}
	assertContains(t, string(smokeContent), "method userGetCount exists")
}
