// go2hal v0.3.2
// Copyright (c) 2017 Patrick Moule
// License: MIT

package hal

// NamedMap simply links a name with PropertyMap
type NamedMap struct {
	Name    string
	Content PropertyMap
}

// PropertyMap simply maps a string to any kind of value.
type PropertyMap map[string]interface{}
