package repository

import (
	"context"
	"database/sql"
	"errors"
	"note-api/entity"
)

var ErrNoteNotFound = errors.New("note not found")

type NoteRepo interface {
	Create(ctx context.Context, note *entity.Note) error
	GetByID(ctx context.Context, id int) (*entity.Note, error)
	GetAll(ctx context.Context, page, limit int) ([]entity.Note, error)
	DeleteByID(ctx context.Context, id int) error
	Update(ctx context.Context, note *entity.Note) error
}

type mysqlNoteRepository struct {
	db *sql.DB
}

func NewmysqlNoteRepository(db *sql.DB) NoteRepo {
	return &mysqlNoteRepository{db: db}
}

func (r *mysqlNoteRepository) Create(ctx context.Context, note *entity.Note) error {
	query := "INSERT INTO notes (user_id, title, content, created_at) VALUES (?, ?, ?, ?)"

	result, err := r.db.ExecContext(ctx, query, note.UserID, note.Title, note.Content, note.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	note.ID = int(id)
	return nil
}

func (r *mysqlNoteRepository) GetByID(ctx context.Context, id int) (*entity.Note, error) {
	query := "SELECT id, title, content, created_at FROM notes WHERE id = ?"
	var note entity.Note
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}

	return &note, nil
}

func (r *mysqlNoteRepository) GetAll(ctx context.Context, page, limit int) ([]entity.Note, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, title, content, created_at
		FROM notes
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []entity.Note
	for rows.Next() {
		var note entity.Note
		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
func (r *mysqlNoteRepository) DeleteByID(ctx context.Context, id int) error {

	query := "DELETE FROM notes WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNoteNotFound
	}
	return nil

}

func (r *mysqlNoteRepository) Update(ctx context.Context, note *entity.Note) error {
	query := "UPDATE notes SET title = ?, content = ? WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, note.Title, note.Content, note.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNoteNotFound
	}
	return nil
}
