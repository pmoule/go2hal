// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

type embeddedResourceMapper map[linkRelation]interface{}

func (erm *embeddedResourceMapper) ToMap() NamedMap {
	embeddedMap := PropertyMap{}

	for key, val := range *erm {
		embeddedMap[key.Name()] = val
	}

	for key, val := range embeddedMap {
		if resourceArray, isSlice := val.([]Resource); isSlice {
			var resources []interface{}

			for _, resource := range resourceArray {
				mapper, ok := resource.(mapper)

				if !ok {
					continue
				}

				namedMap := mapper.ToMap()
				resources = append(resources, namedMap.Content)
			}

			embeddedMap[key] = resources

		} else {
			resource := val.(Resource)
			mapper, ok := resource.(mapper)

			if !ok {
				continue
			}

			namedMap := mapper.ToMap()
			embeddedMap[key] = namedMap.Content
		}
	}

	return NamedMap{Name: "_embedded", Content: embeddedMap }
}