package metrics

import (
	"errors"
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
	}
)

func NewCustomMetric(account string, project string) CustomMetrics {
	return &metrics{
		account: account,
		project: project,
	}
}

func (m *metrics) Submit(obj gjson.Result) error {
	if clickhouseCLI == nil {
		return errors.New("clickhouse client is not initialized")
	}

	objStr := obj.String()
	return clickhouseCLI.Insert(fmt.Sprintf(`INSERT INTO ws_metrics.customized_metrics VALUES (
			now(), '%s', '%s', '%s')`, m.account, m.project, objStr))
}
