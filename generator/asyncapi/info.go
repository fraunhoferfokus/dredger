package generator

type Info struct {
	Title          string `json:"title"`
	Version        string `json:"version"`
	Description    string `json:"description"`
	TermsOfService string `json:"termsOfService"`
	Contact        string `json:"contact"`
	Tags           []*Tag `json:"tags"`
	//TODO rest of the fields for
}
