// go2hal v0.6.0
// Copyright (c) 2021 Patrick Moule
// License: MIT

package halforms

import "github.com/skhaz/go2hal/hal"

// NewHALFormsRelation creates a relation from provided relation name and href.
// The contained LinkObject's type property has application/prs.hal-forms+json as
// it's value to indicate a HAL-FORMS document as the expected media type.
func NewHALFormsRelation(name string, href string) (hal.LinkRelation, error) {
	relation, err := hal.NewLinkRelation(name)

	if err != nil {
		return nil, err
	}

	link, err := hal.NewLinkObject(href)

	if err != nil {
		return nil, err
	}

	link.Type = MediaTypeIdentifier
	relation.SetLink(link)

	return relation, err
}
