package generator

type Specification struct {
	Version            string                `json:"asyncapi"`
	Info               Info                  `json:"info"`
	Servers            map[string]*Server    `json:"servers,omitempty"`
	Channels           map[string]*Channel   `json:"channels,omitempty"`
	DefaultContentType string                `json:"defaultContentType,omitempty"`
	Operations         map[string]*Operation `json:"operations,omitempty"`
	Components         Components            `json:"components"`
}
