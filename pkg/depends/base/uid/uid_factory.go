package uid

import "time"

type factory struct {
	bitsWorker uint
	bitsTs     uint
	bitsSeq    uint
	unit       time.Duration
}
