package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

func ParseInterfaceToInt(value interface{}) (int, error) {
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int(v.Float()), nil
	case reflect.String:
		return strconv.Atoi(v.String())
	default:
		return 0, fmt.Errorf("unsupported type: %v", v.Kind())
	}
}
