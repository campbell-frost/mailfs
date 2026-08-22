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

func (i *Index) Chunks(ctx context.Context, fileID string) ([]Chunk, error) {
	q := `select * from chunks where file_id = ?`
	rows, err := i.db.QueryContext(ctx, q, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chunks []Chunk
	for rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		var c Chunk
		if err := rows.Scan(&c.FileID, &c.Seq, &c.Size, &c.Sha256, &c.Ref); err != nil {
			return nil, err
		}
		chunks = append(chunks, c)
	}
	return chunks, nil
}
