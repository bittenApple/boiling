package boiling

import (
	"context"
	"errors"
	"time"

	ecli "github.com/etcd-io/etcd/clientv3"
)

const (
	defaultKey = "BoilingKey"
	defaultVal = "BoilingValue"
)

type Options struct {
	Key       string   // the key name in etcd
	Buffer    int64    // the buffered ids interactive once with etcd
	Start     int64    // the start value of id sequences，which will be ignored if the value already be set,
	Endpoints []string // etcd endpoints
	Initial   bool     // Whether restart from start value
}

type Client struct {
	options *Options
	ecli    *ecli.Client
	ch      chan int64
}

// NewClient returns a client which generates incremented id.
func NewClient(opt *Options) (*Client, error) {
	var err error
	cfg := ecli.Config{
		Endpoints:   opt.Endpoints,
		DialTimeout: time.Second,
	}

	ecli, err := ecli.New(cfg)
	if err != nil {
		LogErrf("init etcd client failed: %s", err)
		return nil, err
	}

	if opt.Buffer < 0 {
		LogErrf("buffer option should not be negative: %d", err)
		return nil, errors.New("error step option")
	}

	if opt.Buffer == 0 {
		opt.Buffer = 1
	}

	if opt.Key == "" {
		opt.Key = defaultKey
	}

	c := Client{
		options: opt,
		ecli:    ecli,
		ch:      make(chan int64, opt.Buffer),
	}

	if c.options.Initial {
		err = c.Reset()
		if err != nil {
			LogErrf("reset etcd key failed: %s", err)
			return nil, err
		}
	}

	go c.run()
	return &c, nil
}

func (c *Client) Exit() {
	c.ecli.Close()
}

func (c *Client) run() {
	for {
		resp, err := c.ecli.Put(context.Background(), c.options.Key, defaultVal, ecli.WithPrevKV())
		if err != nil {
			LogErrf("put request error: %s", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		var ver, i int64
		if resp.PrevKv == nil { // Key not exited, first put
			ver = 0
		} else {
			ver = resp.PrevKv.Version
		}

		start := ver*c.options.Buffer + c.options.Start
		for i = 0; i < c.options.Buffer; i++ {
			c.ch <- i + start
		}
	}
}

// GetId returns the generated id in sequence
func (c *Client) GetId() int64 {
	return <-c.ch
}

// Reset deletes the corresponding etcd key, resets the id to start value
func (c *Client) Reset() error {
	_, err := c.ecli.Delete(context.Background(), c.options.Key)
	return err
}
