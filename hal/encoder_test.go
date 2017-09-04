// go2hal v0.3.2
// Copyright (c) 2017 Patrick Moule
// License: MIT

package hal

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pmoule/go2hal/hal/relationtype"
)

type Actor struct {
	Id   int
	Name string
}

func TestEncoder(t *testing.T) {
	wanted := []interface{}{2, "{}"}

	encoder := NewEncoder()

	resource := NewResourceObject()
	bytes, _ := encoder.ToJSON(resource)

	if count := len(bytes); count != wanted[0] {
		t.Errorf("ToJSON returns zero byte == %d, want %d", count, wanted[0])
	}

	if value := string(bytes[:]); value != wanted[1] {
		t.Errorf("JSON value == %s, want %s", value, wanted[1])
	}
}

func TestEncoderWithSelfLink(t *testing.T) {
	wanted := []string{"/docwhoapi/doctors"}

	encoder := NewEncoder()

	root := NewResourceObject()
	link := &LinkObject{Href: wanted[0]}
	self, _ := NewLinkRelation(relationtype.Self)
	self.SetLink(link)
	root.AddLink(self)

	bytes, _ := encoder.ToJSON(root)

	var decoded interface{}
	err := json.Unmarshal(bytes, &decoded)

	if err != nil {
		t.Errorf("Unmarshalling JSON returns error: %d", err)
	}

	result := decoded.(map[string]interface{})

	val, ok := result[LinksProperty].(map[string]interface{})

	if !ok {
		t.Errorf("Generated JSON does not contain link property %s:", LinksProperty)
	}

	val, ok = val[relationtype.Self].(map[string]interface{})

	if !ok {
		t.Errorf("Generated JSON does not contain self relation: %s", relationtype.Self)
	}

	uri, uriOk := val["href"].(string)

	if !uriOk {
		t.Errorf("Generated JSON does not contain href property: %s", "href")
	}

	if uri != wanted[0] {
		t.Errorf("Generated JSON uri: %s, want:  %s", uri, wanted[0])
	}
}

func TestEncoderWithEmbeddedResources(t *testing.T) {
	wanted := []string{"/docwhoapi/doctors", "doctors", "name"}

	encoder := NewEncoder()

	root := NewResourceObject()
	link := &LinkObject{Href: wanted[0]}
	self, _ := NewLinkRelation(relationtype.Self)
	self.SetLink(link)
	root.AddLink(self)

	actors := []Actor{
		Actor{1, "William Hartnell"},
		Actor{2, "Patrick Troughton"},
	}

	var embeddedActors []Resource

	for _, actor := range actors {
		href := fmt.Sprintf("%s/%d", wanted[0], actor.Id)
		selfLink, _ := NewLinkObject(href)

		self, _ := NewLinkRelation(relationtype.Self)
		self.SetLink(selfLink)

		embeddedActor := NewResourceObject()
		embeddedActor.AddLink(self)
		embeddedActor.Data()[wanted[2]] = actor.Name
		embeddedActors = append(embeddedActors, embeddedActor)
	}

	doctors, _ := NewResourceRelation(wanted[1])
	doctors.SetResources(embeddedActors)

	root.AddResource(doctors)

	bytes, _ := encoder.ToJSON(root)

	var decoded interface{}
	err := json.Unmarshal(bytes, &decoded)

	if err != nil {
		t.Errorf("Unmarshalling JSON returns error: %d", err)
	}

	result := decoded.(map[string]interface{})

	val, ok := result[EmbeddedProperty].(map[string]interface{})

	if !ok {
		t.Errorf("Generated JSON does not contain embedded property %s:", EmbeddedProperty)
	}

	doctorsVal, doctorsOk := val[wanted[1]].([]interface{})

	if !doctorsOk {
		t.Errorf("Generated JSON does not contain doctors relation: %s", wanted[1])
		t.Errorf("JSON value == %s", string(bytes[:]))
		t.Errorf("JSON unmarshalled value == %s", result)
	}

	if len(doctorsVal) != len(actors) {
		t.Errorf("Generated JSON does not contain expected amount of doctors: %d", len(actors))
	}

	if doctorsVal[0].(map[string]interface{})["name"] != actors[0].Name {
		t.Errorf("Generated JSON does not contain expected doctor: %s", actors[0].Name)
	}

	if doctorsVal[1].(map[string]interface{})["name"] != actors[1].Name {
		t.Errorf("Generated JSON does not contain expected doctor: %s", actors[1].Name)
	}
}
