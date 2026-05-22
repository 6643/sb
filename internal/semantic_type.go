package internal

import "fmt"

func (r *resolver) resolveType(t TypeRef, allowNil bool) (IRType, error) {
	if t.Name == "" {
		return IRType{}, fmt.Errorf("类型名为空")
	}

	if t.Name == "nil" {
		if !allowNil || t.IsList {
			return IRType{}, fmt.Errorf("nil 仅允许作为 API 非数组返回类型")
		}
		return IRType{Name: "nil", Kind: IRKindNil}, nil
	}

	if _, ok := baseTypes[t.Name]; ok {
		return IRType{Name: t.Name, Kind: IRKindBase, IsList: t.IsList}, nil
	}
	if _, ok := r.structDefs[t.Name]; ok {
		return IRType{Name: t.Name, Kind: IRKindStruct, IsList: t.IsList}, nil
	}
	if _, ok := r.enumDefs[t.Name]; ok {
		return IRType{Name: t.Name, Kind: IRKindEnum, IsList: t.IsList}, nil
	}

	return IRType{}, fmt.Errorf("未定义类型: %s", t.Name)
}
