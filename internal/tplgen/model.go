package tplgen

import "sb/internal/ir"

// TypeKind 类型分类。
type TypeKind int

const (
	KindBase TypeKind = iota
	KindStruct
	KindEnum
	KindNil
)

// Type 模板层抽象类型。
type Type struct {
	Name   string
	Kind   TypeKind
	IsList bool
}

// StructField 结构体字段。
type StructField struct {
	Name string
	Type Type
	Tag  string
	Note string
}

// Struct 结构体定义。
type Struct struct {
	Name   string
	Fields []StructField
	Note   string
}

// EnumChild 枚举成员。
type EnumChild struct {
	ID   uint8
	Name string
	Note string
}

// Enum 枚举定义。
type Enum struct {
	Name     string
	Children []EnumChild
	Note     string
}

// ApiArg API 参数。
type ApiArg struct {
	Name string
	Type Type
}

// Api API 定义。
type Api struct {
	Name   string
	Args   []ApiArg
	Result Type
	Note   string
}

// Schema 模板后端输入。
type Schema struct {
	Structs []Struct
	Enums   []Enum
	Apis    []Api
	Note    string
}

func fromIRSchema(src *ir.Schema) *Schema {
	if src == nil {
		return &Schema{}
	}

	dst := &Schema{Note: src.Note}
	dst.Structs = convertStructs(src.Structs)
	dst.Enums = convertEnums(src.Enums)
	dst.Apis = convertAPIs(src.APIs)
	return dst
}

func convertStructs(items []ir.Struct) []Struct {
	if len(items) == 0 {
		return nil
	}

	out := make([]Struct, 0, len(items))
	for _, it := range items {
		out = append(out, Struct{
			Name:   it.Name,
			Fields: convertFields(it.Fields),
			Note:   it.Note,
		})
	}
	return out
}

func convertFields(items []ir.Field) []StructField {
	if len(items) == 0 {
		return nil
	}

	out := make([]StructField, 0, len(items))
	for _, it := range items {
		out = append(out, StructField{
			Name: it.Name,
			Type: convertType(it.Type),
			Tag:  it.Tag,
			Note: it.Note,
		})
	}
	return out
}

func convertEnums(items []ir.Enum) []Enum {
	if len(items) == 0 {
		return nil
	}

	out := make([]Enum, 0, len(items))
	for _, it := range items {
		out = append(out, Enum{
			Name:     it.Name,
			Children: convertEnumMembers(it.Members),
			Note:     it.Note,
		})
	}
	return out
}

func convertEnumMembers(items []ir.EnumMember) []EnumChild {
	if len(items) == 0 {
		return nil
	}

	out := make([]EnumChild, 0, len(items))
	for _, it := range items {
		out = append(out, EnumChild{ID: it.ID, Name: it.Name, Note: it.Note})
	}
	return out
}

func convertAPIs(items []ir.API) []Api {
	if len(items) == 0 {
		return nil
	}

	out := make([]Api, 0, len(items))
	for _, it := range items {
		out = append(out, Api{
			Name:   it.Name,
			Args:   convertAPIArgs(it.Args),
			Result: convertType(it.Result),
			Note:   it.Note,
		})
	}
	return out
}

func convertAPIArgs(items []ir.APIArg) []ApiArg {
	if len(items) == 0 {
		return nil
	}

	out := make([]ApiArg, 0, len(items))
	for _, it := range items {
		out = append(out, ApiArg{Name: it.Name, Type: convertType(it.Type)})
	}
	return out
}

func convertType(t ir.Type) Type {
	kind := KindBase
	switch t.Kind {
	case ir.KindStruct:
		kind = KindStruct
	case ir.KindEnum:
		kind = KindEnum
	case ir.KindNil:
		kind = KindNil
	}

	if kind == KindNil {
		return Type{Name: "nil", Kind: KindNil, IsList: false}
	}
	return Type{Name: t.Name, Kind: kind, IsList: t.IsList}
}
