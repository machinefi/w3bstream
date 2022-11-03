package httptransport_test

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	. "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
)

func TestServiceMeta(t *testing.T) {
	meta := &ServiceMeta{}

	name, version := "srv-test", "1.1.1"

	os.Setenv(consts.EnvProjectName, name)
	meta.SetDefault()
	NewWithT(t).Expect(meta.String()).To(Equal(name))

	os.Setenv(consts.EnvProjectVersion, version)
	meta.SetDefault()
	NewWithT(t).Expect(meta.String()).To(Equal(name + "@" + version))
}

func TestServiceMetaWithContext(t *testing.T) {
	meta := ServiceMeta{Name: "test"}
	ctx := ContextWithServiceMeta(context.Background(), meta)
	got := ServiceMetaFromContext(ctx)
	NewWithT(t).Expect(got).To(Equal(meta))
}

var (
	bgctx = context.Background()
)
