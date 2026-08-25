package service

import (
	"context"
	"errors"
	"note-api/entity"
	"note-api/repository"
	"strings"
	"time"
)

type NoteService struct {
	repo repository.NoteRepo
}

func NewNoteService(repo repository.NoteRepo) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) CreateNote(ctx context.Context, title, content string) (*entity.Note, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title is required")
	}
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("content is required")
	}
	var note entity.Note
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
		return nil, errors.New("id number is required")
	}
	note, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return note, nil

}

func (s *NoteService) GetAll(ctx context.Context) ([]entity.Note, error) {
	notes, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (s *NoteService) DeleteByID(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("id number is required")
	}
	err := s.repo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
func (s *NoteService) Update(ctx context.Context, note *entity.Note) error {
	if note == nil {
		return errors.New("note is required")
	}
	if note.ID <= 0 {
		return errors.New("id number is required")
	}
	if strings.TrimSpace(note.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(note.Content) == "" {
		return errors.New("content is required")
	}
	if len(note.Title) > 100 {
		return errors.New("title must be less than 100 characters")
	}
	return s.repo.Update(ctx, note)

}
