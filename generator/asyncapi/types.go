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

//TODO: Für jedes individuelle Go-File wird ein passendes Config benötigt
// Immer wenn es immer diesselbe Struktur haben soll und die Files immer gleich sein sollen, dann eine generateBlubFile(conf BlubConf) dafür erstellen
// in /ModuleName/src/cmd/publisher/  kommt rein, was alles für den Channel/Topic benötigt wird: PublisherConfig { ModuleName, MessageName, ChannelName}
// vorher muss vorbereitet werden und in einer map alle Messages, die zu einem Channel gehören gefunden werden
// TODO: Je Channel alle Messages in einem, wenn dann in einer Go-File durch alle Messages iteriert wird, dann kann einfach durch die Liste aller zu diesem Channel gehörigen Messages
// iteriert werden

type PublisherConfig struct {
	ModuleName  string
	MessageName []string //wenn es meherere Messages für den Publisher gibt
	ChannelName string
}

// in /ModuleName/src/cmd/publisher/  kommt rein, was alles für den Pushlisher benötigt wird

// in generator/asyncapi alle unnötigen Go-Files löschen, da nun die Lib Lerenn verwendet wird, also Structs unnötig geworden
// Struktur von openapi mit generateFiles-Funktionen nachzumachen
