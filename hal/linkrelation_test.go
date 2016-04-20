// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import (
	"testing"
)

func TestNewLinkRelation(t *testing.T) {
	wantedRelationName := "relation"
	relation, _ := NewLinkRelation(wantedRelationName)

	wantedName := "link"
	wantedHref := "http://{rel}"

	wantedCurieLink, _ := NewCurieLink(wantedName, wantedHref)
	wantedCuriedRelationName := wantedName + ":" + wantedRelationName

	if relation.Name() != wantedRelationName {
		t.Errorf("LinkRelation name == %q, want %q", relation.Name(), wantedRelationName)
	}

	if relation.FullName() != wantedRelationName {
		t.Errorf("LinkRelation name == %q, want %q", relation.FullName(), wantedRelationName)
	}

	relation.SetCurieLink(wantedCurieLink)
	result := relation.CurieLink()
	result.Name = "changed"

	if wantedCurieLink.Name == result.Name {
		t.Errorf("Identical curie link")
	}

	_, invalidNameError := NewLinkRelation("")

	if invalidNameError == nil {
		t.Errorf("NewCurieLink should return an error due to an invalid name value.")
	}

	if relation.Name() != wantedRelationName {
		t.Errorf("LinkRelation name == %q, want %q", relation.Name(), wantedRelationName)
	}

	if relation.FullName() != wantedCuriedRelationName {
		t.Errorf("LinkRelation name == %q, want %q", relation.FullName(), wantedCuriedRelationName)
	}
}