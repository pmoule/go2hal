// go2hal v0.6.0
// Copyright (c) 2021 Patrick Moule
// License: MIT

package halforms

import "encoding/json"

type Encoder interface {
	ToJSON(document Document) ([]byte, error)
}

type standardEncoder struct {
}

// NewEncoder creates a JSON encoder
func NewEncoder() Encoder {
	return new(standardEncoder)
}

// ToJSON generates a HAL-FORMS document from provided Document.
func (enc *standardEncoder) ToJSON(document Document) ([]byte, error) {
	namedMap := document.ToMap()

	return json.Marshal(namedMap.Content)
}
