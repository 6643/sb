package ast

type TypeKind int

const (
	KindBase TypeKind = iota
	KindStruct
	KindEnum
)

type Type struct {
	Name   string
	Kind   TypeKind
	IsList bool
}

type StructField struct {
	Name string
	Type Type
	Tag  string
	Note string
}

type Struct struct {
	Name   string
	Fields []StructField
	Note   string
}

type EnumChild struct {
	ID   uint8
	Name string
	Note string
}

type Enum struct {
	Name     string
	Children []EnumChild
	Note     string
}

type ApiArg struct {
	Name string
	Type Type
}

type Api struct {
	Name   string
	Args   []ApiArg
	Result Type // Success return type, nil means void
	Note   string
}

type Schema struct {
	Structs []Struct
	Enums   []Enum
	Apis    []Api
	Note    string
}
