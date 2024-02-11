// go2hal v0.6.0
// Copyright (c) 2021 Patrick Moule
// License: MIT

package hal

// EmbeddedProperty is a reserved name for embedding resources in HAL documents.
const EmbeddedProperty string = "_embedded"

// LinksProperty is a reserved name for embedding Link Objects in HAL documents.
const LinksProperty string = "_links"

// HALJSONMimeType is the mimetype of HAL JSON documents
const HALJSONMimeType string = "application/hal+json; charset=utf-8"
