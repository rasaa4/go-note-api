import { apiRequest } from "./client";
import type { CreateNoteInput, Note, UpdateNoteInput } from "../types/note";

/**
 * These functions map 1:1 onto the routes actually wired up in
 * note-api/cmd/main.go and handled in note-api/handler/handler.go.
 * Nothing here is invented — see README.md "Backend API reference"
 * for the full trace from route to handler to service.
 *
 * Two of them (updateNote, deleteNote) call handler methods that
 * exist in the Go code but are NOT currently registered on any route
 * in main.go. Calling them against an unmodified backend will 404.
 * See README.md "Backend limitations" for the exact fix.
 */

/** GET /notes?page=&limit=  → 200 Note[] */
export async function getNotes(page: number, limit: number): Promise<Note[]> {
  return apiRequest<Note[]>("/notes", {
    method: "GET",
    query: { page, limit },
  });
}

/** GET /notes/get?id=  → 200 Note | 404 */
export async function getNoteById(id: number): Promise<Note> {
  return apiRequest<Note>("/notes/get", {
    method: "GET",
    query: { id },
  });
}

/** POST /notes  → 201 Note */
export async function createNote(data: CreateNoteInput): Promise<Note> {
  return apiRequest<Note>("/notes", {
    method: "POST",
    body: data,
  });
}

/**
 * PUT /notes/update?id=  → 204 No Content
 *
 * NOT REGISTERED in main.go on the inspected backend — see
 * README.md "Backend limitations". Implemented here against
 * handler.Update's actual signature so it works the moment the
 * route is added.
 */
export async function updateNote(id: number, data: UpdateNoteInput): Promise<void> {
  return apiRequest<void>("/notes/update", {
    method: "PUT",
    query: { id },
    body: data,
  });
}

/**
 * DELETE /notes/delete?id=  → 204 No Content
 *
 * NOT REGISTERED in main.go on the inspected backend — see
 * README.md "Backend limitations".
 */
export async function deleteNote(id: number): Promise<void> {
  return apiRequest<void>("/notes/delete", {
    method: "DELETE",
    query: { id },
  });
}
