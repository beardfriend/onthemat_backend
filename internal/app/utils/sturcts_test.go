package utils

import (
	"fmt"
	"reflect"
	"testing"

	"onthemat/internal/app/transport/request"

	"github.com/fatih/structs"
)

func TestSturcts(t *testing.T) {
}

func TestStructLibs(t *testing.T) {
	req := &request.YogaUpdateBody{
		Level:       Int(2),
		Description: String("description"),
	}
	s := structs.New(req)
	for _, v := range s.Fields() {
		if v.Kind() == reflect.Pointer {
			if v.IsZero() {
				fmt.Println(v.Tag("json"))
			}
		}
	}
}

func BenchmarkReflect(b *testing.B) {
	b.Run("raw", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req := &request.YogaUpdateBody{
				NameKor:     "kor",
				Level:       Int(2),
				Description: String("description"),
			}

			var nilableFields []string

			// field := v.Type()
			v := reflect.ValueOf(req)
			v = reflect.Indirect(v)
			ts := reflect.TypeOf(req)
			for i := 0; i < v.NumField(); i++ {
				if v.Field(i).Kind() == reflect.Pointer {
					if v.Field(i).IsNil() {
						nilableFields = append(nilableFields, ts.Name())
					}
				}
			}
		}
	})

	b.Run("lib", func(b *testing.B) {
		req := &request.YogaUpdateBody{
			NameKor:     "kor",
			Level:       Int(2),
			Description: String("description"),
		}
		s := structs.New(req)
		for _, v := range s.Fields() {
			var nilableFields []string
			if v.Kind() == reflect.Pointer {
				if v.IsZero() {
					nilableFields = append(nilableFields, v.Name())
				}
			}
		}
	})
}

func String(s string) *string {
	result := &s
	return result
}

func Int(s int) *int {
	result := &s
	return result
}
