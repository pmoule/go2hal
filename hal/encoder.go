// go2hal v0.2.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import (
	"encoding/json"
	"errors"
)

// Encoder to encode a Resource Object into a valid HAL document.
type Encoder interface {
	ToJSON(resource Resource) ([]byte, error)
}

type standardEncoder struct {
}

type advancedEncoder struct {
	standardEncoder
}

func NewEncoder() Encoder {
	return new(standardEncoder)
}

func NewAdvancedEncoder() Encoder {
	return new(advancedEncoder)
}

// ToJSON generates a HAL document from given Resource Object.
func (enc *standardEncoder) ToJSON(resource Resource) ([]byte, error) {
	if mapper, ok := resource.(mapper); ok {
		namedMap := mapper.ToMap()
		return json.Marshal(namedMap.Content)
	}

	return nil, errors.New("Resource is not of type mapper")
}
