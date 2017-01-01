package generators

import "bytes"

func GenerateAndSearchClause(filter map[string]interface{}) (string, []interface{}) {
	return GenerateSearchClause(filter, true)
}

func GenerateOrSearchClause(filter map[string]interface{}) (string, []interface{}) {
	return GenerateSearchClause(filter, false)
}

func GenerateSearchClause(filter map[string]interface{}, withAnd bool) (string, []interface{}) {
	const WHERE = " WHERE "
	const PREPARE = " = ? "
	const AND = "AND "
	const OR = "OR "

	if filter == nil {
		return "", nil
	}

	filterLength := len(filter)
	if filterLength == 0 {
		return "", nil
	}
	lastIndex := filterLength - 1

	separator := AND
	if !withAnd {
		separator = OR
	}

	var whereBuffer bytes.Buffer
	whereBuffer.WriteString(WHERE)

	values := make([]interface{}, filterLength)
	index := 0
	for name, value := range filter {
		values[index] = value
		whereBuffer.WriteString(name)
		whereBuffer.WriteString(PREPARE)
		if index != lastIndex {
			whereBuffer.WriteString(separator)
		}
		index++
	}

	return whereBuffer.String(), values
}

func GenerateIn(ids []int64) (string, []interface{}) {
	if ids == nil {
		return "", nil
	}

	values := make([]interface{}, len(ids))
	var inBuffer bytes.Buffer

	for index := len(ids) - 1; index != 0; index-- {
		values[index] = ids[index]
		inBuffer.WriteString("?,")
	}
	inBuffer.WriteString("?")
	values[0] = ids[0]

	return inBuffer.String(), values
}
