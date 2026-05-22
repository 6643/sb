package internal

import "fmt"

func (r *resolver) collectDefs(schema *Schema) error {
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
