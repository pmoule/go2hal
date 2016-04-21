// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import (
	"encoding/json"
	"errors"
)

type Encoder struct {
}

func (enc *Encoder) ToJSON(resource Resource) ([]byte, error) {
	if mapper, ok := resource.(mapper); ok {
		namedMap := mapper.ToMap()
		return json.Marshal(namedMap.Content)
	}

	return nil, errors.New("Resource is not of type mapper.")
}