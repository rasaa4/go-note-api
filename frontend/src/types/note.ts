/**
 * Mirrors the Go backend's entity.Note struct (note-api/entity/entity.go).
 *
 *   type Note struct {
 *       ID        int       `json:"id"`
 *       Title     string    `json:"title"`
 *       Content   string    `json:"content"`
 *       CreatedAt time.Time `json:"created_at"`
 *   }
 */
export interface Note {
  id: number;
  title: string;
  content: string;
  created_at: string;
}

/** Body for POST /notes (note-api/handler.CreateNoteRequest). */
export interface CreateNoteInput {
  title: string;
  content: string;
}

/**
 * Body for PUT /notes/update?id=.
 * The backend decodes the full entity.Note from the body and then
 * overwrites note.ID with the "id" query parameter, so id/created_at
 * in the body are ignored by the server — only title/content matter.
 */
export interface UpdateNoteInput {
  title: string;
  content: string;
}

/** Shape of every error response from note-api/handler/error.go's ErrorResponse. */
export interface ApiErrorBody {
  error: string;
}
