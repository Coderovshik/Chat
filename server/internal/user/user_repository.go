package user

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Repository struct {
	db DBTX
}

func NewRepository(db DBTX) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	const query = "INSERT INTO users(username, email, passhash) VALUES ($1, $2, $3) RETURNING id"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, user.Username, user.Email, user.Passhash).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	const query = "SELECT id, username, email, passhash FROM users WHERE email = $1"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	u := User{}
	err = stmt.QueryRowContext(ctx, email).Scan(&u.ID, &u.Username, &u.Email, &u.Passhash)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
