package http_test

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/conf/http"
)

type SomeChecker struct{}

func (SomeChecker) LivenessCheck() map[string]string {
	return map[string]string{
		"mw1_1": "ok",
		"mw1_2": "failed",
	}
}

type WrapChecker struct{ SomeChecker }

func (WrapChecker) LivenessCheck() map[string]string {
	return map[string]string{
		"mw1_1": "ok",
		"mw1_2": "failed",
		"mw2_1": "failed",
		"mw2_2": "ok",
	}
}

type Config struct {
	SomeChecker
	WrapChecker
}

type ConfigWithChecker struct {
	SomeChecker
	WrapChecker
}

func (ConfigWithChecker) LivenessCheck() map[string]string {
	return map[string]string{
		"SomeChecker/mw1_1": "ok",
		"SomeChecker/mw1_2": "failed",
		"WrapChecker/mw1_1": "ok",
		"WrapChecker/mw1_2": "failed",
		"WrapChecker/mw2_1": "failed",
		"WrapChecker/mw2_2": "ok",
	}
}

func TestRegisterCheckerBy(t *testing.T) {
	t.Run("ForEachInputChecker", func(t *testing.T) {
		RegisterCheckerBy(SomeChecker{}, WrapChecker{})

		statuses, err := (&Liveness{}).Output(bgctx)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(statuses).To(Equal(map[string]string{
			"SomeChecker/mw1_1": "ok",
			"SomeChecker/mw1_2": "failed",
			"WrapChecker/mw1_1": "ok",
			"WrapChecker/mw1_2": "failed",
			"WrapChecker/mw2_1": "failed",
			"WrapChecker/mw2_2": "ok",
		}))
	})

	t.Run("ForEachField", func(t *testing.T) {
		RegisterCheckerBy(Config{})

		statuses, err := (&Liveness{}).Output(bgctx)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(statuses).To(Equal(map[string]string{
			"SomeChecker/mw1_1": "ok",
			"SomeChecker/mw1_2": "failed",
			"WrapChecker/mw1_1": "ok",
			"WrapChecker/mw1_2": "failed",
			"WrapChecker/mw2_1": "failed",
			"WrapChecker/mw2_2": "ok",
		}))
	})

	t.Run("CoveredSubLivenessChecker", func(t *testing.T) {
		ResetRegistered()
		RegisterCheckerBy(ConfigWithChecker{})

		statuses, err := (&Liveness{}).Output(bgctx)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(statuses).To(Equal(map[string]string{
			"ConfigWithChecker/SomeChecker/mw1_1": "ok",
			"ConfigWithChecker/SomeChecker/mw1_2": "failed",
			"ConfigWithChecker/WrapChecker/mw1_1": "ok",
			"ConfigWithChecker/WrapChecker/mw1_2": "failed",
			"ConfigWithChecker/WrapChecker/mw2_1": "failed",
			"ConfigWithChecker/WrapChecker/mw2_2": "ok",
		}))
	})
}

var bgctx = context.Background()
