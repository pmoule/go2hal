// go2hal v0.3.2
// Copyright (c) 2017 Patrick Moule
// License: MIT

package mapping

import (
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
		K string `json:"k, omitempty"`
	}

	type Test2 struct {
		F string `json:"f"`
		g string
		H Test3  `json:"h"`
		I [2]int `json:"i, omitempty"`
		J Test3  `json:"j, omitempty"`
	}

	type Test1 struct {
		Test7
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
		Timestamp time.Time `json:"timestamp, omitempty"`
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
		Test5 Test5 `json:"test5, omitempty"`
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
}
