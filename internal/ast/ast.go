package ast

// TypeRef 表示 schema 中出现的类型引用。
type TypeRef struct {
	Name   string
	IsList bool
}

// Field 表示结构体字段。
// Embedded=true 且 Name 为空时表示嵌入字段。
type Field struct {
	Name     string
	Type     TypeRef
	Tag      string
	Note     string
	Embedded bool
}

// Struct 表示结构体定义。
type Struct struct {
	Name   string
	Fields []Field
	Note   string
}

// EnumMemberRaw 表示枚举成员，Value=nil 表示自动分配。
type EnumMemberRaw struct {
	Name  string
	Value *uint8
	Note  string
}

// Enum 表示枚举定义。
type Enum struct {
	Name    string
	Members []EnumMemberRaw
	Note    string
}

// APIArg 表示 API 参数。
type APIArg struct {
	Name string
	Type TypeRef
}

// API 表示 API 定义。
type API struct {
	Name   string
	Args   []APIArg
	Result TypeRef
	Note   string
}

// Schema 是 AST 根节点。
type Schema struct {
	Structs []Struct
	Enums   []Enum
	APIs    []API
	Note    string
}
