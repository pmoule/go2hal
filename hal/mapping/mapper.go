// go2hal v0.4.0
// Copyright (c) 2020 Patrick Moule
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
		if v == reflect.Zero(reflect.TypeOf(v)).Interface() {
			return PropertyMap{}
		}
	}

	vType := v.Type()

	// force the real type
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

		// force the real type
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
		if vField.Kind() == reflect.Slice {
			element := tField.Type.Elem()

			// force the real type
			if element.Kind() == reflect.Ptr {
				element = element.Elem()
			}

			if element.Kind() == reflect.Struct {
				sliceValuePtr := createSlice([]PropertyMap{})

				for i := 0; i < vField.Len(); i++ {
					v := vField.Index(i)
					value := readDataFields(v)
					sliceValuePtr.Set(reflect.Append(sliceValuePtr, reflect.ValueOf(value)))
				}

				return fieldName, sliceValuePtr.Interface(), true
			} else {
				sliceValuePtr := createSlice([]string{})

				for i := 0; i < vField.Len(); i++ {
					v := vField.Index(i)
					vType := v.Type()

					if vType.Kind() == reflect.Ptr {
						v = v.Elem()
					}

					if m, ok := v.Addr().Interface().(json.Marshaler); ok {
						b, err := m.MarshalJSON()

						if err != nil {
							continue
						}

						value := string(b)
						isZeroValue := len(value) == 0

						if isZeroValue {
							continue
						}

						sliceValuePtr.Set(reflect.Append(sliceValuePtr, reflect.ValueOf(value)))
					} else {
						sliceValuePtr.Set(reflect.Append(sliceValuePtr, v))
					}
				}

				if isZeroValue(sliceValuePtr) && omitEmpty {
					return "", nil, false
				}

				return fieldName, sliceValuePtr.Interface(), true
			}
		}

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

			if isZeroValue && omitEmpty {
				return "", nil, false
			}

			return fieldName, value, true
		}

		if vField.Kind() == reflect.Struct {
			value := readDataFields(reflect.ValueOf(vField.Interface()))
			isZeroValue := len(value) == 0

			if isZeroValue && omitEmpty {
				return "", nil, false
			}

			return fieldName, value, true
		}
	}

	if isZeroValue(vField) && omitEmpty {
		return "", nil, false
	}

	return fieldName, vField.Interface(), true
}

func createSlice(sliceType interface{}) reflect.Value {
	reflection := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(sliceType).Elem()), 0, 0)
	reflectionValue := reflect.New(reflection.Type())
	reflectionValue.Elem().Set(reflection)
	slicePtr := reflect.ValueOf(reflectionValue.Interface())
	sliceValuePtr := slicePtr.Elem()

	return sliceValuePtr
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
	case reflect.Func:
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
	case reflect.Array, reflect.Slice:
		isZero := true

		for i := 0; i < val.Len(); i++ {
			isZero = isZero && isZeroValue(val.Index(i))
		}

		return isZero
	case reflect.Map:
		isZero := true

		for _, e := range val.MapKeys() {
			isZero = isZero && isZeroValue(val.MapIndex(e))
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
