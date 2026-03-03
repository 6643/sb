package semantic

import (
	"testing"

	"sb/internal/ast"
)

func TestResolveBasic(t *testing.T) {
	schema := &ast.Schema{
		Enums: []ast.Enum{{
			Name: "Status",
			Members: []ast.EnumMemberRaw{
				{Name: "Ok"},
				{Name: "Err"},
			},
		}},
		Structs: []ast.Struct{{
			Name:   "User",
			Fields: []ast.Field{{Name: "id", Type: ast.TypeRef{Name: "u32"}}},
		}},
		APIs: []ast.API{{
			Name:   "user.get",
			Args:   []ast.APIArg{{Name: "id", Type: ast.TypeRef{Name: "u32"}}},
			Result: ast.TypeRef{Name: "User"},
		}},
	}

	irSchema, err := Resolve(schema)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if len(irSchema.Enums) != 1 || irSchema.Enums[0].Members[1].ID != 1 {
		t.Fatalf("unexpected enum resolve result: %+v", irSchema.Enums)
	}
}

func TestResolveUndefinedType(t *testing.T) {
	schema := &ast.Schema{
		Structs: []ast.Struct{{
			Name:   "User",
			Fields: []ast.Field{{Name: "info", Type: ast.TypeRef{Name: "NotExist"}}},
		}},
	}
	_, err := Resolve(schema)
	if err == nil {
		t.Fatal("expected undefined type error, got nil")
	}
}
