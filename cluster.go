package grmln

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// Cluster config defaults
const (
	DefaultBackoffBase = time.Millisecond * 100
	DefaultBackoffMax  = time.Second * 10
	DefaultMimeType    = "application/json"
)

type (
	// OnConnectError is called when there is an error connecting
	OnConnectError func(addr string, err error, attempts int)
)

// Cluster represents a cluster of Gremlin servers
type Cluster struct {
	conns       chan *Conn
	backoffBase time.Duration
	backoffMax  time.Duration

	onConnectError OnConnectError

	closing chan struct{}
}

// ClusterConfig contains configuration options for the cluster
type ClusterConfig struct {
	// MimeType is the mime type to use when sending requests. Defaults to application/json
	MimeType string

	// ConnectionsPerAddress is the number of connections to open to each address. Defaults to 1
	ConnectionsPerAddress int

	// DefaultScriptEvaluationTimeout is the default script evaluation timeout. Defaults to 3000ms
	DefaultScriptEvaluationTimeout time.Duration

	// DefaultEvalLanguage is the default eval language for use. Defaults to "gremlin-groovy"
	DefaultEvalLanguage string

	// DefaultBatchSize is the default size of batched responses. 0 uses server default
	DefaultBatchSize int

	// BackoffBase is the base used when calculating connection retry backoff
	BackoffBase time.Duration

	// BackoffMax is the maximum time we can wait between connect attempts
	BackoffMax time.Duration

	// OnConnectError is called when there is an error connecting
	OnConnectError OnConnectError
}

// NewCluster creates a new cluster
func NewCluster(config ClusterConfig, addrs ...string) *Cluster {
	config = setDefaults(config)

	cluster := Cluster{
		conns:          make(chan *Conn, len(addrs)*config.ConnectionsPerAddress),
		backoffBase:    config.BackoffBase,
		backoffMax:     config.BackoffMax,
		onConnectError: config.OnConnectError,
		closing:        make(chan struct{}),
	}

	for _, addr := range addrs {
		for i := 0; i < config.ConnectionsPerAddress; i++ {
			go cluster.putConn(cluster.connect(addr, config.MimeType), nil)
		}
	}

	return &cluster
}

func setDefaults(config ClusterConfig) ClusterConfig {
	if config.BackoffBase == 0 {
		config.BackoffBase = DefaultBackoffBase
	}

	if config.BackoffMax == 0 {
		config.BackoffMax = DefaultBackoffMax
	}

	if config.OnConnectError == nil {
		config.OnConnectError = func(addr string, err error, attempts int) {}
	}

	if config.ConnectionsPerAddress == 0 {
		config.ConnectionsPerAddress = 1
	}

	if config.MimeType == "" {
		config.MimeType = DefaultMimeType
	}

	return config
}

func (c *Cluster) getConn(ctx context.Context) (*Conn, error) {
	select {
	case <-c.closing:
		return nil, clusterErrorClusterClosed
	default:
	}

	select {
	case conn, ok := <-c.conns:
		if !ok {
			return nil, clusterErrorClusterClosed
		}
		return conn, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Cluster) putConn(conn *Conn, err error) {
	if err != nil {
		go func() {
			conn := c.connect(conn.addr, conn.mimeType)
			if conn != nil {
				c.conns <- conn
			}
		}()
		return
	}

	c.conns <- conn
}

func (c *Cluster) randomJitter(sleep time.Duration) time.Duration {
	base := int64(c.backoffBase)
	return time.Duration(
		math.Min(
			float64(c.backoffMax),
			float64(rand.Int63n(int64(sleep)-base)+base),
		),
	)
}

func (c *Cluster) connect(addr, mimeType string) *Conn {
	sleep := c.backoffBase
	attempts := 1
	for {
		conn, err := Dial(context.Background(), addr, mimeType)
		if err == nil {
			// connected!
			return conn
		}

		c.onConnectError(addr, err, attempts)

		select {
		case <-c.closing:
			return nil
		case <-time.After(sleep):
		}
		sleep = c.randomJitter(sleep * 3)
		attempts++
	}
}

// ProcessRequest can process a raw gremlin request
func (c *Cluster) ProcessRequest(ctx context.Context, r Request, onResponse OnResponse) error {
	conn, err := c.getConn(ctx)
	if err != nil {
		return err
	}

	err = conn.sendRequest(ctx, r)
	if err != nil {
		c.putConn(conn, err)
		return err
	}

	err = conn.readResponse(ctx, onResponse)
	if err != nil {
		c.putConn(conn, err)
		return err
	}

	c.putConn(conn, nil)
	return nil
}

// Close closes the cluster
func (c *Cluster) Close() error {
	close(c.closing) // prevent further gets outside of close function

	close(c.conns)
	for conn := range c.conns {
		conn.Close()
	}
	return nil
}
