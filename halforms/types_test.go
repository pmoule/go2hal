package halforms

import (
	"encoding/json"
	"testing"

	"github.com/pmoule/go2hal/hal"
	"github.com/pmoule/go2hal/hal/relationtype"
)

func TestDocument(t *testing.T) {
	href := "www.example.com"
	document := NewDocument(href)
	template := NewTemplate()
	document.AddTemplate(template)

	options := &Options{}
	options.Inline = append(options.Inline, struct {
		Number int
		Text   string
	}{
		42,
		"Test value",
	})
	property := NewProperty("property")
	property.Options = options
	template.Properties = append(template.Properties, property)

	encoder := NewEncoder()
	jsonDocument, _ := encoder.ToJSON(document)

	var decoded interface{}
	err := json.Unmarshal(jsonDocument, &decoded)

	if err != nil {
		t.Errorf("Unmarshalling JSON returns error: %d", err)
	}

	result := decoded.(map[string]interface{})

	val, ok := result[hal.LinksProperty].(map[string]interface{})

	if !ok {
		t.Errorf("Generated JSON does not contain property %s:", hal.LinksProperty)
	}

	selfValue, ok := val[relationtype.Self].(map[string]interface{})

	if !ok {
		t.Errorf("Generated JSON does not contain property %s", relationtype.Self)
	}

	hrefValue, ok := selfValue["href"].(string)

	if !ok {
		t.Errorf("Generated JSON does not contain property: %s", "href")
	}

	if hrefValue != href {
		t.Errorf("Generated JSON href: %s, want:  %s", hrefValue, href)
	}

	typeValue, ok := selfValue["type"].(string)

	if !ok {
		t.Errorf("Generated JSON does not contain property: %s", "type")
	}

	if typeValue != MediaTypeIdentifier {
		t.Errorf("Generated JSON type: %s, want:  %s", typeValue, MediaTypeIdentifier)
	}
}
