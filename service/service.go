package service

import (
	"context"
	"errors"
	"note-api/entity"
	"note-api/repository"
	"time"
)

type NoteService struct {
	repo repository.NoteRepo
}

func NewNoteService(repo repository.NoteRepo) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) CreateNote(ctx context.Context, title, content string) (*entity.Note, error) {
	if title == "" || content == "" {
		return nil, errors.New("title and content is required")
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

	if id == 0 {
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
