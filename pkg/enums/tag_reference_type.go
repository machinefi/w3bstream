package enums

//go:generate toolkit gen enum TagReferenceType
type TagReferenceType uint8

const (
	TAG_REFERENCE_TYPE_UNKNOWN TagReferenceType = iota
	TAG_REFERENCE_TYPE__PROJECT
)
