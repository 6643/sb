package ir

// TypeKind 表示类型类别。
type TypeKind int

const (
	KindBase TypeKind = iota
	KindStruct
	KindEnum
	KindNil
)

// Type 是语义解析后的类型。
type Type struct {
	Name   string
	Kind   TypeKind
	IsList bool
}

// Field 是扁平化后的结构体字段。
type Field struct {
	Name string
	Type Type
	Tag  string
	Note string
}

// Struct 表示结构体。
type Struct struct {
	Name   string
	Fields []Field
	Note   string
}

// EnumMember 表示最终枚举成员。
type EnumMember struct {
	ID   uint8
	Name string
	Note string
}

// Enum 表示枚举。
type Enum struct {
	Name    string
	Members []EnumMember
	Note    string
}

// APIArg 表示 API 参数。
type APIArg struct {
	Name string
	Type Type
}

// API 表示 API。
type API struct {
	Name   string
	Args   []APIArg
	Result Type
	Note   string
}

// Schema 是后端输入。
type Schema struct {
	Structs []Struct
	Enums   []Enum
	APIs    []API
	Note    string
}
