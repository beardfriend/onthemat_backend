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
		if data[strcase.ToLowerCamel(col)] != nil {

			id, ok := data[strcase.ToLowerCamel(col)].(json.Number)
			if ok {
				d, _ := strconv.Atoi(id.String())
				data[strcase.ToLowerCamel(col)] = d
			}

			data := data[strcase.ToLowerCamel(col)]

			if data == "" {
				result[col] = nil
			} else {
				result[col] = data
			}
		}
	}
	return
}
