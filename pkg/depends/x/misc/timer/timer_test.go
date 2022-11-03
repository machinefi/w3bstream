package timer_test

import (
	"testing"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/x/misc/timer"
)

func TestStartSpan(t *testing.T) {
	span := timer.StartSpan()

	t.Log(span.StartedAt())
	time.Sleep(time.Second)
	t.Log(span.Cost())
	time.Sleep(time.Second)
	t.Log(span.Cost())
}
