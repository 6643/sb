package semantic

import (
	"fmt"

	"sb/internal/ast"
	"sb/internal/ir"
)

var baseTypes = map[string]struct{}{
	"i8": {}, "u8": {}, "i16": {}, "u16": {},
	"i32": {}, "u32": {}, "i64": {}, "u64": {},
	"f32": {}, "f64": {}, "bool": {}, "text": {}, "bin": {},
}

// Resolve 把 AST 解析为可直接代码生成的 IR。
func Resolve(schema *ast.Schema) (*ir.Schema, error) {
	resolver := &resolver{
		structDefs: make(map[string]ast.Struct),
		enumDefs:   make(map[string]ast.Enum),
		fieldCache: make(map[string][]ir.Field),
	}
	return resolver.resolve(schema)
}

type resolver struct {
	structDefs map[string]ast.Struct
	enumDefs   map[string]ast.Enum
	fieldCache map[string][]ir.Field
}

func (r *resolver) resolve(schema *ast.Schema) (*ir.Schema, error) {
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

	return &ir.Schema{
		Structs: structs,
		Enums:   enums,
		APIs:    apis,
		Note:    schema.Note,
	}, nil
}

func (r *resolver) collectDefs(schema *ast.Schema) error {
	for _, st := range schema.Structs {
		if _, ok := r.structDefs[st.Name]; ok {
			return fmt.Errorf("结构体重复定义: %s", st.Name)
		}
		if _, ok := r.enumDefs[st.Name]; ok {
			return fmt.Errorf("名称冲突: %s", st.Name)
		}
		r.structDefs[st.Name] = st
	}

	for _, e := range schema.Enums {
		if _, ok := r.enumDefs[e.Name]; ok {
			return fmt.Errorf("枚举重复定义: %s", e.Name)
		}
		if _, ok := r.structDefs[e.Name]; ok {
			return fmt.Errorf("名称冲突: %s", e.Name)
		}
		r.enumDefs[e.Name] = e
	}

	return nil
}

func (r *resolver) resolveEnums(defs []ast.Enum) ([]ir.Enum, error) {
	result := make([]ir.Enum, 0, len(defs))
	for _, e := range defs {
		memberNames := make(map[string]struct{})
		memberIDs := make(map[uint8]struct{})
		members := make([]ir.EnumMember, 0, len(e.Members))

		var nextID uint8
		hasNext := false
		for _, m := range e.Members {
			if _, ok := memberNames[m.Name]; ok {
				return nil, fmt.Errorf("枚举 %s 成员重名: %s", e.Name, m.Name)
			}
			memberNames[m.Name] = struct{}{}

			id, ok, err := r.resolveEnumID(m.Value, nextID, hasNext)
			if err != nil {
				return nil, fmt.Errorf("枚举 %s 成员 %s: %w", e.Name, m.Name, err)
			}
			if _, exists := memberIDs[id]; exists {
				return nil, fmt.Errorf("枚举 %s 成员 ID 重复: %d", e.Name, id)
			}
			memberIDs[id] = struct{}{}

			nextID = id
			hasNext = ok
			members = append(members, ir.EnumMember{ID: id, Name: m.Name, Note: m.Note})
		}

		result = append(result, ir.Enum{Name: e.Name, Members: members, Note: e.Note})
	}
	return result, nil
}

func (r *resolver) resolveEnumID(explicit *uint8, lastID uint8, hasLast bool) (uint8, bool, error) {
	if explicit != nil {
		return *explicit, true, nil
	}
	if !hasLast {
		return 0, true, nil
	}
	if lastID == 255 {
		return 0, true, fmt.Errorf("枚举值溢出")
	}
	return lastID + 1, true, nil
}

func (r *resolver) resolveStructs(defs []ast.Struct) ([]ir.Struct, error) {
	result := make([]ir.Struct, 0, len(defs))
	for _, st := range defs {
		fields, err := r.expandStructFields(st.Name, make(map[string]bool))
		if err != nil {
			return nil, err
		}
		if len(fields) > 255 {
			return nil, fmt.Errorf("结构体 %s 字段数超过 255", st.Name)
		}
		if err := checkDuplicateFieldNames(st.Name, fields); err != nil {
			return nil, err
		}
		result = append(result, ir.Struct{Name: st.Name, Fields: fields, Note: st.Note})
	}
	return result, nil
}

func (r *resolver) expandStructFields(structName string, visiting map[string]bool) ([]ir.Field, error) {
	if cached, ok := r.fieldCache[structName]; ok {
		return cached, nil
	}
	if visiting[structName] {
		return nil, fmt.Errorf("检测到循环嵌入: %s", structName)
	}

	st, ok := r.structDefs[structName]
	if !ok {
		return nil, fmt.Errorf("未定义结构体: %s", structName)
	}

	visiting[structName] = true
	result := make([]ir.Field, 0, len(st.Fields))
	for _, f := range st.Fields {
		if !f.Embedded {
			t, err := r.resolveType(f.Type, false)
			if err != nil {
				return nil, fmt.Errorf("结构体 %s 字段 %s: %w", structName, f.Name, err)
			}
			result = append(result, ir.Field{Name: f.Name, Type: t, Tag: f.Tag, Note: f.Note})
			continue
		}

		if f.Type.IsList {
			return nil, fmt.Errorf("结构体 %s: 嵌入字段不能是数组", structName)
		}
		if _, ok := r.structDefs[f.Type.Name]; !ok {
			return nil, fmt.Errorf("结构体 %s: 嵌入类型 %s 不是结构体", structName, f.Type.Name)
		}

		expanded, err := r.expandStructFields(f.Type.Name, visiting)
		if err != nil {
			return nil, err
		}
		result = append(result, expanded...)
	}
	visiting[structName] = false

	r.fieldCache[structName] = result
	return result, nil
}

func checkDuplicateFieldNames(structName string, fields []ir.Field) error {
	seen := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		if _, ok := seen[f.Name]; ok {
			return fmt.Errorf("结构体 %s 展开后字段重名: %s", structName, f.Name)
		}
		seen[f.Name] = struct{}{}
	}
	return nil
}

func (r *resolver) resolveAPIs(defs []ast.API) ([]ir.API, error) {
	result := make([]ir.API, 0, len(defs))
	apiNames := make(map[string]struct{}, len(defs))

	for _, api := range defs {
		if _, ok := apiNames[api.Name]; ok {
			return nil, fmt.Errorf("API 重复定义: %s", api.Name)
		}
		apiNames[api.Name] = struct{}{}

		args := make([]ir.APIArg, 0, len(api.Args))
		argNames := make(map[string]struct{}, len(api.Args))
		for _, a := range api.Args {
			if _, ok := argNames[a.Name]; ok {
				return nil, fmt.Errorf("API %s 参数重名: %s", api.Name, a.Name)
			}
			argNames[a.Name] = struct{}{}

			t, err := r.resolveType(a.Type, false)
			if err != nil {
				return nil, fmt.Errorf("API %s 参数 %s: %w", api.Name, a.Name, err)
			}
			args = append(args, ir.APIArg{Name: a.Name, Type: t})
		}

		resType, err := r.resolveType(api.Result, true)
		if err != nil {
			return nil, fmt.Errorf("API %s 返回类型: %w", api.Name, err)
		}

		result = append(result, ir.API{Name: api.Name, Args: args, Result: resType, Note: api.Note})
	}

	return result, nil
}

func (r *resolver) resolveType(t ast.TypeRef, allowNil bool) (ir.Type, error) {
	if t.Name == "" {
		return ir.Type{}, fmt.Errorf("类型名为空")
	}

	if t.Name == "nil" {
		if !allowNil || t.IsList {
			return ir.Type{}, fmt.Errorf("nil 仅允许作为 API 非数组返回类型")
		}
		return ir.Type{Name: "nil", Kind: ir.KindNil}, nil
	}

	if _, ok := baseTypes[t.Name]; ok {
		return ir.Type{Name: t.Name, Kind: ir.KindBase, IsList: t.IsList}, nil
	}
	if _, ok := r.structDefs[t.Name]; ok {
		return ir.Type{Name: t.Name, Kind: ir.KindStruct, IsList: t.IsList}, nil
	}
	if _, ok := r.enumDefs[t.Name]; ok {
		return ir.Type{Name: t.Name, Kind: ir.KindEnum, IsList: t.IsList}, nil
	}

	return ir.Type{}, fmt.Errorf("未定义类型: %s", t.Name)
}
