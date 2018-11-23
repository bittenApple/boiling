package boiling

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestKey = "TestBoilingKey"

var endpoints []string

func init() {
	env := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}
	endpoints = strings.Split(env("ETCD_ENDPOINTS", ""), ",")
}

func TestGetId(t *testing.T) {
	o := &Options{
		Endpoints: endpoints,
		Key:       TestKey,
		Buffer:    10,
		Initial:   true,
	}

	cli, err := NewClient(o)
	assert.Nil(t, err)

	var i int64
	for i = 0; i < 100; i++ {
		assert.Equal(t, i, cli.GetId())
	}
}

func TestStop(t *testing.T) {
	o := &Options{
		Endpoints: endpoints,
		Key:       TestKey,
		Initial:   true,
	}

	cli, err := NewClient(o)
	assert.Nil(t, err)
	cli.Stop()
	for i := 0; i < 10; i++ {
		assert.Equal(t, int64(0), cli.GetId())
	}
}
