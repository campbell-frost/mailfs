package index

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserRepo struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

type User struct {
	ID        int
	Sub       string
	Email     string
	Name      string
	CreatedAt time.Time
}

// Upsert creates the user on first login and refreshes email/name on later logins
func (ur *UserRepo) Upsert(ctx context.Context, u User) (User, error) {
	q := `insert into users (sub, email, name, created_at) values (?, ?, ?, ?)
	      on conflict (sub) do update set email = excluded.email, name = excluded.name
	      returning id, created_at`
	var createdAt string
	err := ur.db.
		QueryRowContext(ctx, q, u.Sub, u.Email, u.Name, time.Now().UTC().Format(time.RFC3339Nano)).
		Scan(&u.ID, &createdAt)
	if err != nil {
		return User{}, err
	}

	u.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
	return u, nil
}

func (ur *UserRepo) Get(ctx context.Context, id int) (User, error) {
	var u User
	var createdAt string
	err := ur.db.QueryRowContext(ctx,
		`select id, sub, email, name, created_at from users where id = ?`, id,
	).Scan(&u.ID, &u.Sub, &u.Email, &u.Name, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, err
	}

	u.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
	return u, nil
}
