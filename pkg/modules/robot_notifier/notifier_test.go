package robot_notifier_test

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/modules/robot_notifier"
	"github.com/machinefi/w3bstream/pkg/modules/robot_notifier/lark"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestPush(t *testing.T) {
	config := &types.RobotNotifierConfig{
		URL:    "https://open.larksuite.com/open-apis/bot/v2/hook/f8d7cd45-4b45-40fe-9635-5e2f85e19155",
		Secret: "vztL7BIOyDw10XEd9H5B6",
		Env:    "dev-staging",
	}
	config.Init()

	ctx := types.WithRobotNotifierConfig(context.Background(), config)

	body, err := lark.Build(ctx, "test_title", "warning", "message content")
	NewWithT(t).Expect(err).To(BeNil())
	t.Log(string(body))

	err = robot_notifier.Push(ctx, body, lark.ResponseHook)
	t.Log(err)
}
