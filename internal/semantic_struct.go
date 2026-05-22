package internal

import "fmt"

func (r *resolver) resolveStructs(defs []Struct) ([]IRStruct, error) {
	result := make([]IRStruct, 0, len(defs))
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
		result = append(result, IRStruct{Name: st.Name, Fields: fields, Note: st.Note})
	}
	return result, nil
}

func (r *resolver) expandStructFields(structName string, visiting map[string]bool) ([]IRField, error) {
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
	result := make([]IRField, 0, len(st.Fields))
	for _, f := range st.Fields {
		if !f.Embedded {
			t, err := r.resolveType(f.Type, false)
			if err != nil {
				return nil, fmt.Errorf("结构体 %s 字段 %s: %w", structName, f.Name, err)
			}
			result = append(result, IRField{Name: f.Name, Type: t, Tag: f.Tag, Note: f.Note})
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

func checkDuplicateFieldNames(structName string, fields []IRField) error {
	seen := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		if _, ok := seen[f.Name]; ok {
			return fmt.Errorf("结构体 %s 展开后字段重名: %s", structName, f.Name)
		}
		seen[f.Name] = struct{}{}
	}
	return nil
}
