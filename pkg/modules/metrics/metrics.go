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
	EventMtc = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: _eventMtcName,
			Help: "received events metrics.",
		},
		[]string{"account", "project", "publisher", "eventtype"},
	)

	PublisherMtc = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: _publisherMtcName,
			Help: "registered publishers for the project.",
		},
		[]string{"account", "project"},
	)

	BlockChainTxMtc = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: _blockChainTxMtcName,
		Help: "blockchain transaction counter metrics.",
	}, []string{"project", "chainID"})
)

func init() {
	prometheus.MustRegister(EventMtc)
	prometheus.MustRegister(PublisherMtc)
	prometheus.MustRegister(BlockChainTxMtc)
}

func RemoveMetrics(ctx context.Context, account string, project string) {
	EventMtc.DeletePartialMatch(prometheus.Labels{"account": account, "project": project})
	PublisherMtc.DeletePartialMatch(prometheus.Labels{"account": account, "project": project})
	BlockChainTxMtc.DeletePartialMatch(prometheus.Labels{"project": project})

	if err := eraseDataInServer(ctx, account, project); err != nil {
		l := types.MustLoggerFromContext(ctx)
		// the metrics server isn't essential for the core service
		l.Warn(err)
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
