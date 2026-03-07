package internal

type TsGenerator struct {
	Config Config
}

func NewTsGenerator(cfg Config) *TsGenerator {
	return &TsGenerator{Config: cfg}
}

func (g *TsGenerator) getTsType(t TplType) string {
	switch t.Name {
	case "i8", "u8", "i16", "u16", "i32", "u32", "f32", "f64":
		return "number"
	case "i64", "u64":
		return "bigint"
	case "bool":
		return "boolean"
	case "text":
		return "string"
	case "bin":
		return "Uint8Array"
	}
	return PascalCase(t.Name)
}

func (g *TsGenerator) getTsLogicType(t TplType) string {
	if t.Name == "nil" {
		return "void"
	}
	base := g.getTsType(TplType{Name: t.Name, Kind: t.Kind})
	if !t.IsList {
		return base
	}
	return base + "[]"
}

func (g *TsGenerator) getTsValue(name string) string {
	switch name {
	case "text":
		return `""`
	case "bin":
		return "new Uint8Array(0)"
	case "bool":
		return "false"
	case "i64", "u64":
		return "0n"
	default:
		return "0"
	}
}
