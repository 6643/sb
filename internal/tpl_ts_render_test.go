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
	assertContains(t, typeText, "export interface StructField<T>")
	assertContains(t, typeText, "export interface StructMeta<T>")
	assertContains(t, typeText, "export const defineStruct = <T>(meta: StructMeta<T>): StructMeta<T> => meta;")
	assertContains(t, typeText, "export const getStruct = <T>(meta: StructMeta<T>, buf: Buffer): [T, Err] => {")
	assertContains(t, typeText, "export const setStruct = <T>(meta: StructMeta<T>, buf: Buffer, value: T | null | undefined): Err => {")
	assertContains(t, typeText, "export const validateStruct = <T>(meta: StructMeta<T>, value: T | null | undefined): Err => {")
	assertContains(t, typeText, "export const eqStruct = <T>(meta: StructMeta<T>, a: T | null | undefined, b: T | null | undefined): boolean => {")

	enumContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "enum.ts"))
	if err != nil {
		t.Fatalf("read enum.ts failed: %v", err)
	}
	enumText := string(enumContent)
	assertContains(t, enumText, "const simOperatorMeta = rt.defineEnum<SimOperator>(")
	assertContains(t, enumText, "export const DefaultSimOperator")
	assertContains(t, enumText, "export const NormalizeSimOperator = (v: SimOperator): SimOperator => rt.normalizeEnum(simOperatorMeta, v);")
	assertContains(t, enumText, "export const IsAssignableSimOperator = (v: SimOperator): boolean => rt.isAssignableEnum(simOperatorMeta, v);")
	assertContains(t, enumText, "export const getSimOperatorListBody = (buf: rt.Buffer, state: number): [SimOperator[], rt.Err] => rt.getEnumList(simOperatorMeta, buf, state);")
	assertContains(t, enumText, "export const setSimOperatorListBody = (buf: rt.Buffer, state: number, v: SimOperator[]): rt.Err => rt.setEnumList(simOperatorMeta, buf, state, v);")
	assertNotContains(t, enumText, "switch (v) {")

	structContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "struct_sim.ts"))
	if err != nil {
		t.Fatalf("read struct_sim.ts failed: %v", err)
	}
	structText := string(structContent)
	assertContains(t, structText, "const simHeaderWidths = [1, 2, 1, 1, 2, 2, 2] as const;")
	assertContains(t, structText, "import {")
	assertNotContains(t, structText, "import * as _ from \"./_\"")
	assertContains(t, structText, "const simMeta = rt.defineStruct<Sim>({")
	assertContains(t, structText, "export const newSim = (): Sim => rt.newStruct(simMeta, getSim, setSim);")
	assertContains(t, structText, "export const isZeroSim = (s: Sim | null | undefined): boolean => rt.isZeroStruct(simMeta, s);")
	assertContains(t, structText, "export const validateSim = (s: Sim | null | undefined): rt.Err => rt.validateStruct(simMeta, s);")
	assertContains(t, structText, "export const getSim = (buf: rt.Buffer): [Sim, rt.Err] => rt.getStruct(simMeta, buf);")
	assertContains(t, structText, "export const setSim = (buf: rt.Buffer, s: Sim): rt.Err => rt.setStruct(simMeta, buf, s);")
	assertContains(t, structText, "export const eqSim = (a: Sim | null | undefined, b: Sim | null | undefined): boolean => rt.eqStruct(simMeta, a, b);")
	assertNotContains(t, structText, "new rt.BitWriter(")
	assertNotContains(t, structText, "new rt.BitReader(")
	assertContains(t, structText, "rt.zeroListField<Sim, number>(")
	assertContains(t, structText, "rt.textListField<Sim>(\"titles\", \"Titles\")")
	assertContains(t, structText, "rt.defaultListField<Sim, SimInfo>(")
	assertNotContains(t, structText, "const [headerStates, errHeaderState] = rt.readHeader(")
	assertNotContains(t, structText, "const [header, errHeader] = rt.writeHeader(")
	assertNotContains(t, structText, "const headerStates = [")

	structInfoContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "struct_sim_info.ts"))
	if err != nil {
		t.Fatalf("read struct_sim_info.ts failed: %v", err)
	}
	structInfoText := string(structInfoContent)
	assertContains(t, structInfoText, "const simInfoMeta = rt.defineStruct<SimInfo>({")
	assertContains(t, structInfoText, "export const getSimInfo = (buf: rt.Buffer): [SimInfo, rt.Err] => rt.getStruct(simInfoMeta, buf);")
	assertContains(t, structInfoText, "export const setSimInfo = (buf: rt.Buffer, s: SimInfo): rt.Err => rt.setStruct(simInfoMeta, buf, s);")
	assertNotContains(t, structInfoText, "import * as _ from \"./_\"")
	assertNotContains(t, structInfoText, "const [headerStates, errHeaderState] = rt.readHeader(")
	assertNotContains(t, structInfoText, "const headerStates = [")

	if _, err := os.Stat(filepath.Join(tmpDir, "sb", "struct_api_user_get_count_req.ts")); !os.IsNotExist(err) {
		t.Fatalf("expected legacy api wrapper file removed, got err=%v", err)
	}

	rpcContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc.ts"))
	if err != nil {
		t.Fatalf("read rpc.ts failed: %v", err)
	}
	rpcText := string(rpcContent)
	assertContains(t, rpcText, "import * as rt from \"./type\"")
	assertNotContains(t, rpcText, "import * as _ from \"./_\"")
	assertContains(t, rpcText, "export type RpcStatus = RpcErrCode | number;")
	assertContains(t, rpcText, "async function readResponseBytes(res: Response, maxRespBytes: number): Promise<[Uint8Array | null, RpcErrCode | null]> {")
	assertContains(t, rpcText, "const reader = res.body.getReader();")
	assertContains(t, rpcText, "await reader.cancel();")
	assertContains(t, rpcText, "const buf = new rt.Buffer();")
	assertContains(t, rpcText, "rt.setU8(buf, page")
	assertContains(t, rpcText, "const [value, err] = rt.getU8(respBuf);")
	assertContains(t, rpcText, "return [null, res.status];")
	assertNotContains(t, rpcText, "newApiUserGetCountReq")
	assertNotContains(t, rpcText, "readApiUserGetCountResp")

	smokeContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc_smoke.test.ts"))
	if err != nil {
		t.Fatalf("read rpc_smoke.test.ts failed: %v", err)
	}
	assertContains(t, string(smokeContent), "process.env.SB_BASE_URL")
	assertNotContains(t, string(smokeContent), "process.env.SIMPLE_BIN_BASE_URL")
	assertContains(t, string(smokeContent), "method userGetCount exists")
}
