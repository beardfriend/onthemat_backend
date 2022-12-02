package utils

import (
	"encoding/json"
	"strconv"
	"strings"

	"onthemat/internal/app/transport"
	"onthemat/pkg/ent"

	"github.com/fatih/structs"
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

func GetUpdateableDataV2(s *structs.Struct, columns []string) (result map[string]interface{}) {
	result = make(map[string]interface{}, 0)
	for _, f := range s.Fields() {

		name := strcase.ToSnake(f.Name())

		isUpdateable := Contains(columns, name)
		if !isUpdateable {
			continue
		}

		switch f.Value().(type) {

		case string:
			value := f.Value().(string)
			result[name] = value

		case *string:
			ptrValue := f.Value().(*string)
			if ptrValue != nil {
				value := *ptrValue
				result[name] = value
			}

		case int:
			value := f.Value().(int)
			result[name] = value

		case *int:
			ptrValue := f.Value().(*int)
			if ptrValue != nil {
				value := *ptrValue
				result[name] = value
			}
		case bool:
			value := f.Value().(bool)
			result[name] = value

		case *bool:
			ptrValue := f.Value().(*bool)
			if ptrValue != nil {
				value := *ptrValue
				result[name] = value
			}

		case transport.TimeString:
			value := f.Value().(transport.TimeString)
			result[name] = value

		case *transport.TimeString:
			ptrValue := f.Value().(*transport.TimeString)
			if ptrValue != nil {
				value := *ptrValue
				result[name] = value
			}
		}
	}
	return
}

func MakeDataForCondition(requestIds []int, existIds []int) (createable []int, updateable []int, deleteable []int) {
	updateable = Intersection(requestIds, existIds)
	deleteable = Difference(existIds, requestIds)
	createable = Difference(requestIds, existIds)
	return
}
