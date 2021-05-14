// go2hal v0.5.0
// Copyright (c) 2021 Patrick Moule
// License: MIT

package hal

import (
	"encoding/json"
)

// Encoder to encode a Resource into a valid HAL document.
type Encoder interface {
	ToJSON(resource Resource) ([]byte, error)
}

type standardEncoder struct {
}

// NewEncoder creates a JSON encoder
func NewEncoder() Encoder {
	return new(standardEncoder)
}

// ToJSON generates a HAL document from provided Resource.
func (enc *standardEncoder) ToJSON(resource Resource) ([]byte, error) {
	namedMap := resource.ToMap()

	return json.Marshal(namedMap.Content)
}
