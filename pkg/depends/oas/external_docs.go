package oas

func NewExternalDoc(url string, desc string) *ExternalDoc {
	return &ExternalDoc{
		URL:         url,
		Description: desc,
	}
}

type ExternalDoc struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}
