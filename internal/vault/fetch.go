package vault

import (
	"context"
	"io"

	"golang.org/x/sync/errgroup"
)

func (v *Vault) Fetch(ctx context.Context, id string, w io.Writer) (error) {
	const readahead = 4

	chunks, err := v.idx.Chunks(ctx, id)
	if err != nil {
		return err
	}
	for start := 0; start < len(chunks); start += readahead {
		batch := chunks[start:min(start+readahead, len(chunks))]
		bufs := make([][]byte, len(batch))
		g, gCtx := errgroup.WithContext(ctx)
		for i, c := range batch {
			g.Go(func() error {
				data, err := v.gmail.Get(gCtx, c.Ref)
				if err != nil {
					return err
				}

				bufs[i] = data
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

		for _, b := range bufs {
			if _, err := w.Write(b); err != nil {
				return err
			}
		}
	}

	return nil
}
