// go2hal v0.6.0
// Copyright (c) 2021 Patrick Moule
// License: MIT

package halforms

import (
	"net/http"

	"github.com/pmoule/go2hal/hal"
	"github.com/pmoule/go2hal/hal/mapping"
)

// InlineItem is a default item for the Options inline property.
//
// Properties:
//
// Prompt: display values of option.
//
// Value: value of option.
type InlineItem struct {
	Prompt string `json:"prompt"`
	Value  string `json:"value"`
}

// Options is a list of possible values for a Property.
//
// Properties:
//
// MaxItems: indicates the maximum number of items to return in the SelectedValues.
//
// MinItems: indicates the minimum number of items to return in the SelectedValues.
//
// Inline: contains the list of possible values. This can be an array od strings,
// an array of InlineItem elements or any custom type.
//
// Link: contains a href to an external resource contaning the list of possible values.
//
// PromptField: name of inline or link elements to use as prompt.
//
// SelectedValues: contains the list of preselected values from possible values.
//
// ValueField name of inline or link elements to use as value.
type Options struct {
	Inline         []interface{}   `json:"inline"`
	Link           *hal.LinkObject `json:"link,omitempty"`
	MaxItems       string          `json:"maxItems"`
	MinItems       string          `json:"minItems"`
	PromptField    string          `json:"promptField"`
	SelectedValues []string        `json:"selectedValues"`
	ValueField     string          `json:"valueField"`
}

// Property decribes details of a state transition element.
//
// Properties:
//
// Name: the Property name.
//
// Prompt: human readable name to display.
//
// ReadOnly: indicates wheter it' s a read-only property.
//
// Required: indicates wheter it' s a required property.
//
// Regex: a regular expression string to be applied to the value.
//
// Templated: indicates whether value property conatins URI template to be resolved by client.
//
// Value: The property value. May be a templated URI.
//
// Cols: the maximum number of characters per line in multiline types.
//
// Max: maximum value for numeric values.
//
// MaxLength: maximum number of characters for string values.
//
// Min: minimum value for numeric values.
//
// MinLength: minimum number of characters for string values.
//
// Options: contains information about sets of possible values.
//
// Placeholder: hint to describe form value.
//
// Rows: the maximum number of rows in multiline types.
//
// Step: interval between numeric values.
//
// Type: the type to use for rendering the value.
type Property struct {
	Name        string   `json:"name"`
	Prompt      string   `json:"prompt"`
	ReadOnly    bool     `json:"readOnly"`
	Regex       string   `json:"regex"`
	Required    bool     `json:"required"`
	Templated   bool     `json:"templated"`
	Value       string   `json:"value"`
	Cols        uint     `json:"cols"`
	Max         string   `json:"max"`
	MaxLength   string   `json:"maxLength"`
	Min         string   `json:"min"`
	MinLength   string   `json:"minLength"`
	Options     *Options `json:"options,omitempty"`
	Placeholder string   `json:"placeholder"`
	Rows        uint     `json:"rows"`
	Step        int      `json:"step"`
	Type        string   `json:"type"`
}

// NewProperty returns an initialised Property with a provided name.
//
// Default values are:
//
// Prompt: same as provided name
//
// ReadOnly: false
//
// Required: false
//
// Templated: false
//
// Cols: 40
//
// Rows: 5
//
// Type: "text"
func NewProperty(name string) *Property {
	property := &Property{Name: name, Prompt: name, Cols: 40, Rows: 5, Type: PropertyDefaultType}

	return property
}

// Template describes the state transition details including the HTTP method, message
// content-type, and arguments for the transition.
//
// Properties:
//
// ContentType: media type to use when sending a request body to the server.
//
// Key: a unique identifier for the template.
//
// Method: HTTP method to use for requests.
//
// Properties: contains a href to an external resource contaning the list of possible values.
//
// Target: target URL for submitting a HAL-FORMS values.
//
// Title: a human readable title for the template.
type Template struct {
	ContentType string      `json:"contentType"`
	Key         string      `json:"key"`
	Method      string      `json:"method"`
	Properties  []*Property `json:"properties"`
	Target      string      `json:"target"`
	Title       string      `json:"title"`
}

// NewTemplate returns an initialised Template.
//
// Default values are:
//
// ContentType: "application/json"
//
// Key: "default"
//
// Method: "GET"
//
// Title: "default"
func NewTemplate() *Template {
	template := &Template{ContentType: jsonMediaTypeIdentifier, Key: TemplateDefaultKey, Method: http.MethodGet, Properties: []*Property{}, Title: TemplateDefaultKey}

	return template
}

type templates map[string]*Template

func (t templates) ToMap() mapping.NamedMap {
	properties := mapping.PropertyMap{}

	for _, template := range t {
		properties[template.Key] = template
	}

	return mapping.NamedMap{Name: TemplatesProperty, Content: properties}
}

// Document contains links and state transition details.
type Document struct {
	links     hal.Links
	templates templates
}

// Links returns a "_links" named map of link relations and assigned links.
func (d *Document) Links() mapping.NamedMap {
	return d.links.ToMap()
}

// AddLinke adds a link relation to HAL-FORMS document.
func (d *Document) AddLink(rel hal.LinkRelation) {
	d.links[rel.Name()] = rel
}

// Templates returns a "_templates" named map of templates.
func (d *Document) Templates() mapping.NamedMap {
	return d.templates.ToMap()
}

// AddTemplate adds a template to HAL-FORMS document.
func (d *Document) AddTemplate(template *Template) {
	d.templates[template.Key] = template
}

// ToMap converts Document to a "root" named map for JSON conversion.
func (d *Document) ToMap() mapping.NamedMap {
	properties := mapping.PropertyMap{}
	namedMaps := []mapping.NamedMap{}
	namedMaps = append(namedMaps, d.Links())
	namedMaps = append(namedMaps, d.Templates())

	for _, namedMap := range namedMaps {
		properties[namedMap.Name] = namedMap.Content
	}

	return mapping.NamedMap{Name: "root", Content: properties}
}

// NewDocument returns an initialised Document with "self" link relation to provided href.
func NewDocument(href string) Document {
	document := Document{links: map[string]hal.LinkRelation{}, templates: map[string]*Template{}}
	link := &hal.LinkObject{Href: href, Type: MediaTypeIdentifier}
	self := hal.NewSelfLinkRelation()
	self.SetLink(link)
	document.AddLink(self)

	return document
}
