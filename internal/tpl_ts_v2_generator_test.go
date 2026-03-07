package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTsV2GeneratorWritesSchemaFile(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewTsV2Generator(Config{TsDir: tmpDir})

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
		t.Fatalf("generate ts v2 failed: %v", err)
	}

	typeContent, err := os.ReadFile(filepath.Join(tmpDir, "sbv2", "type.ts"))
	if err != nil {
		t.Fatalf("read type.ts failed: %v", err)
	}
	typeText := string(typeContent)
	assertContains(t, typeText, "export class BitWriter")
	assertContains(t, typeText, "export const getBitmapListCompact")
	assertContains(t, typeText, "export const textState")

	enumContent, err := os.ReadFile(filepath.Join(tmpDir, "sbv2", "enum.ts"))
	if err != nil {
		t.Fatalf("read enum.ts failed: %v", err)
	}
	enumText := string(enumContent)
	assertContains(t, enumText, "export const DefaultSimOperator")
	assertContains(t, enumText, "export const NormalizeSimOperator")
	assertContains(t, enumText, "export const getSimOperatorListBody")

	structContent, err := os.ReadFile(filepath.Join(tmpDir, "sbv2", "struct_sim.ts"))
	if err != nil {
		t.Fatalf("read struct_sim.ts failed: %v", err)
	}
	structText := string(structContent)
	assertContains(t, structText, "const header = new rt.BitWriter(")
	assertContains(t, structText, "rt.getBitmapListCompact<number>(buf, tagsState")
	assertContains(t, structText, "_.getSimInfoListBody(buf, infosState)")
	assertContains(t, structText, "rt.setTextListCompact(buf, titlesState, s.titles)")

	if _, err := os.Stat(filepath.Join(tmpDir, "sbv2", "struct_api_user_get_count_req.ts")); !os.IsNotExist(err) {
		t.Fatalf("expected legacy api wrapper file removed, got err=%v", err)
	}

	rpcContent, err := os.ReadFile(filepath.Join(tmpDir, "sbv2", "rpc.ts"))
	if err != nil {
		t.Fatalf("read rpc.ts failed: %v", err)
	}
	rpcText := string(rpcContent)
	assertContains(t, rpcText, "_.setU8(buf, page")
	assertContains(t, rpcText, "const [value, err] = _.getU8(respBuf);")
	assertNotContains(t, rpcText, "newApiUserGetCountReq")
	assertNotContains(t, rpcText, "readApiUserGetCountResp")

	smokeContent, err := os.ReadFile(filepath.Join(tmpDir, "sbv2", "rpc_smoke.test.ts"))
	if err != nil {
		t.Fatalf("read rpc_smoke.test.ts failed: %v", err)
	}
	assertContains(t, string(smokeContent), "method userGetCount exists")
}
