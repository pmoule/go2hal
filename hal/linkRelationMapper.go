package hal

type linkRelationMapper map[linkRelation]interface{}

func (lrm *linkRelationMapper) ToMap() NamedMap {
	linkMap := PropertyMap{}

	for key, val := range *lrm {
		linkMap[key.Name()] = val
	}

	return NamedMap{Name: "_links", Content: linkMap}
}