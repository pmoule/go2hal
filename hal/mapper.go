package hal

type mapper interface {
	ToMap() NamedMap
}

type NamedMap struct {
	Name    string
	Content PropertyMap
}

type PropertyMap map[string]interface{}