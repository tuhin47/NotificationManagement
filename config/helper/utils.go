package helper

import (
	"reflect"
	"strconv"
)

func ToInt(s string) *int {
	if s == "" {
		return nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &i
}

func ToBool(s string) *bool {
	if s == "" {
		return nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil
	}
	return &b
}

func FirstNonEmpty[T any](values ...T) *T {
	for _, value := range values {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Ptr:
			if !v.IsNil() {
				return &value
			}
		case reflect.String:
			if v.String() != "" {
				return &value
			}
		default:
			return &value
		}
	}
	return nil
}
