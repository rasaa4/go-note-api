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
	GetAll(ctx context.Context) ([]entity.Note, error)
}

type mysqlNoteRepository struct {
	db *sql.DB
}

func NewmysqlNoteRepository(db *sql.DB) NoteRepo {
	return &mysqlNoteRepository{db: db}
}

func (r *mysqlNoteRepository) Create(ctx context.Context, note *entity.Note) error {
	query := "INSERT INTO notes (title, content, created_at) VALUES (?, ?, ?)"

	result, err := r.db.ExecContext(ctx, query, note.Title, note.Content, note.CreatedAt)
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

func (r *mysqlNoteRepository) GetAll(ctx context.Context) ([]entity.Note, error) {
	query := "SELECT id, title, content, created_at FROM notes"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []entity.Note
	for rows.Next() {
		var note entity.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}
