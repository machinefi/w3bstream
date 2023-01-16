package types

type PullPolicy string

const (
	PullAlways       PullPolicy = "Always"
	PullNever        PullPolicy = "Never"
	PullIfNotPresent PullPolicy = "IfNotPresent"
)
