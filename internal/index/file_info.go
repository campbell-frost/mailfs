package index

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type FileInfo struct {
	ID        string
	Filename  string
	Size      int64
	CreatedAt time.Time
	Status    string
	TempPath  string
	Chunks    int
}

func (i *Index) CreateFileInfo(ctx context.Context, f FileInfo) error {
	q := `insert into files (id, filename, temp_path, size, status, created_at) values (?, ?, ?, ?, ?, ?)`
	_, err := i.db.ExecContext(ctx, q, f.ID, f.Filename, f.TempPath, f.Size, f.Status, f.CreatedAt.UTC().Format(time.RFC3339Nano))
	return err
}

func (i *Index) Stat(ctx context.Context, id string) (FileInfo, error) {
	var f FileInfo
	var createdAt string
	err := i.db.QueryRowContext(ctx,
		`select id, filename, temp_path, size, status, created_at, (select count(*) from chunks where file_id = files.id)
		 from files where id = ?`, id,
	).Scan(&f.ID, &f.Filename, &f.TempPath, &f.Size, &f.Status, &createdAt, &f.Chunks)
	if errors.Is(err, sql.ErrNoRows) {
		return FileInfo{}, ErrNotFound
	}
	if err != nil {
		return FileInfo{}, err
	}

	f.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
	return f, nil
}

func (i *Index) SetStatus(ctx context.Context, id string, status string) error {
	q := `update files set status = ? where id = ?`
	res, err := i.db.ExecContext(ctx, q, status, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

func (i *Index) Files() ([]FileInfo, error) {
	q := `select id, filename, temp_path, size, status, created_at, (select count(*) from chunks where file_id = files.id)
		 from files`
	rows, err := i.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]FileInfo, 0)
	for rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		var f FileInfo
		var createdAt string
		if err := rows.Scan(&f.ID, &f.Filename, &f.TempPath, &f.Size, &f.Status, &createdAt, &f.Chunks); err != nil {
			return nil, err
		}
		f.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
		files = append(files, f)
	}
	return files, nil
}

func (i *Index) Delete(ctx context.Context, id string) error {
	cq := `delete from chunks where file_id = ?;`
	_, err := i.db.ExecContext(ctx, cq, id)
	if err != nil {
		return err
	}

	fq := `delete from files where id = ?;`
	_, err = i.db.ExecContext(ctx, fq, id)
	return err
}
