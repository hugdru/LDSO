package generators

func GenerateSearchClause(filter map[string]string) (string, []interface{}) {
	if filter == nil {
		return "", nil
	}
	filterLength := len(filter)
	if filterLength == 0 {
		return "", nil
	}
	var where string
	values := make([]interface{}, filterLength)
	offset := filterLength - 1
	for name, value := range filter {
		values[offset] = value
		where = name + " = ? " + where
		if offset == 0 {
			where = "WHERE " + where
		} else {
			where = " AND " + where
		}
		offset--
	}

	return where, values
}
