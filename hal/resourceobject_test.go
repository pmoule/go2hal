// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import (
	"testing"
	"fmt"
	"reflect"
)

func TestAddLinkObject(t *testing.T) {
	want, _ := NewLinkRelation("relation")
	want.SetLink(&LinkObject{})
	wantedName := "_links"

	resourceObject := NewResourceObject()
	resourceObject.AddLink(want)

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
	want, _ := NewResourceRelation("relation")
	want.SetResource(NewResourceObject())
	wantedName := "_embedded"

	resourceObject := NewResourceObject()
	resourceObject.AddResource(want)

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

	resourceObject.AddCurieLinks([]*LinkObject {curieLink})

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
	resourceObject.AddCurieLinks([]*LinkObject {curieLink, curieLink2})

	fmt.Println("length: " + string(len(resourceObject.Links().Content)))
	val2 := resourceObject.Links().Content["curies"]
	result2, _ := val2.([]*LinkObject)

	if count := len(result2); count != 2 {
		t.Errorf("CurieLink count == %d, want %d", count, 2)
	}
}