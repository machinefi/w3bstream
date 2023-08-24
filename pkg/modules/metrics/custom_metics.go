package metrics

import (
	"context"
	"fmt"

	"github.com/tidwall/gjson"

	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
)

type (
	CustomMetrics interface {
		Submit(gjson.Result) error
	}
)

type (
	metrics struct {
		account string // account use wallet address (if exists) or account id
		project string // project use project name
		writer  *SQLBatcher
	}
)

func NewCustomMetric(account string, project string) CustomMetrics {
	return &metrics{
		account: account,
		project: project,
		writer:  NewSQLBatcher("INSERT INTO ws_metrics.customized_metrics VALUES"),
	}
}

func (m *metrics) Submit(obj gjson.Result) error {
	objStr := obj.String()
	ctx, l := logger.NewSpanContext(context.Background(), "modules.metrics.Submit")
	defer l.End()
	return m.writer.Insert(ctx, fmt.Sprintf(`now(), '%s', '%s', '%s'`, m.account, m.project, objStr))
}
