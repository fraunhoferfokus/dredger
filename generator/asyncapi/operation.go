package generator

type Operation struct {
	Message MessageRef `json:"message"`
}

type MessageRef struct {
	Ref string `json:"$ref"`
}
