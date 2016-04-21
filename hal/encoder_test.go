// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import "testing"

func TestEncoder(t *testing.T) {
	wanted := []interface{}{2, "{}"}

	encoder := new(Encoder)

	resource := NewResourceObject()
	bytes, _ := encoder.ToJSON(resource)

	if count := len(bytes); count != wanted[0] {
		t.Errorf("ToJSON retuns zero byte == %d, want %d", count, wanted[0])
	}

	if value := string(bytes[:]); value != wanted[1] {
		t.Errorf("LinkRelation name == %s, want %s", value, wanted[1])
	}
}