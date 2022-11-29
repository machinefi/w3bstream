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

type Endpoint struct {
	Endpoint types.Endpoint `env:""`
	Wait     bool
	Prefix   string
	pool     *redis.Pool
}

func (r *Endpoint) Get() redis.Conn {
	if r.pool != nil {
		return r.pool.Get()
	}
	return nil
}

func (r *Endpoint) Key(key string) string {
	return fmt.Sprintf("%s:%s", r.Prefix, key)
}

func (r *Endpoint) LivenessCheck() map[string]string {
	m := map[string]string{}

	conn := r.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		m[r.Endpoint.Host()] = err.Error()
	} else {
		m[r.Endpoint.Host()] = "ok"
	}

	return m
}

func (r *Endpoint) SetDefault() {
	if r.Endpoint.Scheme == "" {
		r.Endpoint.Scheme = "tcp"
	}
	if r.Endpoint.Hostname == "" {
		r.Endpoint.Hostname = "127.0.0.1"
	}
	if r.Endpoint.Port == 0 {
		r.Endpoint.Port = 6379
	}
	if !r.Wait {
		r.Wait = true
	}
	if r.Prefix == "" {
		r.Prefix = fmt.Sprintf("%s:%s:",
			strings.ToLower(os.Getenv(consts.GoRuntimeEnv)),
			strings.ToLower(os.Getenv(consts.EnvProjectName)),
		)
	}
}

func (r *Endpoint) Init() {
	if r.pool == nil {
		r.init()
	}
}

func (r *Endpoint) init() {
	opt := struct {
		ConnectTimeout types.Duration `name:"connectTimeout" default:"10s"`
		ReadTimeout    types.Duration `name:"readTimeout"    default:"10s"`
		WriteTimeout   types.Duration `name:"writeTimeout"   default:"10s"`
		IdleTimeout    types.Duration `name:"idleTimeout"    default:"240s"`
		MaxActive      int            `name:"maxActive"      default:"5"`
		MaxIdle        int            `name:"maxIdle"        default:"3"`
		DB             int            `name:"dB"             default:"1"`
	}{}

	err := types.UnmarshalExtra(r.Endpoint.Param, &opt)
	if err != nil {
		panic(err)
	}

	dialer := func() (c redis.Conn, err error) {
		options := []redis.DialOption{
			redis.DialDatabase(opt.DB),
			redis.DialConnectTimeout(time.Duration(opt.ConnectTimeout)),
			redis.DialWriteTimeout(time.Duration(opt.WriteTimeout)),
			redis.DialReadTimeout(time.Duration(opt.ReadTimeout)),
		}

		if r.Endpoint.Password != "" {
			options = append(options, redis.DialPassword(r.Endpoint.Password.String()))
		}

		return redis.Dial("tcp", r.Endpoint.Host(), options...)
	}

	r.pool = &redis.Pool{
		Dial:        dialer,
		MaxIdle:     opt.MaxIdle,
		MaxActive:   opt.MaxActive,
		IdleTimeout: time.Duration(opt.IdleTimeout),
		Wait:        true,
	}
}

func (r *Endpoint) Exec(cmd *Cmd, others ...*Cmd) (interface{}, error) {
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
