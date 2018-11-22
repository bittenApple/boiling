package boiling

import (
	"errors"
	"time"

	ecli "github.com/etcd-io/etcd/client"
)

const defaultKey = "BoilingKey"

type Options struct {
	Key       string   // the key name in etcd
	Step      int64    // the gap number batched once interactive with etcd
	Start     int64    // the start value of id sequence，it will be ignored if the value already set,
	Endpoints []string // etcd endpoints
}

type Client struct {
	opt  *Options
	ecli ecli.Client
}

// NewClient returns a client that generates id sequence.
func NewClient(o *Options) (*Client, error) {
	var err error
	cfg := ecli.Config{
		Endpoints:               o.Endpoints,
		Transport:               ecli.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	ecli, err := ecli.New(cfg)

	if err != nil {
		Logf("Init etcd client failed: %s", err)
		return nil, err
	}

	if o.Step < 0 {
		Logf("Step option should not be negative: %d", err)
		return nil, errors.New("error step option")
	}

	c := Client{
		opt:  o,
		ecli: ecli,
	}
	return &c, nil
}

func (c *Client) Next() int64 {
	return 0
}
