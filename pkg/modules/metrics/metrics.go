package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/machinefi/w3bstream/pkg/types"
)

const (
	_eventMtcName        = "inbound_events_metrics"
	_publisherMtcName    = "publishers_metrics"
	_blockChainTxMtcName = "w3b_blockchain_tx_metrics"
)

var (
	eventMtc = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: _eventMtcName,
			Help: "received events metrics.",
		},
		[]string{"account", "project", "publisher", "eventtype"},
	)
	eventClickhouseCli = NewSQLBatcher("INSERT INTO ws_metrics.inbound_events_metrics VALUES")

	publisherMtc = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: _publisherMtcName,
			Help: "registered publishers for the project.",
		},
		[]string{"account", "project"},
	)
	publisherClickhouseCli = NewSQLBatcher("INSERT INTO ws_metrics.publishers_metrics VALUES")

	BlockChainTxMtc = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: _blockChainTxMtcName,
		Help: "blockchain transaction counter metrics.",
	}, []string{"project", "chainID"})
)

func init() {
	prometheus.MustRegister(eventMtc)
	prometheus.MustRegister(publisherMtc)
	prometheus.MustRegister(BlockChainTxMtc)
}

func RemoveMetrics(ctx context.Context, account string, project string) {
	eventMtc.DeletePartialMatch(prometheus.Labels{"account": account, "project": project})
	publisherMtc.DeletePartialMatch(prometheus.Labels{"account": account, "project": project})
	BlockChainTxMtc.DeletePartialMatch(prometheus.Labels{"project": project})

	// erase data in metrics server
	if err := eraseDataInServer(ctx, account, project); err != nil {
		l := types.MustLoggerFromContext(ctx)
		// the metrics server isn't essential for the core service
		l.Warn(err)
	}
	if clickhouseCLI != nil {
		if err := clickhouseCLI.Insert(fmt.Sprintf(`DELETE FROM ws_metrics.inbound_events_metrics WHERE (
			account = '%s') AND (project = '%s')`, account, project)); err != nil {
			l := types.MustLoggerFromContext(ctx)
			l.Warn(err)
		}
		if err := clickhouseCLI.Insert(fmt.Sprintf(`DELETE FROM ws_metrics.publishers_metrics WHERE (
				account = '%s') AND (project = '%s')`, account, project)); err != nil {
			l := types.MustLoggerFromContext(ctx)
			l.Warn(err)
		}
		if err := clickhouseCLI.Insert(fmt.Sprintf(`DELETE FROM ws_metrics.customized_metrics WHERE (
			account = '%s') AND (project = '%s')`, account, project)); err != nil {
			l := types.MustLoggerFromContext(ctx)
			l.Warn(err)
		}
	}
}

func EventMetricsInc(ctx context.Context, account, project, publisher, eventtype string) {
	eventMtc.WithLabelValues(account, project, publisher, eventtype).Inc()
	if clickhouseCLI != nil {
		if err := eventClickhouseCli.Insert(fmt.Sprintf(`now(), '%s', '%s', '%s', 
		'%s', %d`, account, project, publisher, eventtype, 1)); err != nil {
			l := types.MustLoggerFromContext(ctx)
			l.Error(err)
		}
	}
}

func PublisherMetricsInc(ctx context.Context, account, project string) {
	publisherMtc.WithLabelValues(account, project).Inc()
	if clickhouseCLI != nil {
		if err := publisherClickhouseCli.Insert(fmt.Sprintf(`now(), '%s', '%s', %d`, account, project, 1)); err != nil {
			l := types.MustLoggerFromContext(ctx)
			l.Error(err)
		}
	}
}

func PublisherMetricsDec(ctx context.Context, account, project string) {
	publisherMtc.WithLabelValues(account, project).Dec()
	if clickhouseCLI != nil {
		if err := publisherClickhouseCli.Insert(fmt.Sprintf(`now(), '%s', '%s', %d`, account, project, -1)); err != nil {
			l := types.MustLoggerFromContext(ctx)
			l.Error(err)
		}
	}
}

func eraseDataInServer(ctx context.Context, account string, project string) error {
	cfg, existed := types.MetricsCenterConfigFromContext(ctx)
	if !existed {
		return errors.New("fail to get the url of metrics center")
	}
	baseURL := cfg.Endpoint

	if err := httpReq(
		fmt.Sprintf("%s/api/v1/admin/tsdb/delete_series?match[]=%s",
			baseURL,
			url.QueryEscape(fmt.Sprintf(`%s{account="%s", project="%s"}`, _eventMtcName, account, project))),
	); err != nil {
		return err
	}

	if err := httpReq(
		fmt.Sprintf("%s/api/v1/admin/tsdb/delete_series?match[]=%s",
			baseURL,
			url.QueryEscape(fmt.Sprintf(`%s{account="%s", project="%s"}`, _publisherMtcName, account, project))),
	); err != nil {
		return err
	}

	if err := httpReq(
		fmt.Sprintf("%s/api/v1/admin/tsdb/delete_series?match[]=%s",
			baseURL,
			url.QueryEscape(fmt.Sprintf(`%s{project="%s"}`, _blockChainTxMtcName, project))),
	); err != nil {
		return err
	}
	return cleanTombstones(baseURL)
}

func cleanTombstones(baseURL string) error {
	return httpReq(fmt.Sprintf("%s/api/v1/admin/tsdb/clean_tombstones", baseURL))
}

func httpReq(url string) error {
	// Create the HTTP client and request
	client := http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return errors.New("the http request to metrics center fails")
	}

	return nil
}
