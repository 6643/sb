package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoGeneratorWritesSchemaFile(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGoGenerator(Config{GoDir: tmpDir, GoTag: "json"})

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
					{Name: "id", Type: TplType{Name: "u32", Kind: TplKindBase}, Tag: "_id"},
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
		t.Fatalf("generate go failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "sb", "schema.gen.go"))
	if err != nil {
		t.Fatalf("read schema.gen.go failed: %v", err)
	}
	text := string(content)

	assertContains(t, text, "package sb")
	assertContains(t, text, "rt \"sb/go/sb/rt\"")
	assertContains(t, text, "func isSimOperator(v SimOperator) bool {")
	assertContains(t, text, "func defaultSimOperator() SimOperator { return SimOperatorZz }")
	assertContains(t, text, "func normalizeSimOperator(v SimOperator) SimOperator {")
	assertContains(t, text, "func isAssignableSimOperator(v SimOperator) bool {")
	assertContains(t, text, "func SetSim(buf *bytes.Buffer, s *Sim) error {")
	assertContains(t, text, "func GetSim(buf *bytes.Buffer, s *Sim) error {")
	assertContains(t, text, "var simHeaderWidths = [...]uint8{1, 2, 1, 1, 2, 2, 2}")
	assertContains(t, text, "if err := rt.ReadHeader(header, simHeaderWidths[:], headerStates[:], \"Sim header\"); err != nil { return err }")
	assertContains(t, text, "if err := rt.WriteHeader(headerData[:], simHeaderWidths[:], headerStates[:]); err != nil { return fmt.Errorf(\"SetSim write header: %w\", err) }")
	assertContains(t, text, "func isZeroSim(s *Sim) bool {")
	assertContains(t, text, "func getSimInfoListBody(buf *bytes.Buffer, state uint8) ([]*SimInfo, error) {")
	assertContains(t, text, "func getSimInfoListBodyReuse(buf *bytes.Buffer, state uint8, dst []*SimInfo) ([]*SimInfo, error) {")
	assertContains(t, text, "return rt.GetDefaultPtrListInto(buf, state, dst, defaultSimInfo, GetSimInfo)")
	assertContains(t, text, "return rt.SetDefaultPtrList(")
	assertContains(t, text, "\"SimInfo\",")
	assertContains(t, text, "validateSimInfo,")
	assertContains(t, text, "return rt.SizeDefaultPtrList(")
	assertContains(t, text, "fieldSizeTitles, err := rt.SizeTextList(s.Titles)")
	assertContains(t, text, "fieldSizeTags, err := rt.SizeZeroList(s.Tags, 4)")
	assertContains(t, text, "valueTags, err := rt.GetZeroListInto(buf, tagsState, reuseTags, rt.GetU32)")
	assertContains(t, text, "if err := rt.SetZeroList(buf, tagsState, s.Tags, 4, rt.SetU32); err != nil { return fmt.Errorf(\"SetSim Tags: %w\", err) }")
	assertNotContains(t, text, "func IsSimOperator(v SimOperator) bool {")
	assertNotContains(t, text, "func DefaultSimOperator()")
	assertNotContains(t, text, "func ValidateSimInfo(")
	assertContains(t, text, "`json:\"_id\"`")

	apiContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "api._.go"))
	if err != nil {
		t.Fatalf("read api._.go failed: %v", err)
	}
	apiText := string(apiContent)
	assertContains(t, apiText, "func UserGetCountHandler(w http.ResponseWriter, r *http.Request) {")
	assertContains(t, apiText, "var page uint8")
	assertContains(t, apiText, "value, err := rt.GetU8(buf)")
	assertContains(t, apiText, "page = value")
	assertContains(t, apiText, "if err := rt.SetU8(&body, result); err != nil { w.WriteHeader(http.StatusInternalServerError); return }")
	assertNotContains(t, apiText, "type ApiUserGetCountReq struct {")
	assertNotContains(t, apiText, "type ApiUserGetCountResp struct {")
	assertNotContains(t, apiText, "readApiUserGetCountReq")
	assertNotContains(t, apiText, "setApiUserGetCountResp")
	assertContains(t, apiText, "var _ = slices.Equal[[]int, int]")

	rpcContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc.go"))
	if err != nil {
		t.Fatalf("read rpc.go failed: %v", err)
	}
	rpcText := string(rpcContent)
	assertContains(t, rpcText, "func CallUserGetCount(c *Client, ctx context.Context, page uint8) (result uint8, errCode RpcErrCode) {")
	assertContains(t, rpcText, "if err := rt.SetU8(&buf, page); err != nil {")
	assertContains(t, rpcText, "value, err := rt.GetU8(respBuf)")
	assertNotContains(t, rpcText, "ApiUserGetCountReq")
	assertNotContains(t, rpcText, "readApiUserGetCountResp")

	stubContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "api.user.get_count.go"))
	if err != nil {
		t.Fatalf("read api.user.get_count.go failed: %v", err)
	}
	stubText := string(stubContent)
	assertContains(t, stubText, "func user_get_count(ctx context.Context, page uint8) (result uint8, errCode RpcErrCode) {")
}
