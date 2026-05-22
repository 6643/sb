package internal

import "fmt"

func (r *resolver) resolveAPIs(defs []API) ([]IRAPI, error) {
	result := make([]IRAPI, 0, len(defs))
	apiNames := make(map[string]struct{}, len(defs))

	for _, api := range defs {
		if _, ok := apiNames[api.Name]; ok {
			return nil, fmt.Errorf("API 重复定义: %s", api.Name)
		}
		apiNames[api.Name] = struct{}{}

		args := make([]IRAPIArg, 0, len(api.Args))
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
			args = append(args, IRAPIArg{Name: a.Name, Type: t})
		}

		resType, err := r.resolveType(api.Result, true)
		if err != nil {
			return nil, fmt.Errorf("API %s 返回类型: %w", api.Name, err)
		}

		result = append(result, IRAPI{Name: api.Name, Args: args, Result: resType, Note: api.Note})
	}

	return result, nil
}
