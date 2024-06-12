package helper

import (
	"fmt"
	"reflect"
)

func CopyNonEmptyFields(src, dst interface{}) {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst).Elem() // .Elem() to get the value pointed to by the pointer

	srcType := srcVal.Type()

	// Ensure the src and dst are structs
	if srcType.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
		fmt.Println("The provided objects are not structs")
		return
	}

	// Iterate over the fields of the src struct
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dstVal.FieldByName(srcType.Field(i).Name)

		// Get the zero value for the field's type
		zeroValue := reflect.Zero(srcField.Type()).Interface()
		currentValue := srcField.Interface()

		// Copy the value if the source field is not empty
		if !reflect.DeepEqual(currentValue, zeroValue) {
			if dstField.CanSet() {
				dstField.Set(srcField)
			}
		}
	}
}
