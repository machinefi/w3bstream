package operators

type (
	Generator interface {
		GenCode() (string, error)
		IsValidSuccessor(Generator) bool
	}
	HeadGenerator struct{}
)

func (hg *HeadGenerator) GenCode() (string, error) {
	return "", nil
}

func (hg *HeadGenerator) IsValidSuccessor(Generator) bool {
	return true
}
