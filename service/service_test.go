package service

import (
	"context"
	"errors"
	"note-api/entity"
	"testing"
)

type mockNoteRepo struct{}

func (m *mockNoteRepo) Create(ctx context.Context, note *entity.Note) error {
	note.ID = 1
	return nil
}

func (m *mockNoteRepo) GetByID(ctx context.Context, id int) (*entity.Note, error) {
	return nil, nil
}

func (m *mockNoteRepo) GetAll(ctx context.Context) ([]entity.Note, error) {
	return nil, nil
}

func (m *mockNoteRepo) DeleteByID(ctx context.Context, id int) error {
	return nil
}

func (m *mockNoteRepo) Update(ctx context.Context, note *entity.Note) error {
	return nil
}

func TestCreateNote(t *testing.T) {

	repo := &mockNoteRepo{}
	service := NewNoteService(repo)
	note, err := service.CreateNote(context.Background(), "My Note", "Hello i'm testing")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if note == nil {
		t.Fatalf("expected note, got nil")
	}
	if note.Title != "My Note" {
		t.Fatalf("expected title My note, got %s", note.Title)
	}
	if note.Content != "Hello i'm testing" {
		t.Fatalf("expected content Hello i'm testing, got %s", note.Content)
	}
	if note.ID != 1 {
		t.Errorf("expected ID 1, got %d", note.ID)
	}
}

func TestCreateNote_TitleRequired(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)
	note, err := service.CreateNote(context.Background(), " ", "Hello")

	if !errors.Is(err, ErrTitleRequired) {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}
	if note != nil {
		t.Fatalf("expected nil note, got %+v", note)
	}

}

func TestCreateNote_ContentRequired(t *testing.T) {

	repo := &mockNoteRepo{}
	service := NewNoteService(repo)
	note, err := service.CreateNote(context.Background(), "Mynote", " ")
	if !errors.Is(err, ErrContentRequired) {
		t.Fatalf("expected ErrContentRequired, got %v", err)
	}
	if note != nil {
		t.Fatalf("expected nil note, got %+v", note)
	}

}
