package blockchain

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/enums"
)

const chainUniqFlag = 0

func RemoveMonitor(ctx context.Context, projectName string) error {
	if err := removeContractLogByProject(ctx, projectName); err != nil {
		return err
	}
	if err := removeChainTxByProject(ctx, projectName); err != nil {
		return err
	}
	return removeChainHeightByProject(ctx, projectName)
}

func getEventType(eventType string) string {
	if eventType == "" {
		return enums.MONITOR_EVENTTYPEDEFAULT
	}
	return eventType
}
