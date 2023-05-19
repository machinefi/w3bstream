package stringsx_test

import (
	"testing"

	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

func TestGenRandomVisibleString(t *testing.T) {
	t.Log(stringsx.GenRandomVisibleString(0))
}
