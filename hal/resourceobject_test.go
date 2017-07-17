// go2hal v0.3.0
// Copyright (c) 2017 Patrick Moule
// License: MIT

package hal

import (
	"reflect"
	"testing"
)

func TestNewResource(t *testing.T) {
	wanted := []interface{}{"_links", 0, "_embedded", 0, 0}

	resource := NewResourceObject()
	links := resource.Links()

	if name := links.Name; name != wanted[0] {
		t.Errorf("Links name is %s, want %s", name, wanted[0])
	}

	if count := len(links.Content); count != wanted[1] {
		t.Errorf("Initial link amount %d, want %d", count, wanted[1])
	}

	embeddedResources := resource.EmbeddedResources()

	if name := embeddedResources.Name; name != wanted[2] {
		t.Errorf("Embedded resource name is %s, want %s", name, wanted[2])
	}

	if count := len(embeddedResources.Content); count != wanted[3] {
		t.Errorf("Initial embeddedResources amount %d, want %d", count, wanted[3])
	}

	data := resource.Data()

	if count := len(data); count != wanted[4] {
		t.Errorf("Initial data amount %d, want %d", count, wanted[4])
	}
}

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

	resourceObject.AddCurieLinks([]*LinkObject{curieLink})

	val := resourceObject.Links().Content["curies"]
	result, ok := val.([]*LinkObject)

	if !ok {
		t.Errorf("CurieLink type == %q, want %q", reflect.TypeOf(val), reflect.TypeOf(result))
	}

	if count := len(result); count != 1 {
		t.Errorf("CurieLink count == %d, want %d", count, 1)
	}

	if curieLink != result[0] {
		t.Errorf("CurieLink == %+v, want %+v", val, curieLink)
	}

	curieLink2, _ := NewCurieLink(curieName, curieHref)
	resourceObject.AddCurieLinks([]*LinkObject{curieLink, curieLink2})

	val2 := resourceObject.Links().Content["curies"]
	result2, _ := val2.([]*LinkObject)

	if count := len(result2); count != 2 {
		t.Errorf("CurieLink count == %d, want %d", count, 2)
	}
}

func TestAddData(t *testing.T) {
	resource := NewResourceObject()

	data := resource.Data()

	if count := len(data); count != 0 {
		t.Errorf("Initial data amount %d, want %d", count, 0)
	}

	data["test"] = "test"

	if count := len(data); count != 1 {
		t.Errorf("Data amount %d, want %d", count, 1)
	}

	delete(data, "test")

	if count := len(data); count != 0 {
		t.Errorf("Data amount %d, want %d", count, 0)
	}

	resource.AddData(nil)
	resource.AddData("test")
	resource.AddData(true)
	resource.AddData(1)

	if count := len(data); count != 0 {
		t.Errorf("Data amount %d, want %d", count, 0)
	}

	type Test3 struct {
		J string `json:"j, omitempty"`
	}

	type Test2 struct {
		F string `json:"f"`
		g string
		H Test3  `json:"h"`
		I [2]int `json:"i, omitempty"`
	}

	type Test1 struct {
		Test2
		A string   `json:"a"`
		B []string `json:"b"`
		c string
		D int `json:"d, omitempty"`
		E int `json:"-"`
	}

	test2 := new(Test2)
	test2.F = "F"
	test2.g = "G"
	test := Test1{*test2, "A", []string{"B"}, "C", 0, 1}
	resource.AddData(test)

	if count := len(data); count != 4 {
		t.Errorf("Data amount %d, want %d", count, 4)
	}

	if val, ok := data["a"]; !ok && val != "A" {
		t.Errorf("Expected key %s with value %s in data", "a", "A")
	}

	resource.AddData(&test)

	if count := len(data); count != 4 {
		t.Errorf("Data amount %d, want %d", count, 4)
	}

	if val, ok := data["f"]; !ok && val != "F" {
		t.Errorf("Expected key %s with value %s in data", "f", "F")
	}
}
