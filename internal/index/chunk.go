package index

import "context"

type Chunk struct {
	FileID string
	Seq    int
	Size   int
	Sha256 string
	Ref    []byte
}

func (i *Index) AddChunk(ctx context.Context, c Chunk) error {
	q := `insert or replace into chunks
	      (file_id, seq, size, sha256, ref)
	      values (?, ?, ?, ?, ?)`
	_, err := i.db.ExecContext(ctx, q,
		c.FileID, c.Seq, c.Size, c.Sha256, c.Ref,
	)
	return err
}
