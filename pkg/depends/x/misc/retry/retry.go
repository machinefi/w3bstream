package retry

import (
	"log"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
)

type Retry struct {
	Repeats  int
	Interval types.Duration
}

func (r *Retry) SetDefault() {
	if r.Repeats == 0 {
		r.Repeats = Default.Repeats
	}
	if r.Interval == 0 {
		r.Interval = Default.Interval
	}
}

func (r Retry) Do(exec func() error) (err error) {
	if r.Repeats <= 0 {
		return exec()
	}
	for i := 0; i < r.Repeats; i++ {
		if err = exec(); err != nil {
			log.Printf("retry in seconds %s [err: %v]", r.Interval.Duration(), err)
			time.Sleep(time.Duration(r.Interval))
			continue
		}
		break
	}
	return
}

var Default = &Retry{3, types.Duration(3 * time.Second)}

func Do(retry *Retry, exec func() error) error {
	return retry.Do(exec)
}
