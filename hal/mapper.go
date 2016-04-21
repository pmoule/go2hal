// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

type mapper interface {
	ToMap() NamedMap
}

// NamedMap simply links a name with PropertyMap
type NamedMap struct {
	Name    string
	Content PropertyMap
}

// PropertyMap simply maps a string to any kind of value.
type PropertyMap map[string]interface{}