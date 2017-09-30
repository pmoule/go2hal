// go2hal v0.3.3
// Copyright (c) 2017 Patrick Moule
// License: MIT

package hal

import "errors"

// LinkObject is a hyperlink from the Resource it is attached to.
// A valid LinkObject requires a href value. All other properties are optional.
// See https://tools.ietf.org/html/draft-kelly-json-hal for
// property description.
type LinkObject struct {
	Href        string `json:"href,omitempty"`        //required
	Templated   bool   `json:"templated,omitempty"`   //optional
	Type        string `json:"type,omitempty"`        //optional
	Deprecation string `json:"deprecation,omitempty"` //optional
	Name        string `json:"name,omitempty"`        //optional
	Profile     string `json:"profile,omitempty"`     //optional
	Title       string `json:"title,omitempty"`       //optional
	HrefLang    string `json:"hreflang,omitempty"`    //optional
}

// NewLinkObject initializes a LinkObject with it's required href value.
func NewLinkObject(href string) (*LinkObject, error) {
	if href == "" {
		return nil, errors.New("LinkObject requires a href value")
	}

	return &LinkObject{Href: href}, nil
}

// NewCurieLink initializes a special LinkObject required for establishing CURIEs.
func NewCurieLink(name string, href string) (*LinkObject, error) {
	if name == "" {
		return nil, errors.New("CURIE LinkObject requires a name value")
	}

	linkObject, error := NewLinkObject(href)

	if error != nil {
		return nil, error
	}

	linkObject.Name = name
	linkObject.Templated = true

	return linkObject, nil
}
