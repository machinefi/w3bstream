package metrics

import (
	"context"
	"fmt"

	"github.com/tidwall/gjson"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/models"
)

var autoCollectCli = NewSQLBatcher("INSERT INTO ws_metrics.auto_collect_metrics VALUES")

var (
	lonkeys = []string{"lon", "Lon", "long", "Long", "longitude", "Longitude"}
	latkeys = []string{"lat", "latitude", "Lat", "Latitude"}
)

func GeoCollect(ctx context.Context, v *models.Event) {
	ctx, l := logr.Start(ctx, "metrics.GeoCollect")
	defer l.End()

	var (
		lon    float64
		lat    float64
		haslon = false
		haslat = false
		data   = string(v.Input)
	)

	for _, key := range lonkeys {
		if r := gjson.Get(data, key); r.Exists() {
			lon = r.Float()
			haslon = true
			break
		}
	}
	if !haslon {
		return
	}

	for _, key := range latkeys {
		if r := gjson.Get(data, key); r.Exists() {
			lat = r.Float()
			haslat = true
			break
		}
	}
	if !haslat {
		return
	}

	geodata := fmt.Sprintf(`{"longitude": %f, "latitude": %f}`, lon, lat)
	if err := autoCollectCli.Insert(fmt.Sprintf(`now(), '%s', '%s', '%s', '%s'`, v.AccountID.String(), v.ProjectName, v.PublisherKey, geodata)); err != nil {
		l.WithValues("eid", v.EventID).Error(err)
	}
}
