package oas

type Info struct {
	InfoObject
	SpecExtensions
}

func (i Info) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(i.InfoObject, i.SpecExtensions)
}

func (i *Info) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &i.InfoObject, &i.SpecExtensions)
}

type InfoObject struct {
	Title          string `json:"title"`
	Description    string `json:"description,omitempty"`
	TermsOfService string `json:"termsOfService,omitempty"`
	*Contact       `json:"contact,omitempty"`
	*License       `json:"license,omitempty"`
	Version        string `json:"version"`
}

type Contact struct {
	ContactObject
	SpecExtensions
}

func (i Contact) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(i.ContactObject, i.SpecExtensions)
}

func (i *Contact) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &i.ContactObject, &i.SpecExtensions)
}

type ContactObject struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	LicenseObject
	SpecExtensions
}

func (i License) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(i.LicenseObject, i.SpecExtensions)
}

func (i *License) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &i.LicenseObject, &i.SpecExtensions)
}

type LicenseObject struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}
