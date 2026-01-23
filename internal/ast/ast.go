package ast

// TypeKind 类型分类: 基础类型, 结构体, 枚举
type TypeKind int

const (
	KindBase TypeKind = iota // 基础类型 (如 u8, text)
	KindStruct               // 用户定义的结构体
	KindEnum                 // 用户定义的枚举
)

// Type 抽象类型定义
// 涵盖了基础类型, 引用类型以及数组/列表形式
type Type struct {
	Name   string
	Kind   TypeKind
	IsList bool // 是否为数组/切片 ([T])
}

// StructField 结构体字段定义
type StructField struct {
	Name string
	Type Type
	Tag  string // Go struct tag (如 `json:"id"`)
	Note string // 字段注释
}

// Struct 结构体定义
type Struct struct {
	Name   string
	Fields []StructField
	Note   string
}

// EnumChild 枚举成员定义
type EnumChild struct {
	ID   uint8  // 枚举数值 (0-255)
	Name string
	Note string
}

// Enum 枚举定义 (支持 u8 范围内的数值映射)
type Enum struct {
	Name     string
	Children []EnumChild
	Note     string
}

// ApiArg API 参数定义
type ApiArg struct {
	Name string
	Type Type
}

// Api 远程调用接口定义
type Api struct {
	Name   string
	Args   []ApiArg
	Result Type // 返回类型 (nil 表示 void/无返回值)
	Note   string
}

// Schema 完整的协议描述文件 (AST 根节点)
type Schema struct {
	Structs []Struct
	Enums   []Enum
	Apis    []Api
	Note    string
}
