package service

import (
	"context"
	"errors"
	"note-api/entity"
	"note-api/repository"
	"strings"
	"time"
)

var (
	ErrInvalidID       = errors.New("invalid id")
	ErrTitleRequired   = errors.New("title is required")
	ErrContentRequired = errors.New("content is required")
	ErrTitleTooLong    = errors.New("title must be less than 100 characters")
	ErrNoteRequired    = errors.New("note is required")
	ErrInvalidPage     = errors.New("invalid page")
	ErrInvalidLimit    = errors.New("invalid limit")
)

type NoteService struct {
	repo repository.NoteRepo
}

func NewNoteService(repo repository.NoteRepo) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) CreateNote(ctx context.Context, userID int, title, content string) (*entity.Note, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ErrTitleRequired
	}
	if strings.TrimSpace(content) == "" {
		return nil, ErrContentRequired
	}
	var note entity.Note
	note.UserID = userID
	note.Title = title
	note.Content = content
	note.CreatedAt = time.Now()
	err := s.repo.Create(ctx, &note)
	if err != nil {
		return nil, err
	}
	return &note, nil

}

func (s *NoteService) GetByID(ctx context.Context, id int) (*entity.Note, error) {

	if id <= 0 {
		return nil, ErrInvalidID
	}
	note, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return note, nil

}
func (s *NoteService) GetAll(ctx context.Context, page, limit int) ([]entity.Note, error) {
	if page <= 0 {
		return nil, ErrInvalidPage
	}
	if limit <= 0 {
		return nil, ErrInvalidLimit
	}
	notes, err := s.repo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (s *NoteService) DeleteByID(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidID
	}
	err := s.repo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
func (s *NoteService) Update(ctx context.Context, note *entity.Note) error {
	if note == nil {
		return ErrNoteRequired
	}
	if note.ID <= 0 {
		return ErrInvalidID
	}
	if strings.TrimSpace(note.Title) == "" {
		return ErrTitleRequired
	}
	if strings.TrimSpace(note.Content) == "" {
		return ErrContentRequired
	}
	if len(note.Title) > 100 {
		return ErrTitleTooLong
	}
	return s.repo.Update(ctx, note)

}
