package index

import (
	"context"
	"database/sql"
)

type ChunkRepo struct {
	db *sql.DB
}

func NewChunk(db *sql.DB) *ChunkRepo {
	return &ChunkRepo{db: db}
}

type Chunk struct {
	FileID string
	Seq    int
	Size   int
	Sha256 string
	Ref    []byte
}

func (cr *ChunkRepo) Add(ctx context.Context, c Chunk) error {
	q := `insert or replace into chunks
	      (file_id, seq, size, sha256, ref)
	      values (?, ?, ?, ?, ?)`
	_, err := cr.db.ExecContext(ctx, q,
		c.FileID, c.Seq, c.Size, c.Sha256, c.Ref,
	)
	return err
}

func (cr *ChunkRepo) Get(ctx context.Context, fileID string) ([]Chunk, error) {
	q := `select * from chunks where file_id = ?`
	rows, err := cr.db.QueryContext(ctx, q, fileID)
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
