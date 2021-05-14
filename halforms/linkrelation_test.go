package halforms

import "testing"

func TestNewHALFormsRelation(t *testing.T) {
	relName := ""
	href := ""
	_, err := NewHALFormsRelation(relName, href)

	if err == nil {
		t.Errorf("error should not be nil")
	}

	relName = "rel"

	_, err = NewHALFormsRelation(relName, href)

	if err == nil {
		t.Errorf("error should not be nil")
	}

	relName = ""
	href = "href"

	_, err = NewHALFormsRelation(relName, href)

	if err == nil {
		t.Errorf("error should not be nil")
	}

	relName = "rel"
	href = "href"

	rel, err := NewHALFormsRelation(relName, href)

	if err != nil {
		t.Errorf("error should be nil")
	}

	if rel == nil {
		t.Errorf("relation should not be nil")
	}

	if rel.Name() != relName {
		t.Errorf("relation name is %s, want %s", rel.Name(), relName)
	}

	if count := len(rel.Links()); count != 1 {
		t.Errorf("links count is %d, want %d", count, 1)
	}

	link := rel.Links()[0]

	if link.Href != href {
		t.Errorf("href is %s, want %s", link.Href, href)
	}

	if link.Type != MediaTypeIdentifier {
		t.Errorf("type is %s, want %s", link.Href, MediaTypeIdentifier)
	}
}
