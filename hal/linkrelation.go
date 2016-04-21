// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import "errors"

// Relation provides a descriptive name to add a meaning to
// links. To create a more discoverable API, a Relation
// can optionally be prefixed with a CURIE name.
type Relation interface {
	Name() string
	FullName() string
	SetCurieLink(curieLink *LinkObject)
	CurieLink() LinkObject
	Value() interface{}
}

// LinkRelation is a Relation for Link Object assignment.
type LinkRelation interface {
	Relation
	SetLink(*LinkObject)
	SetLinks([]*LinkObject)
}

// ResourceRelation is a Relation for Resource Object assignment.
type ResourceRelation interface {
	Relation
	SetResource(Resource)
	SetResources([]Resource)
}

// An unexported implementation of the LinkRelation interface.
type linkRelation struct {
	name      string
	curieLink *LinkObject
	value     interface{}
}

// newRelation initializes a valid link relation.
func newRelation(name string) (*linkRelation, error) {
	if name == "" {
		return nil, errors.New("LinkRelation requires a name value.")
	}

	return &linkRelation{name: name}, nil
}

// NewLinkRelation initializes a valid Link Relation for Link Object assignment.
func NewLinkRelation(name string) (LinkRelation, error) {
	return newRelation(name)
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
	lr.value = link
}

// Assign a slice of Link Objects
func (lr *linkRelation) SetLinks(links []*LinkObject) {
	lr.value = links
}

// Assign a Resource Object
func (lr *linkRelation) SetResource(resource Resource) {
	lr.value = resource
}

// Assign a slice of Resource Objects
func (lr *linkRelation) SetResources(resources []Resource) {
	lr.value = resources
}

func (lr *linkRelation) Value() interface{} {
	return lr.value
}

type links map[string]LinkRelation

func (l *links) ToMap() NamedMap {
	linkMap := PropertyMap{}

	for _, val := range *l {
		linkMap[val.FullName()] = val.Value()
	}

	return NamedMap{Name: "_links", Content: linkMap}
}

type embeddedResources map[string]ResourceRelation

func (er *embeddedResources) ToMap() NamedMap {
	embeddedMap := PropertyMap{}

	for _, val := range *er {
		embeddedMap[val.FullName()] = val.Value()
	}

	for key, val := range embeddedMap {
		if resourceArray, isSlice := val.([]Resource); isSlice {
			var resources []interface{}

			for _, resource := range resourceArray {
				if mapper, ok := resource.(mapper); ok {
					namedMap := mapper.ToMap()
					resources = append(resources, namedMap.Content)
				}
			}

			embeddedMap[key] = resources

		} else {
			resource := val.(Resource)

			if mapper, ok := resource.(mapper); ok {
				namedMap := mapper.ToMap()
				embeddedMap[key] = namedMap.Content
			}
		}
	}

	return NamedMap{Name: "_embedded", Content: embeddedMap }
}