package index

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type FileRepo struct {
	db *sql.DB
}

func NewFile(db *sql.DB) *FileRepo {
	return &FileRepo{db: db}
}

type FileInfo struct {
	ID        string
	Filename  string
	Size      int64
	CreatedAt time.Time
	Status    string
	TempPath  string
	Chunks    int
}

func (fr *FileRepo) Create(ctx context.Context, f FileInfo) error {
	q := `insert into files (id, filename, temp_path, size, status, created_at) values (?, ?, ?, ?, ?, ?)`
	_, err := fr.db.ExecContext(ctx, q, f.ID, f.Filename, f.TempPath, f.Size, f.Status, f.CreatedAt.UTC().Format(time.RFC3339Nano))
	return err
}

func (fr *FileRepo) Get(ctx context.Context, id string) (FileInfo, error) {
	var f FileInfo
	var createdAt string
	err := fr.db.QueryRowContext(ctx,
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

type FileStatus struct {
	Status          string
	ProcessedChunks int
}

func (fr *FileRepo) Status(ctx context.Context, id string) (FileStatus, error) {
	var s FileStatus
	const q = `
	select
		f.status,
		(select count(*) from chunks where file_id = f.id)
	from files f
	where f.id = ?`

	if err := fr.db.QueryRowContext(ctx, q, id).Scan(&s.Status, &s.ProcessedChunks); err != nil {
		return FileStatus{}, err
	}

	return FileStatus{
		Status:          s.Status,
		ProcessedChunks: s.ProcessedChunks,
	}, nil
}

func (fr *FileRepo) SetStatus(ctx context.Context, id string, status string) error {
	q := `update files set status = ? where id = ?`
	res, err := fr.db.ExecContext(ctx, q, status, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

func (fr *FileRepo) List() ([]FileInfo, error) {
	q := `select id, filename, temp_path, size, status, created_at, (select count(*) from chunks where file_id = files.id)
		 from files`
	rows, err := fr.db.Query(q)
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

func (fr *FileRepo) Delete(ctx context.Context, id string) error {
	cq := `delete from chunks where file_id = ?;`
	_, err := fr.db.ExecContext(ctx, cq, id)
	if err != nil {
		return err
	}

	fq := `delete from files where id = ?;`
	_, err = fr.db.ExecContext(ctx, fq, id)
	return err
}
