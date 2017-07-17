// go2hal v0.3.0
// Copyright (c) 2017 Patrick Moule
// License: MIT

package hal

import (
	"testing"

	"github.com/pmoule/go2hal/hal/relationtype"
)

func TestNewResourceFactory(t *testing.T) {
	factory := NewResourceFactory(nil)

	if factory == nil {
		t.Errorf("Factory is nil.")
	}

	curieLinks := make([]*LinkObject, 1)
	curieLinks[0] = &LinkObject{Name: "Curie1"}
	factory = NewResourceFactory(curieLinks)

	resourceLink := factory.CreateResourceLink("href", "Curie1")

	if resourceLink.CurieLink().Name != curieLinks[0].Name {
		t.Errorf("Factory is not initialised with CURIE link.")
	}
}

func TestCreateRootResource(t *testing.T) {
	curieLinks := make([]*LinkObject, 1)
	curieLinks[0] = &LinkObject{Name: "Curie1"}
	factory := NewResourceFactory(curieLinks)

	root := factory.CreateRootResource("")

	if root == nil {
		t.Errorf("Root Resource is nil.")
	}

	links := root.Links()

	if links.Name != LinksProperty {
		t.Errorf("Links relation name %s, wanted %s", links.Name, LinksProperty)
	}

	if len(links.Content) != len(curieLinks) {
		t.Errorf("Links count %d, wanted %d", len(links.Content), len(curieLinks))
	}

	if links.Content[relationtype.CURIES] == nil {
		t.Errorf("Missing link relation %s", relationtype.CURIES)
	}

	if links.Content[relationtype.Self] != nil {
		t.Errorf("Not expected link relation %s", relationtype.Self)
	}

	selfLink := "http://self"
	root = factory.CreateRootResource(selfLink)
	links = root.Links()

	if links.Content[relationtype.Self] == nil {
		t.Errorf("Missing link relation %s", relationtype.Self)
	}

	if links.Content[relationtype.Self].(*LinkObject).Href != selfLink {
		t.Errorf("Self link is %s, wanted %s", links.Content[relationtype.Self], selfLink)
	}
}

func TestCreateEmbeddedResource(t *testing.T) {
	curieLinks := make([]*LinkObject, 1)
	curieLinks[0] = &LinkObject{Name: "Curie1"}
	factory := NewResourceFactory(curieLinks)

	resource := factory.CreateEmbeddedResource("")

	if resource == nil {
		t.Errorf("Embedded resource is nil.")
	}

	links := resource.Links()

	if links.Name != LinksProperty {
		t.Errorf("Links relation name %s, wanted %s", links.Name, LinksProperty)
	}

	if len(links.Content) > 0 {
		t.Errorf("Links count %d, wanted %d", len(links.Content), 0)
	}

	selfLink := "http://self"
	resource = factory.CreateEmbeddedResource(selfLink)
	links = resource.Links()

	if links.Content[relationtype.Self] == nil {
		t.Errorf("Missing link relation %s", relationtype.Self)
	}

	if links.Content[relationtype.Self].(*LinkObject).Href != selfLink {
		t.Errorf("Self link is %s, wanted %s", links.Content[relationtype.Self], selfLink)
	}
}

func TestCreateLink(t *testing.T) {
	curieLinks := make([]*LinkObject, 1)
	curieLinks[0] = &LinkObject{Name: "Curie1"}
	factory := NewResourceFactory(curieLinks)

	relationName := ""
	href := ""
	curieLinkName := ""
	link := factory.CreateLink(relationName, href, curieLinkName)

	if link != nil {
		t.Errorf("Link should be nil: %v", link)
	}

	href = "href"
	link = factory.CreateLink(relationName, href, curieLinkName)

	if link != nil {
		t.Errorf("Link should be nil: %v", link)
	}

	relationName = "relationName"
	link = factory.CreateLink(relationName, href, curieLinkName)

	if link == nil {
		t.Errorf("Link is nil")
	}

	if len(link.Links()) != 1 {
		t.Errorf("Link relation with %d links, wanted %d links", len(link.Links()), 1)
	}

	if link.FullName() != relationName {
		t.Errorf("Full name is %s, wanted %s", link.FullName(), relationName)
	}

	if link.Links()[0].Href != href {
		t.Errorf("Href is %s, wanted %s", link.Links()[0].Href, href)
	}

	curieLinkName = curieLinks[0].Name
	link = factory.CreateLink(relationName, href, curieLinkName)

	if link.FullName() != curieLinkName+":"+relationName {
		t.Errorf("Full name is %s, wanted %s", link.FullName(), curieLinkName+":"+relationName)
	}
}

func TestCreateResourceLink(t *testing.T) {
	curieLinks := make([]*LinkObject, 1)
	curieLinks[0] = &LinkObject{Name: "Curie1"}
	factory := NewResourceFactory(curieLinks)

	relationName := ""
	curieLinkName := ""
	link := factory.CreateResourceLink(relationName, curieLinkName)

	if link != nil {
		t.Errorf("Link should be nil: %v", link)
	}

	relationName = "relationName"
	link = factory.CreateResourceLink(relationName, curieLinkName)

	if link == nil {
		t.Errorf("Link is nil")
	}

	curieLinkName = curieLinks[0].Name
	link = factory.CreateResourceLink(relationName, curieLinkName)

	if link.FullName() != curieLinkName+":"+relationName {
		t.Errorf("Full name is %s, wanted %s", link.FullName(), curieLinkName+":"+relationName)
	}
}
