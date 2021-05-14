package halforms

const (
	// MediaTypeIdentifier should be uses as part of HTTP accept header when requesting
	// a HAL-FORMS document. It should appear as the HTTP content-type header when sending
	// a response that contains a HAL-FORMS document.
	MediaTypeIdentifier     = "application/prs.hal-forms+json"
	jsonMediaTypeIdentifier = "application/json"
	// TemplatesProperty is a reserved name for templates in HAL-FORMS documents.
	TemplatesProperty   string = "_templates"
	TemplateDefaultKey  string = "default"
	PropertyDefaultType string = "text"
)
