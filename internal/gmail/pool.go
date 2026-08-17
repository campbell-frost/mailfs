package gmail

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

type conn struct {
	cl          *imapclient.Client
	uidValidity uint32
}

var ErrClosed = errors.New("gmail: client closed")

// do runs fn on a pooled IMAP connection.  It blocks until one is free.
func (c *Client) do(ctx context.Context, fn func(*conn) error) error {
	var cn *conn
	select {
	case cn = <-c.pool:
	case <-c.closed:
		return ErrClosed
	case <-ctx.Done():
		return ctx.Err()
	}

	defer func() { c.pool <- cn }()

	for attempt := range 2 {
		if cn.cl == nil {
			if err := c.dial(cn); err != nil {
				return err
			}
		}

		err := fn(cn)
		if err == nil {
			return nil
		}

		if _, ok := errors.AsType[*imap.Error](err); ok {
			return err
		}

		cn.cl.Close()
		cn.cl = nil
		if attempt == 1 {
			return err
		}

	}
	return nil
}

// dial dials TLS, logs in, creates the mailbox if it does not exist, selects it,
// and stores the UIDValidity on the conn.
func (c *Client) dial(cn *conn) error {
	cl, err := imapclient.DialTLS(c.cfg.Addr, nil)
	if err != nil {
		return fmt.Errorf("dial %v: %w", c.cfg.Addr, err)
	}

	if err := cl.Login(c.cfg.Username, c.cfg.Password).Wait(); err != nil {
		cl.Close()
		return fmt.Errorf("login %v: %w", c.cfg.Username, err)
	}

	if err = cl.Create(c.cfg.Mailbox, nil).Wait(); err != nil && !alreadyExists(err) {
		cl.Close()
		return fmt.Errorf("create %v: %w", c.cfg.Mailbox, err)
	}

	sel, err := cl.Select(c.cfg.Mailbox, nil).Wait()
	if err != nil {
		cl.Close()
		return fmt.Errorf("select %v: %w", c.cfg.Mailbox, err)
	}

	cn.cl = cl
	cn.uidValidity = sel.UIDValidity

	return nil
}

func alreadyExists(err error) bool {
	var status *imap.Error
	if errors.As(err, &status) && status.Code == imap.ResponseCodeAlreadyExists {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "already exists")
}
