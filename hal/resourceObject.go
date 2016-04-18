package hal

import (
	"encoding/json"
)

//todo: add reference to https://tools.ietf.org/html/draft-kelly-json-hal-07
type Resource interface {
	Data() PropertyMap
	Links() NamedMap
	EmbeddedResources() NamedMap
	AddLinkObject(rel *linkRelation, linkObject *LinkObject)
	AddLinkObjects(rel *linkRelation, linkObjects []*LinkObject)
	AddResourceObject(rel *linkRelation, resource Resource)
	AddResourceObjects(rel *linkRelation, resources []Resource)
	AddCurieLink(link *LinkObject)
	ToJson() ([]byte, error)
}

type resourceObject struct {
	data     PropertyMap `json:"-"`
	links    linkRelationMapper `json:"_links,omitempty"`
	embedded embeddedResourceMapper `json:"_embedded,omitempty"`
}

func NewResourceObject() Resource {
	return &resourceObject{data: PropertyMap{}}
}

func (r *resourceObject) Data() PropertyMap {
	return r.data
}

func (r *resourceObject) Links() NamedMap {
	return r.links.ToMap()
}

func (r *resourceObject) EmbeddedResources() NamedMap {
	return r.embedded.ToMap()
}

//todo: rename to ToJSON
func (r *resourceObject) ToJson() ([]byte, error) {
	resourceMap := r.ToMap()
	return json.Marshal(resourceMap.Content)
}

func (r *resourceObject) ToMap() NamedMap {
	resourceMap := make(map[string]interface{})

	mappers := []mapper{&r.links, &r.embedded}

	for _, mapper := range mappers {
		namedMap := mapper.ToMap()

		if len(namedMap.Content) > 0 {
			resourceMap[namedMap.Name] = namedMap.Content
		}
	}

	for key, val := range r.data {
		resourceMap[key] = val
	}

	return NamedMap{Name: "root", Content: resourceMap}
}

func (r *resourceObject) AddCurieLink(link *LinkObject) {
	var linkSlice []*LinkObject

	rel := linkRelation{name: "curies"}
	existingLinkSlice, ok := r.links[rel]

	if !ok {
		linkSlice = []*LinkObject{}
	} else {
		linkSlice = existingLinkSlice.([]*LinkObject)
	}

	linkSlice = append(linkSlice, link)

	r.AddLinkObjects(&rel, linkSlice)
}

func (r *resourceObject) AddLinkObject(rel *linkRelation, linkObject *LinkObject) {
	if r.links == nil {
		r.links = linkRelationMapper{}
	}

	r.links[*rel] = linkObject
}

func (r *resourceObject) AddLinkObjects(rel *linkRelation, linkObjects []*LinkObject) {
	if r.links == nil {
		r.links = linkRelationMapper{}
	}

	dataSlice := make([]*LinkObject, len(linkObjects))

	for index, linkObject := range linkObjects {
		dataSlice[index] = linkObject
	}

	r.links[*rel] = dataSlice
}

func (r *resourceObject) AddResourceObject(rel *linkRelation, resource Resource) {
	if r.embedded == nil {
		r.embedded = embeddedResourceMapper{}
	}

	r.embedded[*rel] = resource
}

func (r *resourceObject) AddResourceObjects(rel *linkRelation, resources []Resource) {
	if r.embedded == nil {
		r.embedded = embeddedResourceMapper{}
	}

	dataSlice := make([]Resource, len(resources))

	for index, resource := range resources {
		dataSlice[index] = resource
	}

	r.embedded[*rel] = dataSlice
}