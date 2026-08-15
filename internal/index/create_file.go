package index

import (
	"context"
	"time"
)

type File struct {
	ID        string
	Filename  string
	Size      int64
	CreatedAt time.Time
	Status    string
	TempPath  string
	Chunks    int
}

func (i *Index) CreateFile(ctx context.Context, f File) error {
	q := `insert into files (id, filename, temp_path, size, status, created_at) values (?, ?, ?, ?, ?, ?)`
	_, err := i.db.ExecContext(ctx, q, f.ID, f.Filename, f.TempPath, f.Size, f.Status, f.CreatedAt.UTC().Format(time.RFC3339Nano))
	return err
}
