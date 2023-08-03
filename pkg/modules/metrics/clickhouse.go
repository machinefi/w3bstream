package metrics

import (
	"container/list"
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/machinefi/w3bstream/pkg/types"
)

const (
	queueLength      = 5000
	popThreshold     = 3
	concurrentWorker = 1
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

type BatchWorker struct {
	signal   chan *list.List
	preStatm string
	li       *list.List
	mtx      sync.Mutex
}

const (
	batchSize = 50000
)

func NewBatchWorker(preStatm string) *BatchWorker {
	bw := &BatchWorker{
		signal:   make(chan *list.List, queueLength),
		preStatm: preStatm,
		li:       list.New(),
	}
	go bw.run()
	return bw
}

func (b *BatchWorker) Insert(query string) error {
	if clickhouseCLI == nil {
		return errors.New("clickhouse client is not initialized")
	}
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.li.PushBack(query)
	if b.li.Len() > batchSize {
		select {
		case b.signal <- b.li:
			b.li = list.New()
		default:
			return errors.New("batchWorker queue is full")
		}
	}
	return nil
}

func (b *BatchWorker) run() {
	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			b.mtx.Lock()
			li := b.li
			b.li = list.New()
			b.mtx.Unlock()
			arr := pack(li)
			if len(arr) == 0 {
				continue
			}
			if clickhouseCLI == nil {
				log.Println("clickhouse client is not initialized")
				continue
			}
			err := clickhouseCLI.Insert(b.preStatm + "(" + strings.Join(arr, "),(") + ")")
			if err != nil {
				log.Println("batchWorker failed to insert: ", err)
			}
		case li := <-b.signal:
			arr := pack(li)
			if clickhouseCLI == nil {
				log.Println("clickhouse client is not initialized")
				continue
			}
			err := clickhouseCLI.Insert(b.preStatm + "(" + strings.Join(arr, "),(") + ")")
			if err != nil {
				log.Println("batchWorker failed to insert: ", err)
			}
		}
	}
}

func pack(li *list.List) []string {
	arr := make([]string, 0, li.Len())
	for element := li.Front(); element != nil; element = element.Next() {
		arr = append(arr, element.Value.(string))
	}
	return arr
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
