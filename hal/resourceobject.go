// go2hal v0.2.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import (
	"reflect"
	"strings"

	"github.com/pmoule/go2hal/hal/relationtype"
)

// Resource is the root element of a HAL document.
// A Resource can
// - have links - AddLink(LinkRelation)
// - have CURIEs - AddCurieLinks([]*LinkObject)
// - embed other resources - AddResource(ResourceRelation)
type Resource interface {
	Data() PropertyMap
	Links() NamedMap
	EmbeddedResources() NamedMap
	AddData(interface{})
	AddLink(LinkRelation)
	AddResource(ResourceRelation)
	AddCurieLinks([]*LinkObject)
}

type resourceObject struct {
	data     PropertyMap       `json:"-"`
	links    links             `json:"_links,omitempty"`
	embedded embeddedResources `json:"_embedded,omitempty"`
}

// NewResourceObject initialises a valid Resource.
func NewResourceObject() Resource {
	return &resourceObject{data: PropertyMap{}, links: links{}, embedded: embeddedResources{}}
}

func (r *resourceObject) Data() PropertyMap {
	return r.data
}

func (r *resourceObject) AddData(data interface{}) {
	if data == nil {
		return
	}

	r.readDataFields(reflect.ValueOf(data))
}

func (r *resourceObject) readDataFields(v reflect.Value) {
	vType := v.Type()

	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
		v = v.Elem()
	}

	if vType.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < vType.NumField(); i++ {
		tField := vType.Field(i)
		vField := v.Field(i)

		if !vField.CanInterface() {
			continue
		}

		if tField.Anonymous {
			if !vField.CanAddr() {
				anonymValue := reflect.ValueOf(vField.Interface())
				r.readDataFields(anonymValue)
				continue
			}

			r.readDataFields(vField.Addr())
		}

		jsonValue, ok := tField.Tag.Lookup("json")

		if !ok || jsonValue == "-" {
			continue
		}

		tokens := strings.Split(jsonValue, ",")
		fieldName := tokens[0]
		omitEmpty := len(tokens) > 1 && strings.TrimSpace(tokens[1]) == "omitempty"
		isZeroValue := isZeroValue(vField) //value == reflect.Zero(reflect.TypeOf(value)).Interface()

		if omitEmpty && isZeroValue {
			continue
		}

		r.data[fieldName] = vField.Interface()
	}
}

func isZeroValue(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return val.IsNil()
	case reflect.Struct:
		isZero := true

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

	value := val.Interface()
	zeroValue := reflect.Zero(val.Type()).Interface()

	return value == zeroValue
}

func (r *resourceObject) Links() NamedMap {
	return r.links.ToMap()
}

func (r *resourceObject) EmbeddedResources() NamedMap {
	return r.embedded.ToMap()
}

func (r *resourceObject) ToMap() NamedMap {
	properties := PropertyMap{}

	namedMaps := []NamedMap{}
	namedMaps = append(namedMaps, r.Links())
	namedMaps = append(namedMaps, r.EmbeddedResources())

	for _, namedMap := range namedMaps {
		if len(namedMap.Content) > 0 {
			properties[namedMap.Name] = namedMap.Content
		}
	}

	for key, val := range r.data {
		properties[key] = val
	}

	return NamedMap{Name: "root", Content: properties}
}

func (r *resourceObject) AddCurieLinks(linkObjects []*LinkObject) {
	rel, _ := NewLinkRelation(relationtype.CURIES)
	rel.SetLinks(linkObjects)
	r.AddLink(rel)
}

func (r *resourceObject) AddLink(rel LinkRelation) {
	r.links[rel.Name()] = rel
}

func (r *resourceObject) AddResource(rel ResourceRelation) {
	r.embedded[rel.Name()] = rel
}
