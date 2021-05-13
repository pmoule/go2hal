// go2hal v0.4.0
// Copyright (c) 2020 Patrick Moule
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

	linkObject, _ := NewLinkObject("test")
	relation.SetLink(linkObject)

	if count := len(relation.Links()); count != 1 {
		t.Errorf("LinkRelation links == %d, want %d", count, 1)
	}

	if relation.IsLinkSet() {
		t.Errorf("ResourceRelation not should be a link set")
	}

	relation.SetLinks([]*LinkObject{linkObject})

	if count := len(relation.Links()); count != 2 {
		t.Errorf("LinkRelation links == %d, want %d", count, 2)
	}

	if !relation.IsLinkSet() {
		t.Errorf("ResourceRelation should be a link set")
	}
}

func TestNewResourceRelation(t *testing.T) {
	wantedRelationName := "relation"
	relation, _ := NewResourceRelation(wantedRelationName)

	if relation.Name() != wantedRelationName {
		t.Errorf("ResourceRelation name == %q, want %q", relation.Name(), wantedRelationName)
	}

	resource := NewResourceObject()
	relation.SetResource(resource)

	if relation.IsResourceSet() {
		t.Errorf("ResourceRelation should not be a resource set")
	}

	relation.SetResources([]Resource{resource})

	if !relation.IsResourceSet() {
		t.Errorf("ResourceRelation should be a resource set")
	}
}

func TestLinksToMap(t *testing.T) {
	wantedRelationName := "relation"
	relation, _ := NewLinkRelation(wantedRelationName)
	resourceObject := NewResourceObject()
	resourceObject.AddLink(relation)
	links := resourceObject.Links()

	if len(links.Content) != 1 {
		t.Errorf("Link relations count %d, wanted %d", len(links.Content), 1)
	}

	if links.Content[wantedRelationName] != nil {
		t.Errorf("Link relations is %v, wanted %v", links.Content[wantedRelationName], nil)
	}
}

func TestEmbeddedResourcesToMap(t *testing.T) {
	wantedRelationName := "relation"
	relation, _ := NewResourceRelation(wantedRelationName)
	resourceObject := NewResourceObject()
	resourceObject.AddResource(relation)
	er := resourceObject.EmbeddedResources()

	if len(er.Content) != 1 {
		t.Errorf("Link relations count %d, wanted %d", len(er.Content), 1)
	}

	if er.Content[wantedRelationName] != nil {
		t.Errorf("Link relations is %v, wanted %v", er.Content[wantedRelationName], nil)
	}
}
