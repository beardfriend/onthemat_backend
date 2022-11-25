package utils

import (
	"strings"

	"github.com/goccy/go-json"
)

func StructToMap[T any](data interface{}) (reqData map[string]T) {
	s, _ := json.Marshal(data)
	decoder := json.NewDecoder(strings.NewReader(string(s)))
	decoder.UseNumber()
	decoder.Decode(&reqData)
	return
}
