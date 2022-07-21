// go2hal v0.6.0
// Copyright (c) 2021 Patrick Moule
// License: MIT

package hal

import (
	"errors"

	"github.com/skhaz/go2hal/hal/mapping"
	"github.com/skhaz/go2hal/hal/relationtype"
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

// LinkRelation is a Relation for LinkObject assignment.
type LinkRelation interface {
	Relation
	SetLink(*LinkObject)
	SetLinks([]*LinkObject)
	IsLinkSet() bool
	Links() []*LinkObject
}

// ResourceRelation is a Relation for ResourceObject assignment.
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

// newRelation initializes a linkRelation.
func newRelation(name string) (*linkRelation, error) {
	if name == "" {
		return nil, errors.New("LinkRelation requires a name value")
	}

	return &linkRelation{name: name, links: []*LinkObject{}, resources: []Resource{}}, nil
}

// NewLinkRelation initializes a LinkRelation for Link Object assignment.
func NewLinkRelation(name string) (LinkRelation, error) {
	return newRelation(name)
}

// NewSelfLinkRelation initializes a LinkRelation used for targeting
// the URI of the resource it is attached to.
// See http://www.iana.org/assignments/link-relations/link-relations.xhtml.
func NewSelfLinkRelation() LinkRelation {
	relation, _ := NewLinkRelation(relationtype.Self)
	return relation
}

// NewResourceRelation initializes a link relation for Resource Object assignment.
func NewResourceRelation(name string) (ResourceRelation, error) {
	return newRelation(name)
}

// Name returns the assigned name.
func (lr *linkRelation) Name() string {
	return lr.name
}

// FullName returns the assigned name. In case of preceding CURIE link assignment
// the returned name is prefixed with the CURIE's name.
func (lr *linkRelation) FullName() string {
	if lr.curieLink == nil {
		return lr.Name()
	}

	return lr.curieLink.Name + ":" + lr.Name()
}

// SetCurieLink assigns a CURIE.
func (lr *linkRelation) SetCurieLink(curieLink *LinkObject) {
	lr.curieLink = curieLink
}

// CurieLink returns the assigned CURIE link.
func (lr *linkRelation) CurieLink() LinkObject {
	return *lr.curieLink
}

// SetLink assigns a single Link Object.
func (lr *linkRelation) SetLink(link *LinkObject) {
	lr.links = append(lr.links, link)
	lr.isValueSet = false
}

// SetLinks assigns a slice of Link Objects.
func (lr *linkRelation) SetLinks(links []*LinkObject) {
	lr.links = append(lr.links, links...)
	lr.isValueSet = true
}

// Links returns assigned links.
func (lr *linkRelation) Links() []*LinkObject {
	return lr.links
}

// IsLinkSet returns true if assigned links are always structured in an array.
func (lr *linkRelation) IsLinkSet() bool {
	return lr.isValueSet
}

// SetResource assigns a resource.
func (lr *linkRelation) SetResource(resource Resource) {
	lr.resources = append(lr.resources, resource)
	lr.isValueSet = false
}

// SetResources assigns a slice of resources.
func (lr *linkRelation) SetResources(resources []Resource) {
	lr.resources = append(lr.resources, resources...)
	lr.isValueSet = true
}

// Resources returns assigned resources.
func (lr *linkRelation) Resources() []Resource {
	return lr.resources
}

// IsResourceSet returns true if assigned resources are always structured in an array.
func (lr *linkRelation) IsResourceSet() bool {
	return lr.isValueSet
}

// Links
type Links map[string]LinkRelation

// ToMap converts Links to mapping.NamedMap.
func (l Links) ToMap() mapping.NamedMap {
	properties := mapping.PropertyMap{}

	for _, val := range l {
		if val.IsLinkSet() {
			properties[val.FullName()] = val.Links()
		} else {
			if len(val.Links()) > 0 {
				properties[val.FullName()] = val.Links()[0]
			} else {
				properties[val.FullName()] = nil
			}
		}
	}

	return mapping.NamedMap{Name: LinksProperty, Content: properties}
}

type embeddedResources map[string]ResourceRelation

// ToMap converts embeddedResources to mapping.NamedMap.
func (er embeddedResources) ToMap() mapping.NamedMap {
	embeddedProperties := mapping.PropertyMap{}

	for _, val := range er {
		resources := val.Resources()

		var properties = []mapping.PropertyMap{}

		for _, resource := range resources {
			resourceObject := resource.(*resourceObject)
			namedMap := resourceObject.ToMap()
			properties = append(properties, namedMap.Content)
		}

		if val.IsResourceSet() {
			embeddedProperties[val.FullName()] = properties
		} else {
			if len(properties) > 0 {
				embeddedProperties[val.FullName()] = properties[0]
			} else {
				embeddedProperties[val.FullName()] = nil
			}
		}
	}

	return mapping.NamedMap{Name: EmbeddedProperty, Content: embeddedProperties}
}
