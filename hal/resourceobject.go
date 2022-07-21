// go2hal v0.6.0
// Copyright (c) 2021 Patrick Moule
// License: MIT

package hal

import (
	"github.com/skhaz/go2hal/hal/mapping"
	"github.com/skhaz/go2hal/hal/relationtype"
)

// Resource is the root element of a HAL document.
//
// A Resource can
//
// - have links - AddLink(LinkRelation)
//
// - have CURIEs - AddCurieLinks([]*LinkObject)
//
// - embed other resources - AddResource(ResourceRelation)
type Resource interface {
	Data() mapping.PropertyMap
	Links() mapping.NamedMap
	EmbeddedResources() mapping.NamedMap
	AddData(interface{})
	AddLink(LinkRelation)
	AddResource(ResourceRelation)
	AddCurieLinks([]*LinkObject)
	ToMap() mapping.NamedMap
}

type resourceObject struct {
	data     mapping.PropertyMap
	links    Links
	embedded embeddedResources
}

// NewResourceObject initialises a Resource.
func NewResourceObject() Resource {
	return &resourceObject{data: mapping.PropertyMap{}, links: Links{}, embedded: embeddedResources{}}
}

// Data returns mapping.PropertyMap describng the assigned data.
func (r *resourceObject) Data() mapping.PropertyMap {
	return r.data
}

// AddData assigns any type of data to ResourceObject.
func (r *resourceObject) AddData(data interface{}) {
	value := mapping.MapData(data)

	for k, v := range value {
		r.data[k] = v
	}
}

// Links returns a mapping.NamedMap of assigned link relations.
func (r *resourceObject) Links() mapping.NamedMap {
	return r.links.ToMap()
}

// EmbeddedResources returns a mapping.NamedMap of embedded resources.
func (r *resourceObject) EmbeddedResources() mapping.NamedMap {
	return r.embedded.ToMap()
}

// ToMap converts ResourceObject to mapping.NamedMap.
func (r *resourceObject) ToMap() mapping.NamedMap {
	properties := mapping.PropertyMap{}

	namedMaps := []mapping.NamedMap{}
	namedMaps = append(namedMaps, r.Links())
	namedMaps = append(namedMaps, r.EmbeddedResources())

	for _, namedMap := range namedMaps {
		if len(namedMap.Content) > 0 {
			properties[namedMap.Name] = namedMap.Content
		}
	}

	for key, val := range r.data {
		properties[key] = val
	}

	return mapping.NamedMap{Name: "root", Content: properties}
}

// AddCurieLinks adds a set of LinkObjects usable as CURIES to ResourceObject.
func (r *resourceObject) AddCurieLinks(linkObjects []*LinkObject) {
	rel, _ := NewLinkRelation(relationtype.CURIES)
	rel.SetLinks(linkObjects)
	r.AddLink(rel)
}

// AddLink adds a LinRelation to ResourceObject.
func (r *resourceObject) AddLink(rel LinkRelation) {
	r.links[rel.Name()] = rel
}

// AddResource adds a ResourceRelation to ResourceObject.
func (r *resourceObject) AddResource(rel ResourceRelation) {
	r.embedded[rel.Name()] = rel
}
