package generator

type AsyncAPIConfig struct {
	AsyncAPIPath string
	OutputPath   string
	ModuleName   string
}

type BrokerConfig struct {
	Protocol string // z.B. "nats", "kafka", "mqtt"
	URL      string // z.B. "nats://localhost:4222"
}

type ChannelConfig struct {
	ModuleName  string
	ChannelName string // Topic/Subject Name
	Description string
	Action      *OperationConfig // optional
}

type OperationConfig struct {
	Name    string
	Summary string
	Message MessageConfig
}

type MessageConfig struct {
	MessageName string
	//Description string
	//	Schema      string // z.B. Verweis auf eine Go-Struct
}

type AsyncAPISpec struct {
	Title       string
	Version     string
	Description string
	Broker      BrokerConfig
	Channels    []ChannelConfig
}

type AsyncGenConfig struct {
	ModuleName  string
	ChannelName string // Topic/Subject Name
	Description string
	MessageName string
	Channels    []string
	Action      string
}
