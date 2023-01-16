package oas

type WithTags struct {
	Tags []*Tag `json:"tags,omitempty"`
}

func (o *WithTags) AddTag(t *Tag) {
	if t == nil {
		return
	}
	o.Tags = append(o.Tags, t)
}

func NewTag(name string) *Tag {
	return &Tag{
		TagObject: TagObject{
			Name: name,
		},
	}
}

type Tag struct {
	TagObject
	SpecExtensions
}

func (i Tag) MarshalJSON() ([]byte, error) {
	return flattenMarshalJSON(i.TagObject, i.SpecExtensions)
}

func (i *Tag) UnmarshalJSON(data []byte) error {
	return flattenUnmarshalJSON(data, &i.TagObject, &i.SpecExtensions)
}

type TagObject struct {
	Name         string       `json:"name"`
	Description  string       `json:"description,omitempty"`
	ExternalDocs *ExternalDoc `json:"externalDocs,omitempty"`
}
