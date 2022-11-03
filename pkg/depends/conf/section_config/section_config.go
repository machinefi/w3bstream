package section_config

import "fmt"

type Section struct {
	Name  string
	Value string
}

func (s Section) String() string {
	return fmt.Sprintf("[%s:%s]", s.Name, s.Value)
}

type SectionConfig interface {
	GetSection() *Section
	SetSection(string, string)
}
