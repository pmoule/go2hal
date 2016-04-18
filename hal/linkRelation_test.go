package hal

import (
	"testing"
)

func TestNewLinkRelation(t *testing.T) {
	relation, _ := NewLinkRelation("relation")

	wantedName := "link"
	wantedHref := "http://{rel}"

	wantedCurieLink, _ := NewCurieLink(wantedName, wantedHref)

	if wantedCurieLink.Name != wantedName {
		t.Errorf("LinkRelation name == %q, want %q", wantedCurieLink.Name, wantedName)
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
}