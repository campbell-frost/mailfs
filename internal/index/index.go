package index

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "modernc.org/sqlite"
)

type Index struct {
	db *sql.DB
}

//go:embed schema.sql
var schema string

func New(path string) (*Index, error) {
	db, err := initDB(path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("apply schema: %w", err)
	}
	return &Index{db: db}, nil
}

func initDB(path string) (*sql.DB, error) {
	dsn := "file:" + path + "?" +
		"_pragma=busy_timeout(5000)&" +
		"_pragma=journal_mode(WAL)&" +
		"_pragma=foreign_keys(ON)"
	db, err := sql.Open("sqlite", dsn)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return db, err
}

func (i *Index) Close() error {
	return i.db.Close()
}
