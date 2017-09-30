// go2hal v0.3.2
// Copyright (c) 2017 Patrick Moule
// License: MIT

package hal

import (
	"github.com/pmoule/go2hal/hal/mapping"
	"github.com/pmoule/go2hal/hal/relationtype"
)

// Resource is the root element of a HAL document.
// A Resource can
// - have links - AddLink(LinkRelation)
// - have CURIEs - AddCurieLinks([]*LinkObject)
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
	data     mapping.PropertyMap `json:"-"`
	links    links               `json:"_links,omitempty"`
	embedded embeddedResources   `json:"_embedded,omitempty"`
}

// NewResourceObject initialises a valid Resource.
func NewResourceObject() Resource {
	return &resourceObject{data: mapping.PropertyMap{}, links: links{}, embedded: embeddedResources{}}
}

func (r *resourceObject) Data() mapping.PropertyMap {
	return r.data
}

func (r *resourceObject) AddData(data interface{}) {
	value := mapping.MapData(data)

	for k, v := range value {
		r.data[k] = v
	}
}

func (r *resourceObject) Links() mapping.NamedMap {
	return r.links.ToMap()
}

func (r *resourceObject) EmbeddedResources() mapping.NamedMap {
	return r.embedded.ToMap()
}

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

func (r *resourceObject) AddCurieLinks(linkObjects []*LinkObject) {
	rel, _ := NewLinkRelation(relationtype.CURIES)
	rel.SetLinks(linkObjects)
	r.AddLink(rel)
}

func (r *resourceObject) AddLink(rel LinkRelation) {
	r.links[rel.Name()] = rel
}

func (r *resourceObject) AddResource(rel ResourceRelation) {
	r.embedded[rel.Name()] = rel
}
