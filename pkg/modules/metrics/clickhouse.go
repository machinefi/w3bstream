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

const (
	queueLength  = 5000
	popThreshold = 3
)

type (
	ClickhouseClient struct {
		conn     driver.Conn
		sqLQueue chan *queueElement
		cfg      *clickhouse.Options
	}

	queueElement struct {
		query string
		count int
	}
)

var (
	clickhouseCLI *ClickhouseClient
	sleepTime     = 10 * time.Second
)

func Init(ctx context.Context) {
	cfg, existed := types.MetricsCenterConfigFromContext(ctx)
	if !existed || len(cfg.ClickHouseDSN) == 0 {
		log.Println("fail to get the config of metrics center")
		return
	}
	opts, err := clickhouse.ParseDSN(cfg.ClickHouseDSN)
	if err != nil {
		panic(err)
	}
	clickhouseCLI = newClickhouseClient(opts)
}

func newClickhouseClient(cfg *clickhouse.Options) *ClickhouseClient {
	cc := &ClickhouseClient{
		sqLQueue: make(chan *queueElement),
		cfg:      cfg,
	}
	go cc.Run()
	return cc
}

func (c *ClickhouseClient) Insert(query string) error {
	if len(c.sqLQueue) > queueLength {
		return errors.New("the queue of client is full")
	}
	c.sqLQueue <- &queueElement{
		query: query,
		count: 0,
	}
	return nil
}

func (c *ClickhouseClient) Run() {
	for {
		if err := c.connect(); err != nil {
			log.Println("ClickhouseClient failed to connect: ", err)
			time.Sleep(sleepTime)
			continue
		}
		ele := <-c.sqLQueue
		if err := c.conn.AsyncInsert(context.Background(), ele.query, true); err != nil {
			if !c.liveness() {
				c.conn = nil
				log.Printf("ClickhouseClient failed to connect the server: error: %s, query %s\n", err, ele.query)
			} else {
				log.Printf("ClickhouseClient failed to insert data: error %s, query %s\n ", err, ele.query)
			}
			if ele.count > popThreshold {
				log.Printf("the query %s in ClickhouseClient is poped due to %d times failure.", ele.query, ele.count)
				continue
			}
			ele.count++
			// TODO: Double linked list should be used to append the element to the head
			// when the order of the queue is important
			c.sqLQueue <- ele
		}
	}
}

func (c *ClickhouseClient) connect() error {
	if c.conn != nil {
		return nil
	}
	conn, err := clickhouse.Open(c.cfg)
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
