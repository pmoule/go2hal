// go2hal v0.3.3
// Copyright (c) 2017 Patrick Moule
// License: MIT

package mapping

import (
	"encoding/json"
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

// MapData returns a PropertyMap for given data.
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
		tField := vType.Field(i)

		if vField.Kind() == reflect.Ptr {
			vField = vField.Elem()
		}

		if !vField.IsValid() {
			fieldName, omitEmpty, ok := readJSONInfo(tField)

			if !ok || omitEmpty {
				continue
			}

			propertyMap[fieldName] = nil
			continue
		}

		if !vField.CanInterface() {
			continue
		}

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

func readJSONInfo(tField reflect.StructField) (string, bool, bool) {
	jsonValue, ok := tField.Tag.Lookup("json")

	if !ok || jsonValue == "-" {
		return "", true, false
	}

	tokens := strings.Split(jsonValue, ",")
	omitEmpty := len(tokens) > 1 && strings.TrimSpace(tokens[1]) == "omitempty"
	fieldName := tokens[0]

	return fieldName, omitEmpty, true
}

func toJSONValue(tField reflect.StructField, vField reflect.Value) (string, interface{}, bool) {
	fieldName, omitEmpty, ok := readJSONInfo(tField)

	if !ok {
		return "", nil, false
	}

	_, isTime := vField.Interface().(time.Time)

	if !isTime {
		va := vField

		if tField.Type.Kind() == reflect.Ptr {
			va = vField.Addr()
		}

		if m, ok := va.Interface().(json.Marshaler); ok {
			b, err := m.MarshalJSON()

			if err != nil {
				return "", nil, false
			}

			value := string(b)
			isZeroValue := len(value) == 0

			if omitEmpty && isZeroValue {
				return "", nil, false
			}

			return fieldName, value, true
		}
	}

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
