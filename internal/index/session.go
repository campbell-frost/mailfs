package index

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type SessionRepo struct {
	db *sql.DB
}

func NewSession(db *sql.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

type Session struct {
	TokenHash []byte
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
}

func (sr *SessionRepo) Create(ctx context.Context, s Session) error {
	q := `insert into sessions (token_hash, user_id, created_at, expires_at) values (?, ?, ?, ?)`
	_, err := sr.db.ExecContext(ctx, q,
		s.TokenHash, s.UserID, s.CreatedAt.UTC().Format(time.RFC3339Nano), s.ExpiresAt.UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (sr *SessionRepo) Get(ctx context.Context, tokenHash []byte) (Session, error) {
	var s Session
	var createdAt, expiresAt string
	err := sr.db.QueryRowContext(ctx,
		`select token_hash, user_id, created_at, expires_at from sessions where token_hash = ?`, tokenHash,
	).Scan(&s.TokenHash, &s.UserID, &createdAt, &expiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Session{}, ErrNotFound
	}
	if err != nil {
		return Session{}, err
	}

	s.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
	s.ExpiresAt, _ = time.Parse(time.RFC3339Nano, expiresAt)
	return s, nil
}

func (sr *SessionRepo) Delete(ctx context.Context, tokenHash []byte) error {
	_, err := sr.db.ExecContext(ctx, `delete from sessions where token_hash = ?`, tokenHash)
	return err
}
