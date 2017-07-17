// go2hal v0.3.0
// Copyright (c) 2017 Patrick Moule
// License: MIT

package hal

import (
	"testing"
)

func TestNewCurieLink(t *testing.T) {
	wantedName := "link"
	wantedHref := "http://{rel}"

	wantedCurieLink, _ := NewCurieLink(wantedName, wantedHref)

	if wantedCurieLink.Name != wantedName {
		t.Errorf("Link name == %q, want %q", wantedCurieLink.Name, wantedName)
	}

	if wantedCurieLink.Href != wantedHref {
		t.Errorf("Link href == %q, want %q", wantedCurieLink.Href, wantedHref)
	}

	if !wantedCurieLink.Templated {
		t.Errorf("Link templated == %q, want %q", wantedCurieLink.Templated, true)
	}

	_, invalidNameError := NewCurieLink("", wantedHref)

	if invalidNameError == nil {
		t.Errorf("NewCurieLink should return an error due to an invalid name value.")
	}

	_, invalidHrefError := NewCurieLink(wantedName, "")

	if invalidHrefError == nil {
		t.Errorf("NewCurieLink should return an error due to an invalid href value.")
	}
}
