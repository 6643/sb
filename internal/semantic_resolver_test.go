package internal

import (
	"strings"
	"testing"
)

func TestResolveBasic(t *testing.T) {
	schema := &Schema{
		Enums: []Enum{{
			Name: "Status",
			Members: []EnumMemberRaw{
				{Name: "Ok"},
				{Name: "Err"},
			},
		}},
		Structs: []Struct{{
			Name:   "User",
			Fields: []Field{{Name: "id", Type: TypeRef{Name: "u32"}}},
		}},
		APIs: []API{{
			Name:   "user.get",
			Args:   []APIArg{{Name: "id", Type: TypeRef{Name: "u32"}}},
			Result: TypeRef{Name: "User"},
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
	schema := &Schema{
		Structs: []Struct{{
			Name:   "User",
			Fields: []Field{{Name: "info", Type: TypeRef{Name: "NotExist"}}},
		}},
	}
	_, err := Resolve(schema)
	if err == nil {
		t.Fatal("expected undefined type error, got nil")
	}
}

func TestResolveEnumRejectDuplicateIDsInSameEnum(t *testing.T) {
	zero := uint8(0)
	schema := &Schema{
		Enums: []Enum{{
			Name: "Status",
			Members: []EnumMemberRaw{
				{Name: "Unknown", Value: &zero},
				{Name: "Active", Value: &zero},
			},
		}},
	}

	_, err := Resolve(schema)
	if err == nil {
		t.Fatal("expected duplicate enum id error, got nil")
	}
	if !strings.Contains(err.Error(), "成员 ID 重复") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveEnumAllowSameIDAcrossDifferentEnums(t *testing.T) {
	zero := uint8(0)
	schema := &Schema{
		Enums: []Enum{
			{
				Name: "Status",
				Members: []EnumMemberRaw{
					{Name: "Unknown", Value: &zero},
				},
			},
			{
				Name: "Mode",
				Members: []EnumMemberRaw{
					{Name: "Off", Value: &zero},
				},
			},
		},
	}

	irSchema, err := Resolve(schema)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if len(irSchema.Enums) != 2 {
		t.Fatalf("unexpected enum count: %d", len(irSchema.Enums))
	}
	if irSchema.Enums[0].Members[0].ID != 0 || irSchema.Enums[1].Members[0].ID != 0 {
		t.Fatalf("unexpected enum ids: %+v", irSchema.Enums)
	}
}
