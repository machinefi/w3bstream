package mqtt

//go:generate toolkit gen enum QOS
type QOS int8

const (
	QOS_UNKNOWN        QOS = iota - 1
	QOS__ONCE              // 0
	QOS__AT_LEAST_ONCE     // 1
	QOS__ONLY_ONCE         // 2
)

// func (QoS) Offset() int { return -1 }
