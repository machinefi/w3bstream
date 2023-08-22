package metrics

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	"github.com/machinefi/w3bstream/pkg/types"
)

const (
	queueLength      = 5000
	popThreshold     = 3
	concurrentWorker = 10
)

type (
	ClickhouseClient struct {
		workerPool []*connWorker
		sqLQueue   chan *queueElement
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
	{
		opts.Settings["async_insert"] = 1
		opts.Settings["wait_for_async_insert"] = 0
		opts.Settings["async_insert_busy_timeout_ms"] = 100
	}
	clickhouseCLI = newClickhouseClient(opts)
	log.Println("clickhouse client is initialized")
}

func newClickhouseClient(cfg *clickhouse.Options) *ClickhouseClient {
	cc := &ClickhouseClient{
		sqLQueue: make(chan *queueElement, queueLength),
	}
	for i := 0; i < concurrentWorker; i++ {
		cc.workerPool = append(cc.workerPool, &connWorker{
			sqLQueue: cc.sqLQueue,
			cfg:      cfg,
		})
		go cc.workerPool[i].run()
	}
	return cc
}

func (c *ClickhouseClient) Insert(query string) error {
	select {
	case c.sqLQueue <- &queueElement{
		query: query,
		count: 0,
	}:
	default:
		return errors.New("the queue of client is full")
	}
	return nil
}

type SQLBatcher struct {
	signal   chan string
	preStatm string
	buf      []string
}

const (
	batchSize      = 50000
	tickerInterval = 200 * time.Millisecond
)

func NewSQLBatcher(preStatm string) *SQLBatcher {
	bw := &SQLBatcher{
		signal:   make(chan string, queueLength),
		preStatm: preStatm,
		buf:      make([]string, 0, batchSize),
	}
	go bw.run()
	return bw
}

func (b *SQLBatcher) Insert(query string) error {
	_, l := logger.NewSpanContext(context.Background(), "modules.metrics.SQLBatcher.Insert")
	defer l.End()

	if clickhouseCLI == nil {
		return errors.New("clickhouse client is not initialized")
	}
	select {
	case b.signal <- query:
		return nil
	default:
		return errors.New("the queue of SQLBatcher is full")
	}
}

func (b *SQLBatcher) run() {
	ticker := time.NewTicker(tickerInterval)
	for {
		select {
		case <-ticker.C:
			if len(b.buf) == 0 {
				continue
			}
			if clickhouseCLI == nil {
				log.Println("clickhouse client is not initialized")
				continue
			}
			_ = b.insert()
		case str, ok := <-b.signal:
			if !ok {
				return
			}
			if clickhouseCLI == nil {
				log.Println("clickhouse client is not initialized")
				continue
			}
			b.buf = append(b.buf, str)
			if len(b.buf) >= batchSize {
				if b.insert() != nil {
					continue
				}
				ticker.Reset(tickerInterval)
			}
		}
	}
}

func (b *SQLBatcher) insert() error {
	_, l := logger.NewSpanContext(context.Background(), "modules.metrics.SQLBatcher.insert")
	defer l.End()

	err := clickhouseCLI.Insert(b.preStatm + "(" + strings.Join(b.buf, "),(") + ")")
	if err != nil {
		l.Error(err)
		return err
	}
	b.buf = b.buf[0:0]
	return nil
}

type connWorker struct {
	sqLQueue chan *queueElement
	conn     driver.Conn
	cfg      *clickhouse.Options
}

func (c *connWorker) run() {
	for {
		if err := c.connect(); err != nil {
			log.Println("ClickhouseClient failed to connect: ", err)
			time.Sleep(sleepTime)
			continue
		}
		ele := <-c.sqLQueue
		if err := c.conn.Exec(context.Background(), ele.query); err != nil {
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

func (c *connWorker) connect() error {
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

func (c *connWorker) liveness() bool {
	if err := c.conn.Ping(context.Background()); err != nil {
		log.Println("failed to ping clickhouse server: ", err)
		return false
	}
	return true
}
