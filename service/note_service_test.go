package service

import (
	"context"
	"errors"
	"note-api/entity"
	"strings"
	"testing"
)

type mockNoteRepo struct {
	createErr error
	gotUserID int
	gotPage   int
	gotLimit  int
}

func (m *mockNoteRepo) Create(ctx context.Context, note *entity.Note) error {
	if m.createErr != nil {

		return m.createErr
	}
	note.ID = 1
	return nil
}

func (m *mockNoteRepo) GetByID(ctx context.Context, id int) (*entity.Note, error) {
	return &entity.Note{
		ID:      1,
		Title:   "My Note",
		Content: "Hello",
	}, nil
}

func (m *mockNoteRepo) GetAll(ctx context.Context, userID, page, limit int) ([]entity.Note, error) {
	m.gotUserID = userID
	m.gotPage = page
	m.gotLimit = limit
	return []entity.Note{
		{
			ID:      1,
			Title:   "Note 1",
			Content: "Content 1",
		},
		{
			ID:      2,
			Title:   "Note 2",
			Content: "Content 2",
		},
	}, nil
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
	note, err := service.CreateNote(context.Background(), 1, "My Note", "Hello i'm testing")

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
	if note.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", note.UserID)
	}
}

func TestCreateNote_TitleRequired(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)
	note, err := service.CreateNote(context.Background(), 1, " ", "Hello")

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
	note, err := service.CreateNote(context.Background(), 1, "Mynote", " ")
	if !errors.Is(err, ErrContentRequired) {
		t.Fatalf("expected ErrContentRequired, got %v", err)
	}
	if note != nil {
		t.Fatalf("expected nil note, got %+v", note)
	}

}
func TestGetID(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)

	note, err := service.GetByID(context.Background(), 1)
	if err != nil {

		t.Fatalf("expected no error, got %v", err)
	}
	if note == nil {
		t.Fatalf("expected note, got nil")
	}
	if note.ID != 1 {
		t.Errorf("expected ID 1, got %d", note.ID)
	}
	if note.Title != "My Note" {
		t.Errorf("expected title My Note, got %s", note.Title)
	}
	if note.Content != "Hello" {
		t.Errorf("expected content Hello, got %s", note.Content)
	}

}
func TestGetByID_InvalidID(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)
	note, err := service.GetByID(context.Background(), 0)
	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("expected ErrInvalidID, got %v", err)
	}
	if note != nil {
		t.Fatalf("expected nil note,got %+v", note)
	}
}

func TestGetAll(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)
	notes, err := service.GetAll(context.Background(), 1, 1, 10)
	if err != nil {
		t.Fatalf("expecteed no error, got %v", err)
	}
	if notes == nil {
		t.Fatalf("expecteed notes, got nil")
	}
	if len(notes) != 2 {
		t.Fatalf("expecteed 2 notes, got %v", len(notes))
	}

	if notes[0].Title != "Note 1" {
		t.Errorf("expected Note 1, got %s", notes[0].Title)

	}
	if notes[1].Title != "Note 2" {
		t.Errorf("expected Note 2, got %s", notes[1].Title)

	}
}

func TestDeleteByID(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)
	err := service.DeleteByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteByID_InvalidID(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)

	err := service.DeleteByID(context.Background(), 0)
	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("expected ErrInvalid, got %v", err)
	}

}

func TestUpdate(t *testing.T) {

	repo := &mockNoteRepo{}
	service := NewNoteService(repo)

	note := entity.Note{
		ID:      1,
		Title:   "my note",
		Content: "updated world",
	}
	err := service.Update(context.Background(), &note)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

}
func TestUpdate_NoteRequired(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)
	err := service.Update(context.Background(), nil)
	if !errors.Is(err, ErrNoteRequired) {
		t.Fatalf("expected ErrNoteRequired, got %v", err)
	}
}
func TestUpdate_InvalidID(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)
	note := entity.Note{
		ID:      0,
		Title:   "hello",
		Content: "this is content",
	}
	err := service.Update(context.Background(), &note)
	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("expected ErrInvalidID, got %v", err)
	}
}
func TestUpdate_TitleRequired(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)

	note := entity.Note{
		ID:      1,
		Title:   " ",
		Content: "this is content",
	}
	err := service.Update(context.Background(), &note)
	if !errors.Is(err, ErrTitleRequired) {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}

}
func TestUpdate_ContentRequired(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)

	note := entity.Note{
		ID:      1,
		Title:   "title",
		Content: " ",
	}
	err := service.Update(context.Background(), &note)
	if !errors.Is(err, ErrContentRequired) {
		t.Fatalf("expected ErrContentRequired, got %v", err)
	}
}

func TestUpdate_Titletoolong(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)

	note := entity.Note{
		ID:      1,
		Title:   strings.Repeat("a", 101),
		Content: "hello",
	}
	err := service.Update(context.Background(), &note)
	if !errors.Is(err, ErrTitleTooLong) {
		t.Fatalf("expected ErrTitleTooLong, got %v", err)
	}
}

func TestGetAll_InvalidPage(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)
	notes, err := service.GetAll(context.Background(), 1, 0, 10)
	if !errors.Is(err, ErrInvalidPage) {
		t.Fatalf("expected ErrInvalidPage, got %v", err)
	}
	if notes != nil {
		t.Fatalf("expected nil notes, got %v+", notes)
	}

}

func TestGetAll_InvalidLimit(t *testing.T) {

	repo := &mockNoteRepo{}
	service := NewNoteService(repo)
	notes, err := service.GetAll(context.Background(), 1, 1, 0)
	if !errors.Is(err, ErrInvalidLimit) {
		t.Fatalf("expected ErrInvalidlimit, got %v", err)
	}
	if notes != nil {
		t.Fatalf("expected nil notes, got %v+", notes)
	}
}

func TestGetAll_ReturnsNotes(t *testing.T) {
	repo := &mockNoteRepo{}
	service := NewNoteService(repo)
	notes, err := service.GetAll(context.Background(), 7, 2, 10)

	if err != nil {
		t.Fatalf("expected no error, got %v+", err)
	}
	if notes == nil {
		t.Fatalf("expected notes, got nil")
	}
	if repo.gotUserID != 7 {
		t.Errorf("expected userID 7, got %d", repo.gotUserID)
	}
	if repo.gotPage != 2 {
		t.Errorf("expected page 2, got %d", repo.gotPage)
	}
	if repo.gotLimit != 10 {
		t.Errorf("expected limit 10, got %d", repo.gotLimit)
	}
}
