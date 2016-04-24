// go2hal v0.2.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import (
	"encoding/json"
	"errors"
)

// Encoder to encode a Resource Object into a valid HAL document.
type Encoder struct {
}

// ToJSON generates a HAL document from given Reosurce Object.
// The output media type is "application/hal+json"
func (enc *Encoder) ToJSON(resource Resource) ([]byte, error) {
	if mapper, ok := resource.(mapper); ok {
		namedMap := mapper.ToMap()
		return json.Marshal(namedMap.Content)
	}

	return nil, errors.New("Resource is not of type mapper")
}