package metrics

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"

	"github.com/machinefi/w3bstream/pkg/types"
)

var autoCollectCli = NewSQLBatcher("INSERT INTO ws_metrics.auto_collect_metrics VALUES")

func GeoCollect(ctx context.Context, data []byte) {
	var (
		l         = types.MustLoggerFromContext(ctx)
		project   = types.MustProjectFromContext(ctx)
		publisher = types.MustPublisherFromContext(ctx)
		eventID   = types.MustEventIDFromContext(ctx)

		dataStr = string(data)
		rawMap  = make(map[string]interface{})
	)

	// get lat or latitude key from data
	switch {
	case gjson.Get(dataStr, "lat").Exists():
		rawMap["latitude"] = gjson.Get(dataStr, "lat").Float()
	case gjson.Get(dataStr, "latitude").Exists():
		rawMap["latitude"] = gjson.Get(dataStr, "latitude").Float()
	default:
		rawMap["latitude"] = 0
	}

	// get long or longitude key from data
	switch {
	case gjson.Get(dataStr, "long").Exists():
		rawMap["longitude"] = gjson.Get(dataStr, "long").Float()
	case gjson.Get(dataStr, "longitude").Exists():
		rawMap["longitude"] = gjson.Get(dataStr, "long").Float()
	default:
		rawMap["longitude"] = 0
	}

	rawData, err := json.Marshal(rawMap)
	if err != nil {
		l.WithValues("eid", eventID).Error(err)
	}
	if err := autoCollectCli.Insert(fmt.Sprintf(`now(), '%s', '%s', '%s', 
		'%s'`, project.AccountID.String(), project.Name, publisher.Key, string(rawData))); err != nil {
		l.WithValues("eid", eventID).Error(err)
	}
}
