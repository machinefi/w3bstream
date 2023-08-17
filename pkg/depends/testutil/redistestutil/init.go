package redistestutil

import (
	"github.com/machinefi/w3bstream/pkg/depends/conf/app"
	conflogger "github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	"github.com/machinefi/w3bstream/pkg/depends/conf/redis"
)

var (
	Endpoint = &redis.Endpoint{}
	Redis    = &redis.Redis{}
)

func init() {
	app.New(
		app.WithName("test"),
		app.WithLogger(conflogger.Std()),
		app.WithRoot("."),
	).Conf(Redis, Endpoint)
}
