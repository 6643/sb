package internal

import "fmt"

func (r *resolver) resolveEnums(defs []Enum) ([]IREnum, error) {
	result := make([]IREnum, 0, len(defs))
	for _, e := range defs {
		memberNames := make(map[string]struct{})
		memberIDs := make(map[uint8]struct{})
		members := make([]IREnumMember, 0, len(e.Members))

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
			members = append(members, IREnumMember{ID: id, Name: m.Name, Note: m.Note})
		}

		result = append(result, IREnum{Name: e.Name, Members: members, Note: e.Note})
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
