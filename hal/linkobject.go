// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import "errors"

// LinkObject is a hyperlink from the Resource it is attached to.
// A valid LinkObject requires a Href value. All other properties are optional.
// See https://tools.ietf.org/html/draft-kelly-json-hal-07 for
// property description.
type LinkObject struct {
	Href        string `json:"href,omitempty"`        //required
	Templated   bool `json:"templated,omitempty"`     //optional
	Type        string `json:"type,omitempty"`        //optional
	Deprecation string `json:"deprecation,omitempty"` //optional
	Name        string `json:"name,omitempty"`        //optional
	Profile     string `json:"profile,omitempty"`     //optional
	Title       string `json:"title,omitempty"`       //optional
	HrefLang    string `json:"hreflang,omitempty"`    //optional
}

// NewCurieLink initializes a special LinkObject required for establishing CURIEs.
func NewCurieLink(name string, href string) (*LinkObject, error) {
	if name == "" {
		return nil, errors.New("Curie link requires a name value.")
	}

	if href == "" {
		return nil, errors.New("Curie link requires a href value.")
	}

	return &LinkObject{Name: name, Href: href, Templated: true}, nil
}