package internal

// TplTypeKind 类型分类。
type TplTypeKind int

const (
	TplKindBase TplTypeKind = iota
	TplKindStruct
	TplKindEnum
	TplKindNil
)

// TplType 模板层抽象类型。
type TplType struct {
	Name   string
	Kind   TplTypeKind
	IsList bool
}

// TplStructField 结构体字段。
type TplStructField struct {
	Name string
	Type TplType
	Tag  string
	Note string
}

// TplStruct 结构体定义。
type TplStruct struct {
	Name   string
	Fields []TplStructField
	Note   string
}

// TplEnumChild 枚举成员。
type TplEnumChild struct {
	ID   uint8
	Name string
	Note string
}

// TplEnum 枚举定义。
type TplEnum struct {
	Name     string
	Children []TplEnumChild
	Note     string
}

// TplApiArg API 参数。
type TplApiArg struct {
	Name string
	Type TplType
}

// TplApi API 定义。
type TplApi struct {
	Name   string
	Args   []TplApiArg
	Result TplType
	Note   string
}

// TplSchema 模板后端输入。
type TplSchema struct {
	Structs []TplStruct
	Enums   []TplEnum
	Apis    []TplApi
	Note    string
}

func fromIRSchema(src *IRSchema) *TplSchema {
	if src == nil {
		return &TplSchema{}
	}

	dst := &TplSchema{Note: src.Note}
	dst.Structs = convertStructs(src.Structs)
	dst.Enums = convertEnums(src.Enums)
	dst.Apis = convertAPIs(src.APIs)
	return dst
}

func convertStructs(items []IRStruct) []TplStruct {
	if len(items) == 0 {
		return nil
	}

	out := make([]TplStruct, 0, len(items))
	for _, it := range items {
		out = append(out, TplStruct{
			Name:   it.Name,
			Fields: convertFields(it.Fields),
			Note:   it.Note,
		})
	}
	return out
}

func convertFields(items []IRField) []TplStructField {
	if len(items) == 0 {
		return nil
	}

	out := make([]TplStructField, 0, len(items))
	for _, it := range items {
		out = append(out, TplStructField{
			Name: it.Name,
			Type: convertType(it.Type),
			Tag:  it.Tag,
			Note: it.Note,
		})
	}
	return out
}

func convertEnums(items []IREnum) []TplEnum {
	if len(items) == 0 {
		return nil
	}

	out := make([]TplEnum, 0, len(items))
	for _, it := range items {
		out = append(out, TplEnum{
			Name:     it.Name,
			Children: convertEnumMembers(it.Members),
			Note:     it.Note,
		})
	}
	return out
}

func convertEnumMembers(items []IREnumMember) []TplEnumChild {
	if len(items) == 0 {
		return nil
	}

	out := make([]TplEnumChild, 0, len(items))
	for _, it := range items {
		out = append(out, TplEnumChild{ID: it.ID, Name: it.Name, Note: it.Note})
	}
	return out
}

func convertAPIs(items []IRAPI) []TplApi {
	if len(items) == 0 {
		return nil
	}

	out := make([]TplApi, 0, len(items))
	for _, it := range items {
		out = append(out, TplApi{
			Name:   it.Name,
			Args:   convertAPIArgs(it.Args),
			Result: convertType(it.Result),
			Note:   it.Note,
		})
	}
	return out
}

func convertAPIArgs(items []IRAPIArg) []TplApiArg {
	if len(items) == 0 {
		return nil
	}

	out := make([]TplApiArg, 0, len(items))
	for _, it := range items {
		out = append(out, TplApiArg{Name: it.Name, Type: convertType(it.Type)})
	}
	return out
}

func convertType(t IRType) TplType {
	kind := TplKindBase
	switch t.Kind {
	case IRKindStruct:
		kind = TplKindStruct
	case IRKindEnum:
		kind = TplKindEnum
	case IRKindNil:
		kind = TplKindNil
	}

	if kind == TplKindNil {
		return TplType{Name: "nil", Kind: TplKindNil, IsList: false}
	}
	return TplType{Name: t.Name, Kind: kind, IsList: t.IsList}
}
