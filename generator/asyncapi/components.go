package generator

type Components struct {
	Schemas           map[string]*Schema           `json:"schemas"`
	Servers           map[string]*Server           `json:"servers"`
	Channels          map[string]*Channel          `json:"channels"`
	ServerVariables   map[string]*ServerVariable   `json:"serverVariables"`
	Operations        map[string]*Operation        `json:"operations"`
	Messages          map[string]*Message          `json:"messages"`
	SecuritySchemas   map[string]*SecuritySchema   `json:"securitySchemes"`
	Parameters        map[string]*Parameter        `json:"parameters"`
	CorrelationID     map[string]*CorrelationID    `json:"correlationIds"`
	OperationTraits   map[string]*OperationTrait   `json:"operationTraits"`
	MessageTraits     map[string]*MessageTrait     `json:"messageTraits"`
	Replies           map[string]*Reply            `json:"replies"`
	ReplyAddresses    map[string]*ReplyAddress     `json:"replyAddresses"`
	ServerBindings    map[string]*ServerBinding    `json:"serverBindings"`
	ChannelBindings   map[string]*ChannelBinding   `json:"channelBindings"`
	OperationBindings map[string]*OperationBinding `json:"operationBindings"`
	MessageBindings   map[string]*MessageBinding   `json:"messageBindings"`
	Tags              map[string]*Tag              `json:"tags"`
	ExternalDocs      map[string]*ExternalDoc      `json:"externalDocs"`
}
