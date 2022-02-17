package smode

import (
	"reflect"
)

// Source: https://newbedev.com/function-for-converting-a-struct-to-map-in-golang
func structmap(in interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		panic("StructMap only accepts structs")
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if fi.Tag.Get("json") == "-" {
			continue
		}
		// set key of map to value in struct field
		out[fi.Name] = v.Field(i).Interface()
	}
	return out
}
