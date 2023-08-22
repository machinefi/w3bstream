package global

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/pub_bench/types"
	confapp "github.com/machinefi/w3bstream/pkg/depends/conf/app"
	conflogger "github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
)

type config struct {
	Broker   *confmqtt.Broker `env:""`
	Channels []*types.Channel `env:""`
}

var (
	// App command context
	App *confapp.Ctx

	// cfg command global config
	cfg *config

	// WithContext global with context func
	WithContext contextx.WithContext
	// Context global context with empty
	Context context.Context
)

func init() {
	cfg = &config{
		Broker: &confmqtt.Broker{},
	}

	App = confapp.New(
		confapp.WithName("pub_bench"),
		confapp.WithRoot(".."),
		confapp.WithLogger(conflogger.Std()),
	)
	App.Conf(cfg)

	WithContext = contextx.WithContextCompose(
		types.WithMqttBrokerContext(cfg.Broker),
		types.WithChannelsContext(cfg.Channels),
	)
	Context = WithContext(context.Background())
}
