// go2hal v0.3.2
// Copyright (c) 2017 Patrick Moule
// License: MIT

package mapping

import (
	"reflect"
	"strings"
	"time"
)

// NamedMap simply links a name with PropertyMap
type NamedMap struct {
	Name    string
	Content PropertyMap
}

// PropertyMap simply maps a string to any kind of value.
type PropertyMap map[string]interface{}

func MapData(data interface{}) PropertyMap {
	return readDataFields(reflect.ValueOf(data))
}

func readDataFields(v reflect.Value) PropertyMap {
	if isZeroValue(v) {
		return PropertyMap{}
	}

	vType := v.Type()

	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
		v = v.Elem()
	}

	if vType.Kind() != reflect.Struct {
		return PropertyMap{}
	}

	propertyMap := PropertyMap{}

	for i := 0; i < vType.NumField(); i++ {
		vField := v.Field(i)

		if !vField.CanInterface() {
			continue
		}

		tField := vType.Field(i)

		if tField.Anonymous {
			value := readEmbeddedField(vField)

			for key, v := range value {
				propertyMap[key] = v
			}

			continue
		}

		if fieldName, value, ok := toJSONValue(tField, vField); ok {
			propertyMap[fieldName] = value
		}
	}

	return propertyMap
}

func toJSONValue(tField reflect.StructField, vField reflect.Value) (string, interface{}, bool) {
	jsonValue, ok := tField.Tag.Lookup("json")

	if !ok || jsonValue == "-" {
		return "", nil, false
	}

	tokens := strings.Split(jsonValue, ",")
	omitEmpty := len(tokens) > 1 && strings.TrimSpace(tokens[1]) == "omitempty"
	fieldName := tokens[0]

	_, isTime := vField.Interface().(time.Time)

	if vField.Kind() == reflect.Struct && !isTime {
		value := readDataFields(reflect.ValueOf(vField.Interface()))
		isZeroValue := len(value) == 0

		if omitEmpty && isZeroValue {
			return "", nil, false
		}

		return fieldName, value, true
	}

	isZeroValue := isZeroValue(vField)

	if omitEmpty && isZeroValue {
		return "", nil, false
	}

	return fieldName, vField.Interface(), true
}

func readEmbeddedField(v reflect.Value) PropertyMap {
	if isZeroValue(v) {
		return PropertyMap{}
	}

	if !v.CanAddr() {
		value := reflect.ValueOf(v.Interface())
		return readDataFields(value)
	}

	return readDataFields(v.Addr())
}

func isZeroValue(val reflect.Value) bool {
	if val == reflect.Zero(reflect.TypeOf(val)).Interface() {
		return true
	}

	switch val.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return val.IsNil()
	case reflect.Struct:
		isZero := true

		if value, ok := val.Interface().(time.Time); ok {
			return value.IsZero()
		}

		for i := 0; i < val.NumField(); i++ {
			isZero = isZero && isZeroValue(val.Field(i))
		}

		return isZero
	case reflect.Array:
		isZero := true

		for i := 0; i < val.Len(); i++ {
			isZero = isZero && isZeroValue(val.Index(i))
		}

		return isZero
	}

	if val.CanInterface() {
		value := val.Interface()
		zeroValue := reflect.Zero(val.Type()).Interface()
		return value == zeroValue
	}

	return true
}
