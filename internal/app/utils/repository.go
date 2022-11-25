package utils

import (
	"encoding/json"
	"strconv"

	"onthemat/pkg/ent"

	"github.com/iancoleman/strcase"
)

// Struct 구조체에 Nillable로 선언되어 있으면 nil을 뱉는다.
func MakeUseableFieldWithData(data map[string]ent.Value, allowColumns []string) (result map[string]ent.Value) {
	result = make(map[string]ent.Value, 0)

	for _, col := range allowColumns {
		columnCamelCase := strcase.ToLowerCamel(col)
		value := data[columnCamelCase]

		if value == nil {
			continue
		}

		if value == "" {
			result[col] = nil
			continue
		}

		// 숫자 유형이면,
		if id, ok := value.(json.Number); ok {
			d, _ := strconv.Atoi(id.String())
			result[col] = d
			continue
		}

		result[col] = value

	}
	return
}
