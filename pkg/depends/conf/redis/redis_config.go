package redis

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
)

type Redis struct {
	Protocol       string
	Host           string
	Port           int
	Password       types.Password
	ConnectTimeout types.Duration
	WriteTimeout   types.Duration
	ReadTimeout    types.Duration
	IdleTimeout    types.Duration
	MaxActive      int
	MaxIdle        int
	Wait           bool
	DB             int
	pool           *redis.Pool
	Prefix         string
}

func (r *Redis) Get() redis.Conn {
	if r.pool != nil {
		return r.pool.Get()
	}
	return nil
}

func (r *Redis) clone() *Redis {
	cloned := *r
	return &cloned
}

func (r *Redis) WithPrefix(prefix string) *Redis {
	cloned := r.clone()
	cloned.Prefix += "::" + prefix
	return cloned
}

func (r *Redis) WithDBIndex(n int) *Redis {
	cloned := r.clone()
	cloned.DB = n
	return cloned
}

func (r *Redis) Exec(cmd *Cmd, others ...*Cmd) (interface{}, error) {
	c := r.Get()
	defer c.Close()

	if (len(others)) == 0 {
		return c.Do(cmd.Name, cmd.Args...)
	}

	err := c.Send("MULTI")
	if err != nil {
		return nil, err
	}

	err = c.Send(cmd.Name, cmd.Args...)
	if err != nil {
		return nil, err
	}

	for i := range others {
		o := others[i]
		if o == nil {
			continue
		}
		err := c.Send(o.Name, o.Args...)
		if err != nil {
			return nil, err
		}
	}

	return c.Do("EXEC")
}

func (r *Redis) Key(key string) string {
	return fmt.Sprintf("%s:%s", r.Prefix, key)
}

func (r *Redis) LivenessCheck() map[string]string {
	m := map[string]string{}

	conn := r.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	if err != nil {
		m[r.Host] = err.Error()
	} else {
		m[r.Host] = "ok"
	}

	return m
}

func (r *Redis) SetDefault() {
	if r.Protocol == "" {
		r.Protocol = "tcp"
	}
	if r.Host == "" {
		r.Host = "127.0.0.1"
	}
	if r.Port == 0 {
		r.Port = 6379
	}
	if r.ConnectTimeout == 0 {
		r.ConnectTimeout = types.Duration(10 * time.Second)
	}
	if r.ReadTimeout == 0 {
		r.ReadTimeout = types.Duration(10 * time.Second)
	}
	if r.WriteTimeout == 0 {
		r.WriteTimeout = types.Duration(10 * time.Second)
	}
	if r.IdleTimeout == 0 {
		r.IdleTimeout = types.Duration(240 * time.Second)
	}
	if r.MaxActive == 0 {
		r.MaxActive = 5
	}
	if r.MaxIdle == 0 {
		r.MaxIdle = 3
	}
	if !r.Wait {
		r.Wait = true
	}
	if r.DB == 0 {
		r.DB = 1
	}
	if r.Prefix == "" {
		r.Prefix = fmt.Sprintf("%s:%s:",
			strings.ToLower(os.Getenv(consts.GoRuntimeEnv)),
			strings.ToLower(os.Getenv(consts.EnvProjectName)),
		)
	}
}

func (r *Redis) Init() {
	if r.pool == nil {
		r.init()
	}
}

func (r *Redis) init() {
	dialer := func() (c redis.Conn, err error) {
		c, err = redis.Dial(
			r.Protocol,
			fmt.Sprintf("%s:%d", r.Host, r.Port),

			redis.DialWriteTimeout(time.Duration(r.WriteTimeout)),
			redis.DialConnectTimeout(time.Duration(r.ConnectTimeout)),
			redis.DialReadTimeout(time.Duration(r.ReadTimeout)),
			redis.DialPassword(r.Password.String()),
			redis.DialDatabase(r.DB),
		)
		return
	}

	r.pool = &redis.Pool{
		Dial:        dialer,
		MaxIdle:     r.MaxIdle,
		MaxActive:   r.MaxActive,
		IdleTimeout: time.Duration(r.IdleTimeout),
		Wait:        true,
	}
}
