package metrics

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/machinefi/w3bstream/pkg/types"
)

const queueLength = 5000

type ClickhouseClient struct {
	conn     driver.Conn
	sqLQueue chan string
	cfg      *config
}

type config struct {
	addr     string
	database string
	username string
	password string
}

var (
	clickhouseCLI *ClickhouseClient
	sleepTime     = 10 * time.Second
)

func Init(ctx context.Context) {
	cfg, existed := types.MetricsCenterConfigFromContext(ctx)
	if !existed || !validateCfg(cfg) {
		log.Println("fail to get the config of metrics center")
		return
	}
	clickhouseCLI = newClickhouseClient(&config{
		cfg.ClickHouseAddr,
		cfg.ClickHouseDB,
		cfg.ClickHouseUser,
		cfg.ClickHousePassword,
	})
}

func validateCfg(cfg *types.MetricsCenterConfig) bool {
	return len(cfg.ClickHouseAddr) > 0 &&
		len(cfg.ClickHouseDB) > 0 &&
		len(cfg.ClickHouseUser) > 0 &&
		len(cfg.ClickHousePassword) > 0
}

func newClickhouseClient(cfg *config) *ClickhouseClient {
	cc := &ClickhouseClient{
		sqLQueue: make(chan string),
		cfg:      cfg,
	}
	go cc.Run()
	return cc
}

func (c *ClickhouseClient) Insert(query string) error {
	if len(c.sqLQueue) > queueLength {
		return errors.New("the queue of client is full")
	}
	c.sqLQueue <- query
	return nil
}

func (c *ClickhouseClient) Run() {
	for {
		if c.conn == nil {
			if err := c.connect(); err != nil {
				log.Println("ClickhouseClient failed to connect: ", err)
				time.Sleep(sleepTime)
				continue
			}
		}
		query := <-c.sqLQueue
		if err := c.conn.AsyncInsert(context.Background(), query, true); err != nil {
			if !c.liveness() {
				c.conn = nil
				c.sqLQueue <- query
				continue
			}
			log.Println("ClickhouseClient failed to insert data: ", err)
		}
	}
}

func (c *ClickhouseClient) connect() error {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{c.cfg.addr},
		Auth: clickhouse.Auth{
			Database: c.cfg.database,
			Username: c.cfg.username,
			Password: c.cfg.password,
		},
	})
	if err != nil {
		return err
	}
	c.conn = conn
	if !c.liveness() {
		c.conn = nil
		return errors.New("failed to ping clickhouse server")
	}
	log.Println("clickhouse server login successfully")
	return nil
}

func (c *ClickhouseClient) liveness() bool {
	if err := c.conn.Ping(context.Background()); err != nil {
		log.Println("failed to ping clickhouse server: ", err)
		return false
	}
	return true
}
