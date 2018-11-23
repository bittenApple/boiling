// Package boiling provides function to generate incremented id based on etcd
package boiling

import (
	"context"
	"errors"
	"time"

	ecli "github.com/etcd-io/etcd/clientv3"
	"google.golang.org/grpc"
)

const (
	defaultKey = "BoilingKey"
	defaultVal = "BoilingValue"
)

// Options configures how we set up the client.
type Options struct {
	Key       string   // the key name in etcd
	Buffer    int64    // the buffered ids interactive once with etcd
	Start     int64    // the start value of id sequences，which will be ignored if the value already be set,
	Endpoints []string // etcd endpoints
	Initial   bool     // Whether restart from start value
}

// Client provides and manages an boiling client
type Client struct {
	options *Options
	ecli    *ecli.Client
	ch      chan int64

	ctx    context.Context
	cancel context.CancelFunc
}

// NewClient returns a client which generates incremented id.
func NewClient(opt *Options) (*Client, error) {
	var err error
	cfg := ecli.Config{
		Endpoints:   opt.Endpoints,
		DialTimeout: time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	}

	ecli, err := ecli.New(cfg)
	if err != nil {
		logErrf("init etcd client failed: %s", err)
		return nil, err
	}

	if opt.Buffer < 0 {
		logErrf("buffer option should not be negative: %d", err)
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
			logErrf("reset etcd key failed: %s", err)
			return nil, err
		}
	}

	c.ctx, c.cancel = context.WithCancel(context.Background())
	go c.run()
	return &c, nil
}

// Stop stops generating id, and getId() function will always return 0
func (c *Client) Stop() {
	c.cancel()
}

func (c *Client) run() {
	defer func() {
		c.ecli.Close()
		close(c.ch)
	}()
	for {
		resp, err := c.ecli.Put(context.Background(), c.options.Key, defaultVal, ecli.WithPrevKV())
		if err != nil {
			logErrf("put request error: %s", err)
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
			select {
			case <-c.ctx.Done():
				logErrf("boiling was stopped")
				return
			default:
				c.ch <- i + start
			}
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
