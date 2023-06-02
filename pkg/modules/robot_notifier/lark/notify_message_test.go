package lark_test

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/modules/robot_notifier/lark"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestTagElementToMap(t *testing.T) {
	text := "test_text"
	vmap, err := lark.TagElementToMap(&lark.ElementText{Text: text})
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(vmap).To(Equal(map[string]interface{}{
		"tag":  "text",
		"text": text,
	}))
}

func TestBuild(t *testing.T) {
	config := &types.RobotNotifierConfig{
		URL:    "https://open.larksuite.com/open-apis/bot/v2/hook/f8d7cd45-4b45-40fe-9635-5e2f85e19155",
		Secret: "vztL7BIOyDw10XEd9H5B6",
		Env:    "prod",
	}
	config.Init()

	ctx := types.WithRobotNotifierConfig(context.Background(), config)
	msg, err := lark.Build(ctx, "TITLE", "warning", "message content")
	NewWithT(t).Expect(err).To(BeNil())
	t.Log(string(msg))
}
