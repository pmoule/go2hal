// go2hal v0.2.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import (
	"testing"
)

func TestNewResourceFactory(t *testing.T) {
	factory := NewResourceFactory(nil)

	if factory == nil {
		t.Errorf("Factory is not initialised.")
	}

	curieLinks := make([]*LinkObject, 1)
	curieLinks[0] = &LinkObject{Name: "Curie1"}
	factory = NewResourceFactory(curieLinks)

	resourceLink := factory.CreateResourceLink("href", "Curie1")

	if resourceLink.CurieLink().Name != curieLinks[0].Name {
		t.Errorf("Factory is not initialised with CURIE link.")
	}
}
