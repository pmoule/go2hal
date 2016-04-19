// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

type mapper interface {
	ToMap() NamedMap
}

type NamedMap struct {
	Name    string
	Content PropertyMap
}

type PropertyMap map[string]interface{}