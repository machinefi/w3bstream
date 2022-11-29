package redistestutil

import (
	"github.com/machinefi/w3bstream/pkg/depends/conf/app"
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/conf/redis"
)

var (
	Endpoint = &redis.Endpoint{}
	Redis    = &redis.Redis{}
)

func init() {
	app.New(
		app.WithName("test"),
		app.WithLogger(log.Std()),
		app.WithRoot("."),
	).Conf(Redis, Endpoint)
}
