package utils

import (
	"reflect"
	"strings"
)

func TableName(table interface{}) string {
	if t, ok := table.(interface {
		Name() string
	}); ok {
		return t.Name()
	}

	t := reflect.TypeOf(table)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return toSnakeCase(t.Name())
}

func toSnakeCase(s string) string {
	return strings.NewReplacer(" ", "_").Replace(s)
}
