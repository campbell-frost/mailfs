package gmail

import (
	"context"
	"errors"
	"fmt"

	"github.com/emersion/go-imap/v2"
)

func (c *Client) Delete(ctx context.Context, chunkRefs [][]byte) error {
	if len(chunkRefs) == 0 {
		return nil
	}

	refs, err := parseRefs(chunkRefs)
	if err != nil {
		return fmt.Errorf("failed to parse refs: %w", err)
	}
	return c.do(ctx, func(cn *conn) error {
		var uids imap.UIDSet
		for _, r := range refs {
			set, err := c.locate(cn, r)
			if errors.Is(err, errNotFound) {
				continue
			}
			if err != nil {
				return fmt.Errorf("failed to locate ref: %w", err)
			}
			uids.AddSet(set)
		}

		if len(uids) == 0 {
			return nil
		}

		_, err := cn.cl.Move(uids, trash).Wait()
		return err
	})
}
