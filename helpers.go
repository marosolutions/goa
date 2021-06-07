package goa

import (
	"reflect"
	"strings"
)

// contains checks if a string is present in a slice
func SliceContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// DeleteEmpty will remove any empty strings from a slice
func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		str = strings.TrimSpace(str)
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// StructToKeysAndValues: Convert a struct to Keys and Values mainly for SQL Inserts
func StructToKeysAndValues(tag string, s interface{}) (keys []string, vals []interface{}) {
	v := reflect.ValueOf(s)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		// cols = append(cols, typeOfS.Field(i).Name)
		keys = append(keys, typeOfS.Field(i).Tag.Get(tag))
		vals = append(vals, v.Field(i).Interface())
	}
	return
}
