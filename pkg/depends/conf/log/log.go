package log

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
)

type Log struct {
	Name         string
	Level        Level            `env:""`
	Output       LoggerOutputType `env:""`
	Format       LoggerFormatType
	CKEndpoint   string             `env:""`
	Exporter     trace.SpanExporter `env:"-"`
	ReportCaller bool
}

func (l *Log) SetDefault() {
	if l.Level == 0 {
		l.Level = DebugLevel
	}
	if l.Output == 0 {
		l.Output = LOGGER_OUTPUT_TYPE__ALWAYS
	}
	if l.Format == 0 {
		l.Format = LOGGER_FORMAT_TYPE__JSON
	}
	if l.Name == "" {
		l.Name = "unknown"
		if v := os.Getenv(consts.EnvProjectName); v != "" {
			l.Name = v
		}
	}
}

func (l *Log) InitLogrus() {
	if l.Format == LOGGER_FORMAT_TYPE__JSON {
		logrus.SetFormatter(JsonFormatter)
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	}

	logrus.SetLevel(l.Level.LogrusLogLevel())
	logrus.SetReportCaller(l.ReportCaller)
	// TODO add hook with goid meta logrus.AddHook(goid.Default)
	logrus.AddHook(&ProjectAndMetaHook{l.Name})

	if l.CKEndpoint != "" {
		if ckHook, err := newClickhouseHook(l.CKEndpoint); err != nil {
			logrus.Errorf("new ck hook error: %s", err)
		} else {
			logrus.AddHook(ckHook)
		}
	}
	logrus.SetOutput(os.Stdout)
}

func (l *Log) InitSpanLog() {
	if l.Exporter == nil {
		return
	}
	if err := InstallPipeline(l.Output, l.Format, l.Exporter); err != nil {
		panic(err)
	}
}

func (l *Log) Init() {
	l.InitLogrus()
	l.InitSpanLog()
}

var ckBatchCount = 1000

type ClickhouseHook struct {
	CKEndpoint     string
	ckDB           *sql.DB
	insertSql      string
	entries        []logrus.Entry
	lastInsertTime time.Time
	duration       time.Duration
	lock           *sync.Mutex
}

func newClickhouseHook(endpoint string) (*ClickhouseHook, error) {
	clickhouseHook := &ClickhouseHook{CKEndpoint: endpoint, entries: make([]logrus.Entry, 0, ckBatchCount*2)}
	clickhouseHook.insertSql = `INSERT INTO w3b.server_logs (Timestamp, ProjectName, Msg) VALUES (?, ?, ?)`
	clickhouseHook.lastInsertTime = time.Now()
	clickhouseHook.duration = time.Second
	clickhouseHook.lock = &sync.Mutex{}
	db, err := clickhouseHook.connection()
	if err != nil {
		return nil, err
	}
	clickhouseHook.ckDB = db
	return clickhouseHook, nil
}

func (ck *ClickhouseHook) Fire(entry *logrus.Entry) error {
	if !ck.ping() {
		conn, err := ck.connection()
		if err != nil {
			logrus.Errorf("ck connection error: %s", err)
			return err
		}
		ck.ckDB = conn
	}

	ck.lock.Lock()
	defer ck.lock.Unlock()
	ck.entries = append(ck.entries, *entry)
	if len(ck.entries) >= ckBatchCount || time.Since(ck.lastInsertTime) > ck.duration {
		ck.insertLogs(ck.entries)
	}
	return nil
}

func (ck *ClickhouseHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (ck *ClickhouseHook) connection() (*sql.DB, error) {
	conn, err := sql.Open("clickhouse", ck.CKEndpoint)
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (ck *ClickhouseHook) ping() bool {
	if err := ck.ckDB.Ping(); err != nil {
		logrus.Errorf("ck ping error: %s", err)
		return false
	}
	return true
}

func (ck *ClickhouseHook) insertLogs(entries []logrus.Entry) error {
	err := ck.doWithTx(func(tx *sql.Tx) error {
		statement, err := tx.PrepareContext(context.Background(), ck.insertSql)
		if err != nil {
			return fmt.Errorf(" ck prepareContext error: %w", err)
		}
		defer func() {
			if err := statement.Close(); err != nil {
				logrus.Errorf("ck statement close error: %s", err)
			}
		}()
		for _, item := range entries {
			msg, _ := item.Bytes()
			if _, err := statement.ExecContext(context.Background(), item.Time, item.Data["@prj"], msg); err != nil {
				return fmt.Errorf("ck execContext error:%w", err)
			}
		}
		return nil
	})
	if err == nil {
		ck.entries = ck.entries[:0]
	}
	return err
}

func (ck *ClickhouseHook) doWithTx(fn func(tx *sql.Tx) error) error {
	tx, err := ck.ckDB.Begin()
	if err != nil {
		return fmt.Errorf("ck tx begin error: %w", err)
	}
	defer func() {
		tx.Rollback()
	}()
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}

type ProjectAndMetaHook struct {
	Name string
}

func (h *ProjectAndMetaHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		ctx = context.Background()
	}
	meta := metax.GetMetaFrom(ctx)
	if entry.Data["@prj"] == nil {
		entry.Data["@prj"] = h.Name
	}
	for k, v := range meta {
		entry.Data["meta."+k] = v
	}
	return nil
}

func (h *ProjectAndMetaHook) Levels() []logrus.Level { return logrus.AllLevels }

var (
	project       = "unknown"
	JsonFormatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "@lv",
			logrus.FieldKeyTime:  "@ts",
			logrus.FieldKeyFunc:  "@fn",
			logrus.FieldKeyFile:  "@fl",
		},
		CallerPrettyfier: func(f *runtime.Frame) (fn string, file string) {
			return f.Function + " line:" + strconv.FormatInt(int64(f.Line), 10), ""
		},
		TimestampFormat: "20060102-150405.000Z07:00",
	}
)

func init() {
	if v := os.Getenv(consts.EnvProjectName); v != "" {
		project = v
		if version := os.Getenv(consts.EnvProjectVersion); version != "" {
			project = project + "@" + version
		}
	}
	logrus.SetFormatter(JsonFormatter)
	logrus.SetReportCaller(true)
}
