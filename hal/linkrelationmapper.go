// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

type linkRelationMapper map[LinkRelation]interface{}

func (lrm *linkRelationMapper) ToMap() NamedMap {
	linkMap := PropertyMap{}

	for key, val := range *lrm {
		linkMap[key.FullName()] = val
	}

	return NamedMap{Name: "_links", Content: linkMap}
}