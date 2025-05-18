package generator

import (
	"encoding/json"
)

type Channel struct {
	Address   string                       `json:"address"`
	Messages  map[string]MessageDefinition `json:"messages"`
	Publish   *Operation                   `json:"publish,omitempty"`
	Subscribe *Operation                   `json:"subscribe,omitempty"`
}

type MessageDefinition struct {
	ContentType string          `json:"contentType"`
	Payload     json.RawMessage `json:"payload"` // flexible structure
}
