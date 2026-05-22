package internal

var baseTypes = map[string]struct{}{
	"i8": {}, "u8": {}, "i16": {}, "u16": {},
	"i32": {}, "u32": {}, "i64": {}, "u64": {},
	"f32": {}, "f64": {}, "bool": {}, "text": {}, "bin": {},
}

// Resolve 把 AST 解析为可直接代码生成的 IR。
func Resolve(schema *Schema) (*IRSchema, error) {
	resolver := &resolver{
		structDefs: make(map[string]Struct),
		enumDefs:   make(map[string]Enum),
		fieldCache: make(map[string][]IRField),
	}
	return resolver.resolve(schema)
}

type resolver struct {
	structDefs map[string]Struct
	enumDefs   map[string]Enum
	fieldCache map[string][]IRField
}

func (r *resolver) resolve(schema *Schema) (*IRSchema, error) {
	if err := r.collectDefs(schema); err != nil {
		return nil, err
	}

	enums, err := r.resolveEnums(schema.Enums)
	if err != nil {
		return nil, err
	}

	structs, err := r.resolveStructs(schema.Structs)
	if err != nil {
		return nil, err
	}

	apis, err := r.resolveAPIs(schema.APIs)
	if err != nil {
		return nil, err
	}

	return &IRSchema{
		Structs: structs,
		Enums:   enums,
		APIs:    apis,
		Note:    schema.Note,
	}, nil
}
