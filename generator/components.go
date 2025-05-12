package generator

type Components struct {
	schemas           map[string]*Schema
	servers           map[string]*Server
	channels          map[string]*Channel
	serverVariables   map[string]*ServerVariable
	operations        map[string]*Operation
	messages          map[string]*Message
	securitySchemas   map[string]*SecuritySchema
	parameters        map[string]*Parameter
	correlationID     map[string]*CorrelationID
	operationTraits   map[string]*OperationTrait
	messageTraits     map[string]*MessageTrait
	replies           map[string]*Reply
	replyAddresses    map[string]*ReplyAddress
	serverBindings    map[string]*ServerBinding
	channelBindings   map[string]*ChannelBinding
	operationBindings map[string]*OperationBinding
	messageBindings   map[string]*MessageBinding
	tags              map[string]*Tag
	externalDocs      map[string]*ExternalDoc
}
