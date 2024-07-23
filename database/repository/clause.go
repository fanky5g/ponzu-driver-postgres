package repository

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

func (repo *Repository) valueType(value interface{}) (string, error) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Interface().(type) {
	case string:
		return "text", nil
	case bool:
		return "boolean", nil
	case int:
		return "numeric", nil
	case int64:
		return "numeric", nil
	case int32:
		return "numeric", nil
	case int8:
		return "numeric", nil
	case int16:
		return "numeric", nil
	case float64:
		return "numeric", nil
	case float32:
		return "numeric", nil
	case time.Time:
		return "timestamptz", nil
	default:
		return "", fmt.Errorf("unsupported type %T", value)
	}
}

func (repo *Repository) getComparisonOperator(operator string) (string, error) {
	switch operator {
	case Equal:
        fallthrough
	case LessThan:
        fallthrough
	case LessThanOrEqualTo:
        fallthrough
	case GreaterThan:
        fallthrough
	case GreaterThanOrEqualTo:
        return operator, nil
	default:
		return "", errors.New("unsupported comparison operator")
	}
}
