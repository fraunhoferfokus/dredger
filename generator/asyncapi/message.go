package generator

type Message struct {
	ContentType   string        `json:"contentType"`
	Headers       *Schema       `json:"headers"`
	Payload       *Schema       `json:"payload"`
	Summary       string        `json:"summary"`
	Name          string        `json:"name"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	CorrelationID CorrelationID `json:"correlationID"`
	Tags          []*Tag        `json:"tags"`
	ExternalDocs  *ExternalDoc  `json:"externalDocs"`
}
