package internal

// IRTypeKind 表示类型类别。
type IRTypeKind int

const (
	IRKindBase IRTypeKind = iota
	IRKindStruct
	IRKindEnum
	IRKindNil
)

// IRType 是语义解析后的类型。
type IRType struct {
	Name   string
	Kind   IRTypeKind
	IsList bool
}

// IRField 是扁平化后的结构体字段。
type IRField struct {
	Name string
	Type IRType
	Tag  string
	Note string
}

// IRStruct 表示结构体。
type IRStruct struct {
	Name   string
	Fields []IRField
	Note   string
}

// IREnumMember 表示最终枚举成员。
type IREnumMember struct {
	ID   uint8
	Name string
	Note string
}

// IREnum 表示枚举。
type IREnum struct {
	Name    string
	Members []IREnumMember
	Note    string
}

// IRAPIArg 表示 API 参数。
type IRAPIArg struct {
	Name string
	Type IRType
}

// IRAPI 表示 API。
type IRAPI struct {
	Name   string
	Args   []IRAPIArg
	Result IRType
	Note   string
}

// IRSchema 是后端输入。
type IRSchema struct {
	Structs []IRStruct
	Enums   []IREnum
	APIs    []IRAPI
	Note    string
}
