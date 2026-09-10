package repository

import (
	"context"
	"database/sql"
	"errors"

	"note-api/entity"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepo interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type mysqlUserRepo struct {
	db *sql.DB
}

func NewmysqlUserRepo(db *sql.DB) UserRepo {
	return &mysqlUserRepo{db: db}
}

func (r *mysqlUserRepo) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (email, password_hash, created_at) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, user.Email, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil

}
func (r *mysqlUserRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	query := `SELECT id,email,password_hash,created_at FROM users WHERE id = ?`
	var user entity.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil

}

func (r *mysqlUserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id,email,password_hash,created_at FROM users WHERE email = ?`
	var user entity.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
