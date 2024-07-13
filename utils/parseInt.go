package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var ErrUnSupportType = errors.New("parse unsupported type")

func ParseInterfaceToInt(value interface{}) (int, error) {
	v := reflect.ValueOf(value)

	// nolint:exhaustive
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
		return 0, fmt.Errorf("%w: %v", ErrUnSupportType, v.Kind())
	}
}
