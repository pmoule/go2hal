// go2hal v0.4.0
// Copyright (c) 2020 Patrick Moule
// License: MIT

package mapping

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestMapDataWithPrimitives(t *testing.T) {
	values := []interface{}{
		nil, "test", true, false,
		0, 1, -1, 1.5, time.Now(), []int{1, 2},
	}

	for _, value := range values {
		data := MapData(value)

		if count := len(data); count > 0 {
			t.Errorf("Data amount %d, want %d", count, 0)
		}
	}
}

func TestMapData(t *testing.T) {
	var data PropertyMap

	type Test7 struct {
	}

	type Test3 struct {
		K string `json:"k,omitempty"`
	}

	type Test2 struct {
		F string `json:"f"`
		g string
		H Test3  `json:"h"`
		I [2]int `json:"i,omitempty"`
		J Test3  `json:"j,omitempty"`
	}

	type Test1 struct {
		Test7
		Test2
		A string   `json:"a"`
		B []string `json:"b"`
		c string
		D int `json:"d,omitempty"`
		E int `json:"-"`
	}

	test2 := new(Test2)
	test2.F = "F"
	test2.g = "G"
	test := Test1{Test2: *test2, A: "A", B: []string{"B"}, c: "C", D: 0, E: 1}

	data = MapData(test)

	if count := len(data); count != 4 {
		t.Errorf("Data amount %d, want %d", count, 4)
	}

	if val, ok := data["a"]; !ok && val != "A" {
		t.Errorf("Expected key %s with value %s in data", "a", "A")
	}

	if val, ok := data["b"]; !ok && val != [1]string{"B"} {
		t.Errorf("Expected key %s with value %s in data", "a", [1]string{"B"})
	}

	data = MapData(&test)

	if count := len(data); count != 4 {
		t.Errorf("Data amount %d, want %d", count, 4)
	}

	if val, ok := data["f"]; !ok && val != "F" {
		t.Errorf("Expected key %s with value %s in data", "f", "F")
	}

	type Test4 struct {
		Timestamp time.Time `json:"timestamp,omitempty"`
	}

	test4 := Test4{Timestamp: time.Now()}
	data = MapData(test4)

	if count := len(data); count != 1 {
		t.Errorf("Data amount %d, want %d", count, 1)
	}

	type Test5 struct {
		value int
		Value int `json:"value"`
	}

	type Test6 struct {
		Test5 Test5 `json:"test5,omitempty"`
	}

	test5 := new(Test5)
	test5.Value = 2
	test6 := Test6{*test5}
	data = MapData(test6)

	if count := len(data); count != 1 {
		t.Errorf("Data amount %d, want %d", count, 1)
	}

	if _, ok := data["test5"]; !ok {
		t.Errorf("Expected key %s in data", "test5")
	}

	type YYY struct {
		A string `json:"a"`
	}

	type Test8 struct {
		Y []YYY         `json:"yyyy"`
		Z []*YYY        `json:"zzzz"`
		A []*CustomType `json:"aaaa"`
		C []CustomType  `json:"cccc"`
		B *CustomType   `json:"bbbb"`
		X []string      `json:"string"`
	}

	testY1 := YYY{A: "A value"}
	a := CustomType("test1")
	test8 := Test8{X: []string{"x1", "x2"}, Y: []YYY{testY1}, A: []*CustomType{&a}, B: &a, Z: []*YYY{&testY1}, C: []CustomType{a}}
	data = MapData(test8)

	if count := len(data); count != 6 {
		t.Errorf("Data amount %d, want %d", count, 6)
	}

	val, ok := data["yyyy"]

	if !ok {
		t.Errorf("expected key %s", "yyyy")
	}

	if _, ok := val.([]PropertyMap); !ok {
		t.Errorf("type is %s, want %T", reflect.TypeOf(val), []PropertyMap{})
	}

	val, ok = data["zzzz"]

	if !ok {
		t.Errorf("expected key %s", "zzzz")
	}

	if _, ok := val.([]PropertyMap); !ok {
		t.Errorf("type is %s, want %T", reflect.TypeOf(val), []PropertyMap{})
	}

	val, ok = data["aaaa"]

	if !ok {
		t.Errorf("expected key %s", "aaaa")
	}

	if _, ok := val.([]string); !ok {
		t.Errorf("type is %s, want %T", reflect.TypeOf(val), []string{})
	}

	val, ok = data["cccc"]

	if !ok {
		t.Errorf("expected key %s", "cccc")
	}

	if _, ok := val.([]string); !ok {
		t.Errorf("type is %s, want %T", reflect.TypeOf(val), []string{})
	}

	val, ok = data["string"]

	if !ok {
		t.Errorf("expected key %s", "string")
	}

	if _, ok := val.([]string); !ok {
		t.Errorf("type is %s, want %T", reflect.TypeOf(val), []string{})
	}

	type Test9 struct {
		A []CustomType3 `json:"a"`
		B []CustomType4 `json:"b"`
	}

	test9 := Test9{A: []CustomType3{CustomType3("test")}, B: []CustomType4{CustomType4("")}}
	data = MapData(test9)

	if count := len(data); count != 2 {
		t.Errorf("Data amount %d, want %d", count, 2)
	}

	val, ok = data["a"]

	if !ok {
		t.Errorf("expected key %s", "a")
	}

	if _, ok := val.([]string); !ok {
		t.Errorf("type is %s, want %T", reflect.TypeOf(val), []string{})
	}

	if count := len(val.([]string)); count != 0 {
		t.Errorf("Data amount %d, want %d", count, 0)
	}

	val, ok = data["b"]

	if !ok {
		t.Errorf("expected key %s", "b")
	}

	if _, ok := val.([]string); !ok {
		t.Errorf("type is %s, want %T", reflect.TypeOf(val), []string{})
	}

	if count := len(val.([]string)); count != 0 {
		t.Errorf("Data amount %d, want %d", count, 0)
	}

	type Test10 struct {
		A []CustomType4 `json:"a,omitempty"`
	}

	test10 := Test10{A: []CustomType4{CustomType4("")}}
	data = MapData(test10)

	if count := len(data); count != 0 {
		t.Errorf("Data amount %d, want %d", count, 0)
	}

	test10 = Test10{A: []CustomType4{CustomType4("test")}}
	data = MapData(test10)

	if count := len(data); count != 1 {
		t.Errorf("Data amount %d, want %d", count, 1)
	}

	val, ok = data["a"]

	if !ok {
		t.Errorf("expected key %s", "a")
	}

	if _, ok := val.([]string); !ok {
		t.Errorf("type is %s, want %T", reflect.TypeOf(val), []PropertyMap{})
	}

	if count := len(val.([]string)); count != 1 {
		t.Errorf("Data amount %d, want %d", count, 1)
	}

	type Test11 struct {
		A []string `json:"a"`
	}

	test11 := Test11{A: []string{""}}
	data = MapData(test11)

	if count := len(data); count != 1 {
		t.Errorf("Data amount %d, want %d", count, 1)
	}

	type Test12 struct {
		A []string `json:"a,omitempty"`
	}

	test12 := Test12{A: []string{""}}
	data = MapData(test12)

	if count := len(data); count != 0 {
		t.Errorf("Data amount %d, want %d", count, 0)
	}

	test12.A = append(test12.A, "test value")
	data = MapData(test12)

	if count := len(data); count != 1 {
		t.Errorf("Data amount %d, want %d", count, 1)
	}
}

func TestMapDataWithPointers(t *testing.T) {
	type Test4 struct {
		E string `json:"e,omitempty"`
	}

	type Test3 struct {
		D string `json:"d,omitempty"`
	}

	type Test2 struct {
		C string `json:"c"`
	}

	type Test1 struct {
		*Test2
		A *Test3 `json:"a"`
		B *Test4 `json:"b,omitempty"`
	}

	test1 := &Test1{}
	data := MapData(test1)

	if count := len(data); count != 1 {
		t.Errorf("Data amount %d, want %d", count, 1)
	}

	if _, ok := data["a"]; !ok {
		t.Errorf("Expected key %s in data", "a")
	}

	test2 := new(Test2)
	test2.C = "C"
	test1 = &Test1{Test2: test2}
	data = MapData(test1)

	if count := len(data); count != 2 {
		t.Errorf("Data amount %d, want %d", count, 2)
	}

	if _, ok := data["a"]; !ok {
		t.Errorf("Expected key %s in data", "a")
	}

	if _, ok := data["c"]; !ok {
		t.Errorf("Expected key %s in data", "c")
	}

	test3 := new(Test3)
	test4 := new(Test4)

	test1 = &Test1{Test2: test2, A: test3, B: test4}
	data = MapData(test1)

	if count := len(data); count != 2 {
		t.Errorf("Data amount %d, want %d", count, 2)
	}

	if _, ok := data["a"]; !ok {
		t.Errorf("Expected key %s in data", "a")
	}

	if _, isPropertyMap := data["a"].(PropertyMap); !isPropertyMap {
		t.Errorf("Expected value of key a in data: %v", PropertyMap{})
	}

	if _, ok := data["c"]; !ok {
		t.Errorf("Expected key %s in data", "c")
	}

	test4 = new(Test4)
	test4.E = "E"

	test1 = &Test1{Test2: test2, A: test3, B: test4}
	data = MapData(test1)

	if count := len(data); count != 3 {
		t.Errorf("Data amount %d, want %d", count, 3)
	}

	if _, ok := data["a"]; !ok {
		t.Errorf("Expected key %s in data", "a")
	}

	if _, ok := data["b"]; !ok {
		t.Errorf("Expected key %s in data", "b")
	}

	if val, isPropertyMap := data["b"].(PropertyMap); !isPropertyMap || len(val) != 1 {
		t.Errorf("Expected value of key b in data: %v", PropertyMap{"e": "E"})
	}

	if _, ok := data["c"]; !ok {
		t.Errorf("Expected key %s in data", "c")
	}
}

type CustomType string

func (c *CustomType) MarshalJSON() ([]byte, error) {
	return []byte("test value"), nil
}

type CustomType2 string

func (c CustomType2) MarshalJSON() ([]byte, error) {
	return []byte("test value"), nil
}

type CustomType3 string

func (c CustomType3) MarshalJSON() ([]byte, error) {
	return []byte("test value"), errors.New("error")
}

type CustomType4 string

func (c *CustomType4) MarshalJSON() ([]byte, error) {
	return []byte(string(*c)), nil
}

func TestMapDataWithMaps(t *testing.T) {
	type Test1 struct {
		A map[string]string `json:"a,omitempty"`
	}

	test := &Test1{}
	data := MapData(test)

	if count := len(data); count != 0 {
		t.Errorf("Data amount %d, want %d", count, 0)
	}

	mapData := make(map[string]string)
	test.A = mapData
	data = MapData(test)

	if count := len(data); count != 0 {
		t.Errorf("Data amount %d, want %d", count, 0)
	}

	test.A["test"] = "value 1"
	data = MapData(test)

	if count := len(data); count != 1 {
		t.Errorf("Data amount %d, want %d", count, 1)
	}

	type Test2 struct {
		A map[string]CustomType4 `json:"a,omitempty"`
	}

	test2 := &Test2{}
	mapData2 := make(map[string]CustomType4)
	test2.A = mapData2
	data = MapData(test2)

	if count := len(data); count != 0 {
		t.Errorf("Data amount %d, want %d", count, 0)
	}

	test2.A["test"] = CustomType4("")
	data = MapData(test2)

	if count := len(data); count != 0 {
		t.Errorf("Data amount %d, want %d", count, 0)
	}

	test2.A["test"] = CustomType4("value 1")
	data = MapData(test2)

	if count := len(data); count != 1 {
		t.Errorf("Data amount %d, want %d", count, 1)
	}
}

func TestMapDataWithFuncs(t *testing.T) {
	type Test1 struct {
		A func() `json:"a,omitempty"`
	}

	test := &Test1{}
	data := MapData(test)

	if count := len(data); count != 0 {
		t.Errorf("Data amount %d, want %d", count, 0)
	}

	testFunc := func() {}
	test = &Test1{A: testFunc}
	data = MapData(test)

	if count := len(data); count != 1 {
		t.Errorf("Data amount %d, want %d", count, 1)
	}
}

func TestMapDataWithCustomMarshaler(t *testing.T) {
	type Test1 struct {
		A *CustomType `json:"a"`
		B CustomType  `json:"b"`
	}

	a := CustomType("test1")
	b := CustomType("test2")
	test := &Test1{A: &a, B: b}
	data := MapData(test)

	if count := len(data); count != 2 {
		t.Errorf("Data amount %d, want %d", count, 2)
	}

	if _, ok := data["a"]; !ok {
		t.Errorf("Expected key %s in data", "a")
	}

	if v := data["a"]; v != "test value" {
		t.Errorf("Value is %s, want %s", v, "test value")
	}

	if _, ok := data["b"]; !ok {
		t.Errorf("Expected key %s in data", "b")
	}

	if v := data["b"]; v != CustomType("test2") {
		t.Errorf("Value is %s, want %s", v, "test2")
	}

	type Test2 struct {
		A *CustomType2 `json:"a"`
		B CustomType2  `json:"b"`
	}

	a2 := CustomType2("test1")
	b2 := CustomType2("test2")
	test2 := &Test2{A: &a2, B: b2}
	data = MapData(test2)

	if count := len(data); count != 2 {
		t.Errorf("Data amount %d, want %d", count, 2)
	}

	if _, ok := data["a"]; !ok {
		t.Errorf("Expected key %s in data", "a")
	}

	if v := data["a"]; v != "test value" {
		t.Errorf("Value is %s, want %s", v, "test value")
	}

	if _, ok := data["b"]; !ok {
		t.Errorf("Expected key %s in data", "b")
	}

	if v := data["b"]; v != "test value" {
		t.Errorf("Value is %s, want %s", v, "test value")
	}

	type Test3 struct {
		A *CustomType3 `json:"a"`
	}

	a3 := CustomType3("test1")
	test3 := &Test3{A: &a3}
	data = MapData(test3)

	if count := len(data); count != 0 {
		t.Errorf("Data amount %d, want %d", count, 0)
	}

	type Test4 struct {
		A *CustomType4 `json:"a,omitempty"`
	}

	a4 := CustomType4("")
	test4 := &Test4{A: &a4}
	data = MapData(test4)

	if count := len(data); count != 0 {
		t.Errorf("Data amount %d, want %d", count, 0)
	}
}
