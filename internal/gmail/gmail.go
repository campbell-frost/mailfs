package gmail

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

const (
	poolSize         = 4
	defaultChunkSize = 17 << 20
	defaultMailbox   = "mailfs"
	defaultAddr      = "imap.gmail.com:993"
)

type Client struct {
	cfg    Config
	pool   chan *conn
	closed chan struct{}
	once   sync.Once
}

type Config struct {
	Username  string
	Password  string
	Mailbox   string
	Addr      string
	Chunksize int64
}

func (cfg *Config) defaults() {
	if cfg.Mailbox == "" {
		cfg.Mailbox = defaultMailbox
	}
	if cfg.Addr == "" {
		cfg.Addr = defaultAddr
	}
	if cfg.Chunksize <= 0 {
		cfg.Chunksize = defaultChunkSize
	}
}

func NewClient(cfg Config) (*Client, error) {
	cfg.defaults()
	if cfg.Username == "" || cfg.Password == "" {
		return nil, errors.New("gmail: username and password are required")
	}

	c := &Client{
		cfg:    cfg,
		pool:   make(chan *conn, poolSize),
		closed: make(chan struct{}),
	}

	// setup pool of empty conns
	for range poolSize {
		c.pool <- &conn{}
	}

	// grab one conn from the pool, and dial it to test credentials
	cn := <-c.pool
	defer func() { c.pool <- cn }()
	if err := c.dial(cn); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) Close() error {
	var errs []error
	c.once.Do(func() {
		log.Println("gmail: client closing, bye bye")
		close(c.closed)
		for range poolSize {
			cn := <-c.pool
			if cn.cl == nil {
				continue
			}

			if err := cn.cl.Logout().Wait(); err != nil {
				errs = append(errs, fmt.Errorf("logout: %w", err))
			}

			if err := cn.cl.Close(); err != nil {
				errs = append(errs, fmt.Errorf("close: %w", err))
			}
			cn.cl = nil
		}
	})
	return errors.Join(errs...)
}

func (c *Client) ChunkSize() int64 { return c.cfg.Chunksize }

func (*Client) Concurrency() int { return poolSize }
