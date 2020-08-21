package nrredigo

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// Pool is an interface for representing a pool of Redis connections
type Pool interface {
	GetContext(ctx context.Context) (redis.Conn, error)
	Get() redis.Conn
	Close() error
}

// Wrap returns a wrappedPool that can be used like a normal redis Pool, but sends segments to new relic
func Wrap(p Pool, opts ...Option) Pool {
	return &wrappedPool{
		Pool: p,
		cfg:  createConfig(opts),
	}
}

type wrappedPool struct {
	Pool
	cfg *Config
}

func (p *wrappedPool) GetContext(ctx context.Context) (conn redis.Conn, err error) {
	conn, err = p.Pool.GetContext(ctx)
	if err != nil {
		return
	}

	nrtx := newrelic.FromContext(ctx)
	if nrtx != nil {
		return wrapConn(conn, nrtx, p.cfg), nil
	}

	return
}

func (p *wrappedPool) Get() redis.Conn {
	return p.Pool.Get()
}

func (p *wrappedPool) Close() error {
	return p.Pool.Close()
}
