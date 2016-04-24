// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import (
	"errors"
	"github.com/pmoule/go2hal/hal/relationtype"
)

// Relation provides a descriptive name to add a meaning to
// links. To create a more discoverable API, a Relation
// can optionally be prefixed with a CURIE name.
type Relation interface {
	Name() string
	FullName() string
	SetCurieLink(curieLink *LinkObject)
	CurieLink() LinkObject
}

// LinkRelation is a Relation for Link Object assignment.
type LinkRelation interface {
	Relation
	SetLink(*LinkObject)
	SetLinks([]*LinkObject)
	IsLinkSet() bool
	Links() []*LinkObject
}

// ResourceRelation is a Relation for Resource Object assignment.
type ResourceRelation interface {
	Relation
	SetResource(Resource)
	SetResources([]Resource)
	IsResourceSet() bool
	Resources() []Resource
}

// An unexported implementation of the LinkRelation interface.
type linkRelation struct {
	name       string
	curieLink  *LinkObject
	isValueSet bool
	links      []*LinkObject
	resources  []Resource
}

// newRelation initializes a valid link relation.
func newRelation(name string) (*linkRelation, error) {
	if name == "" {
		return nil, errors.New("LinkRelation requires a name value.")
	}

	return &linkRelation{name: name, links: []*LinkObject{}, resources: []Resource{}}, nil
}

// NewLinkRelation initializes a valid Link Relation for Link Object assignment.
func NewLinkRelation(name string) (LinkRelation, error) {
	return newRelation(name)
}

// NewSelfLinRelation initializes a valid Link Relation used for targeting
// the URI of the resource it is attached to.
// See http://www.iana.org/assignments/link-relations/link-relations.xhtml.
func NewSelfLinkRelation() LinkRelation {
	relation, _ := NewLinkRelation(relationtype.Self)
	return relation
}

// NewResourceRelation initializes a valid link relation for Resource Object assignment.
func NewResourceRelation(name string) (ResourceRelation, error) {
	return newRelation(name)
}

// Returns the assigned name.
func (lr *linkRelation) Name() string {
	return lr.name
}

// Returns the assigned name. In case of preceding CURIE link assignment
// the returned name is prefixed with the CURIE's name.
func (lr *linkRelation) FullName() string {
	if lr.curieLink == nil {
		return lr.Name()
	}

	return lr.curieLink.Name + ":" + lr.Name()
}

// Use CURIES to create a more discoverable API by assigning
// a CURIE link.
func (lr *linkRelation) SetCurieLink(curieLink *LinkObject) {
	lr.curieLink = curieLink
}

// Returns the assigned CURIE link.
func (lr *linkRelation) CurieLink() LinkObject {
	return *lr.curieLink
}

// Assign a single Link Object
func (lr *linkRelation) SetLink(link *LinkObject) {
	lr.links = append(lr.links, link)
	lr.isValueSet = false
}

// Assign a slice of Link Objects
func (lr *linkRelation) SetLinks(links []*LinkObject) {
	for _, link := range links {
		lr.links = append(lr.links, link)
	}

	lr.isValueSet = true
}

func (lr *linkRelation) Links() []*LinkObject {
	return lr.links
}

func (lr *linkRelation) IsLinkSet() bool {
	return lr.isValueSet
}

// Assign a Resource Object
func (lr *linkRelation) SetResource(resource Resource) {
	lr.resources = append(lr.resources, resource)
	lr.isValueSet = false
}

// Assign a slice of Resource Objects
func (lr *linkRelation) SetResources(resources []Resource) {
	for _, resource := range resources {
		lr.resources = append(lr.resources, resource)
	}

	lr.isValueSet = true
}

func (lr *linkRelation) Resources() []Resource {
	return lr.resources
}

func (lr *linkRelation) IsResourceSet() bool {
	return lr.isValueSet
}

type links map[string]LinkRelation

func (l links) ToMap() NamedMap {
	linkMap := PropertyMap{}

	for _, val := range l {
		if val.IsLinkSet() {
			linkMap[val.FullName()] = val.Links()
		} else {
			if len(val.Links()) > 0 {
				linkMap[val.FullName()] = val.Links()[0]
			} else {
				linkMap[val.FullName()] = nil
			}
		}
	}

	return NamedMap{Name: LinksProperty, Content: linkMap}
}

type embeddedResources map[string]ResourceRelation

func (er embeddedResources) ToMap() NamedMap {
	embeddedMap := PropertyMap{}

	for _, val := range er {
		resources := val.Resources()

		var propertyMaps []PropertyMap

		for _, resource := range resources {
			if mapper, ok := resource.(mapper); ok {
				namedMap := mapper.ToMap()
				propertyMaps = append(propertyMaps, namedMap.Content)
			}
		}

		if val.IsResourceSet() {
			embeddedMap[val.FullName()] = propertyMaps
		} else {
			if len(propertyMaps) > 0 {
				embeddedMap[val.FullName()] = propertyMaps[0]
			} else {
				embeddedMap[val.FullName()] = nil
			}
		}
	}

	return NamedMap{Name: EmbeddedProperty, Content:  embeddedMap}
}