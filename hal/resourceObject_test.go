// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import (
	"testing"
	"reflect"
)

func TestAddLinkObject(t *testing.T) {
	want, _ := NewLinkRelation("relation")
	wantedName := "_links"

	resourceObject := NewResourceObject()
	resourceObject.AddLinkObject(want, &LinkObject{})

	namedMap := resourceObject.Links()

	if namedMap.Name != wantedName {
		t.Errorf("Map is named %s, want %s", namedMap.Name, wantedName)
	}

	_, ok := namedMap.Content[want.Name()]

	if !ok {
		t.Errorf("LinkRelation %q does not exist", want)
	}
}

func TestAddResourceObject(t *testing.T) {
	want, _ := NewLinkRelation("relation")
	wantedName := "_embedded"

	embeddedResource := NewResourceObject()
	resourceObject := NewResourceObject()
	resourceObject.AddResourceObject(want, embeddedResource)

	namedMap := resourceObject.EmbeddedResources()

	if namedMap.Name != wantedName {
		t.Errorf("Map is named %s, want %s", namedMap.Name, wantedName)
	}

	val, ok := namedMap.Content[want.Name()]

	if !ok {
		t.Errorf("LinkRelation %q does not exist", want)
	}

	_, isPropertyMap := val.(PropertyMap)

	if !isPropertyMap {
		t.Errorf("LinkRelation value is %[1]T(%[1]p), want %[2]T(%[2]p)", val, PropertyMap{})
	}
}

func TestAddCurieLink(t *testing.T) {
	resourceObject := NewResourceObject()
	curieName := "doc"
	curieHref := "http://doc/{rel}"
	curieLink, _ := NewCurieLink(curieName, curieHref)

	resourceObject.AddCurieLink(curieLink)

	val := resourceObject.Links().Content["curies"]
	result, ok := val.([]*LinkObject)

	if !ok {
		t.Errorf("CurieLink type == %q, want %q", reflect.TypeOf(val), reflect.TypeOf(result))
	}

	if count := len(result); count != 1 {
		t.Errorf("CurieLink count == %d, want %d", count, 1)
	}

	if curieLink != result[0] {
		t.Errorf("CurieLink == %q, want %q", val, curieLink)
	}

	curieLink2, _ := NewCurieLink(curieName, curieHref)
	resourceObject.AddCurieLink(curieLink2)
	val = resourceObject.Links().Content["curies"]
	result, _ = val.([]*LinkObject)

	if count := len(result); count != 2 {
		t.Errorf("CurieLink count == %d, want %d", count, 2)
	}
}