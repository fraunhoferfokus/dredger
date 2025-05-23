package generator

type Server struct {
	Host        string `json:"host"`
	Description string `json:"description"`
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	Bindings    *ServerBinding
	Protocol    string `json:"protocol"`
	//TODO Rest der gesamten Felder
}
