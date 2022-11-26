package utils

import (
	"encoding/json"
	"strconv"
	"strings"

	"onthemat/pkg/ent"

	"github.com/iancoleman/strcase"
)

// 디코딩 과정을 거치면서 비어있는 데이터는 버려진다.
func structToMap[T any](data interface{}) (result map[string]T) {
	s, _ := json.Marshal(data)
	decoder := json.NewDecoder(strings.NewReader(string(s)))
	decoder.UseNumber()
	decoder.Decode(&result)
	return
}

// Struct 구조체에 Nillable로 선언되어 있으면 nil을 뱉는다.
func GetUpdateableData(data interface{}, allowColumns []string) (result map[string]ent.Value) {
	updatableData := structToMap[ent.Value](data)

	result = make(map[string]ent.Value, 0)

	for _, col := range allowColumns {
		columnCamelCase := strcase.ToLowerCamel(col)
		value := updatableData[columnCamelCase]

		if value == nil {
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
