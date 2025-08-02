package models

import (
	"reflect"

	"gorm.io/gorm"
)

type ModelInterface interface {
	UpdateFromModel(source ModelInterface)
}

// copyFields copies fields from source to destination, skipping gorm.Model fields.
func copyFields(destination, source interface{}) {
	destValue := reflect.ValueOf(destination).Elem()
	sourceValue := reflect.ValueOf(source).Elem()

	for i := 0; i < sourceValue.NumField(); i++ {
		sourceField := sourceValue.Type().Field(i)
		if sourceField.Anonymous && sourceField.Type == reflect.TypeOf(gorm.Model{}) {
			// Skip gorm.Model embedded struct
			continue
		}
		if tag := sourceField.Tag.Get("mapper"); tag == "ignore" {
			continue // skip editing
		}

		destField := destValue.FieldByName(sourceField.Name)
		if destField.IsValid() && destField.CanSet() {
			if tag := sourceField.Tag.Get("mapper"); tag == "inherit" {
				copyFields(destField.Addr().Interface(), sourceValue.Field(i).Addr().Interface())
			} else {
				destField.Set(sourceValue.Field(i))
			}
		}
	}
}
