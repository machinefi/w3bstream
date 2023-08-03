package metrics

import (
	"fmt"

	"github.com/tidwall/gjson"
)

type (
	CustomMetrics interface {
		Submit(gjson.Result) error
	}
)

type (
	metrics struct {
		account string
		project string
		writer  *BatchWorker
	}
)

func NewCustomMetric(account string, project string) CustomMetrics {
	return &metrics{
		account: account,
		project: project,
		writer:  NewBatchWorker("INSERT INTO ws_metrics.customized_metrics VALUES"),
	}
}

func (m *metrics) Submit(obj gjson.Result) error {
	objStr := obj.String()
	return m.writer.Insert(fmt.Sprintf(`now(), '%s', '%s', '%s'`, m.account, m.project, objStr))
}
