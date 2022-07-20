// go2hal v0.6.0
// Copyright (c) 2021 Patrick Moule
// License: MIT

package hal

import (
	"encoding/json"
	"net/http"
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

// WriteTo writes the content to a ResponseWriter
func (enc *standardEncoder) WriteTo(w http.ResponseWriter, statusCode int, resource Resource) (int, error) {
	b, err := enc.ToJSON(resource)
	if err != nil {
		return -1, err
	}

	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{JSONMimeType}
	}

	w.WriteHeader(statusCode)
	return w.Write(b)
}
