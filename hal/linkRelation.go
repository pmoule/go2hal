// go2hal v0.1.0
// Copyright (c) 2016 Patrick Moule
// License: MIT

package hal

import "errors"

type linkRelation struct {
	name string
	curieLink *LinkObject
}

func NewLinkRelation(name string) (*linkRelation, error) {
	if name == "" {
		return nil, errors.New("LinkRelation requires a name value.")
	}

	return &linkRelation{name: name}, nil
}

func (lr *linkRelation) Name() string {
	if lr.curieLink == nil {
		return lr.name
	}

	return lr.curieLink.Name + "." +lr.name
}

func (lr *linkRelation) SetCurieLink(curieLink *LinkObject) {
	lr.curieLink = curieLink
}

func (lr *linkRelation) CurieLink() LinkObject {
	return *lr.curieLink
}