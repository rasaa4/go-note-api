# Notes — frontend for `note-api`

A React + TypeScript frontend for the Go `note-api` backend: a notes
dashboard with create, view, edit, and delete, built directly against
the backend's actual handlers (not an assumed or invented API).

This README also documents two backend gaps discovered while
inspecting the Go code, and the exact (minimal) fix for each. The
backend was **not modified** as part of this work.

---

## Tech stack

- React 19 + TypeScript
- Vite
- React Router (client-side routing)
- Tailwind CSS

## Project structure

```
src/
  api/
    client.ts       fetch wrapper: base URL, query params, error normalization
    ApiError.ts     typed error class (status code + message)
    notes.ts        getNotes / getNoteById / createNote / updateNote / deleteNote
  types/
    note.ts         Note, CreateNoteInput, UpdateNoteInput — mirrors entity.Note
  components/
    Layout.tsx          header + page shell
    NoteCard.tsx         a single note in the dashboard grid
    NoteForm.tsx         shared form for create/edit, with validation
    Pagination.tsx       Previous/Next controls + page size
    ConfirmDialog.tsx    confirmation modal for delete
    ToastProvider.tsx    success/error toast notifications
    States.tsx           loading / empty / error state UI
  hooks/
    useNotesList.ts  data fetching + pagination state for the dashboard
  pages/
    NotesDashboard.tsx
    CreateNotePage.tsx
    EditNotePage.tsx
    NoteViewPage.tsx
  lib/
    format.ts        date formatting, text truncation
```

No component calls `fetch` directly — every request goes through
`src/api/notes.ts`, which is the single place that knows about the
backend's URLs and payload shapes.

---

## Running the frontend

```bash
npm install
cp .env.example .env   # edit if your backend runs somewhere other than :8080
npm run dev
```

Opens on `http://localhost:5173` by default.

### Environment variables

| Variable            | Default                 | Purpose                                       |
|----------------------|--------------------------|------------------------------------------------|
| `VITE_API_BASE_URL`  | `http://localhost:8080` | Base URL of the Go backend, no trailing slash  |

## Running the Go backend

From the `note-api` project root:

```bash
go run cmd/main.go
```

This expects a MySQL database reachable at `127.0.0.1:3306`, database
`test_db`, with a `notes` table:

```sql
CREATE TABLE notes (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255),
    content TEXT,
    created_at DATETIME
);
```

(See `cmd/main.go`'s `sql.Open` call for the exact DSN used.)

---

## Backend API reference (as implemented, traced from the code)

Base URL: `http://localhost:8080` (from `cmd/main.go`, `http.ListenAndServe(":8080", ...)`).

| Method | Path            | Handler                  | Success          | Notes |
|--------|-----------------|---------------------------|-------------------|-------|
| POST   | `/notes`        | `NoteHandler.Create`      | `201` + `Note`    | body: `{"title","content"}` |
| GET    | `/notes`        | `NoteHandler.GetAll`      | `200` + `Note[]`  | query: `page`, `limit` (default `1`/`10` if omitted) |
| GET    | `/notes/get`    | `NoteHandler.GetByID`     | `200` + `Note`    | query: `id` |
| PUT    | `/notes/update` | `NoteHandler.Update`      | `204`             | query: `id`; body: `{"title","content"}` — **route not registered, see below** |
| DELETE | `/notes/delete` | `NoteHandler.DeleteByID`  | `204`             | query: `id` — **route not registered, see below** |

`Note` shape (from `entity/entity.go`):
```ts
{ id: number; title: string; content: string; created_at: string }
```

### Error format

Most errors come back as JSON, via `handler/error.go`'s `handleError`:
```json
{ "error": "title is required" }
```
mapped as: validation errors (`service.Err*`) → `400`; `repository.ErrNoteNotFound` → `404`; anything else → `500` with a generic `"internal server error"` message (the real error is not leaked to the client).

**Exception:** malformed JSON bodies and unparsable `id`/`page`/`limit`
query parameters are rejected *before* reaching `handleError`, using
Go's standard `http.Error()`, which returns **plain text**, not JSON
(e.g. `invalid id\n` with no `error` field). The frontend's
`src/api/client.ts` handles both shapes so the UI always shows a clean
message regardless of which path produced the error.

---

## Backend limitations discovered

These are properties of the inspected `note-api` code as it stands.
The frontend is written to work correctly against the *intended*
behavior (matching each handler's actual signature), but two features
will not function until the backend is updated as described. **No
backend files were changed by this task.**

### 1. `Update` and `DeleteByID` handlers exist but have no route

`handler/handler.go` implements `NoteHandler.Update` and
`NoteHandler.DeleteByID` in full (with correct validation and status
codes), but `cmd/main.go` never registers them on the `http.ServeMux`
— only `/notes` and `/notes/get` are wired up. Calling edit or delete
from this frontend against an unmodified backend will currently get a
plain-text `404 page not found` from Go's default mux (not the
backend's own JSON 404).

**Fix** (add to `cmd/main.go`, alongside the existing `mux.HandleFunc` calls):
```go
mux.HandleFunc("/notes/update", noteHandler.Update)
mux.HandleFunc("/notes/delete", noteHandler.DeleteByID)
```

### 2. No CORS headers anywhere in the backend

Nothing in `cmd/main.go` or `handler/*.go` sets
`Access-Control-Allow-Origin` or handles `OPTIONS` preflight requests.
A browser-based frontend running on a different origin (e.g. Vite's
`http://localhost:5173`) will have every request blocked by the
browser's CORS policy — this affects *all* endpoints, including the
already-working `GET /notes`.

**Fix** (minimal, wrap the mux once in `cmd/main.go` before passing it to `ListenAndServe`):
```go
func withCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}
// then: http.ListenAndServe(":8080", withCORS(mux))
```

### 3. `GetAll` has no total-count / has-more metadata

`repository.GetAll` returns a plain `[]entity.Note` with no total
record count, total page count, or "is there a next page" flag. This
makes true numbered pagination impossible to implement correctly from
the frontend alone.

**How the frontend handles it:** Previous/Next controls only. "Next"
is disabled when the current page returns fewer notes than the
requested `limit` — a reasonable heuristic, not a guarantee. If you
want reliable "last page" detection, the backend would need to return
either a total count or a computed `has_more` boolean alongside the
notes array.

### 4. Inconsistent error response format (see "Error format" above)

Not a bug that blocks anything, but worth knowing: some 400s are
JSON, some are plain text, depending on which layer rejects the
request first.

---

## Verification performed

- `npx tsc --noEmit` — no type errors.
- `npm run build` — production build succeeds.
- Full CRUD cycle (create → list → get → update → delete → confirm 404)
  verified against a minimal Go server implementing the exact same
  route/payload/status-code contract as `note-api`'s handlers —
  including both currently-unregistered routes and CORS — to confirm
  the frontend's request/response handling is correct end-to-end.
  This substitute server was used only for testing in the sandbox this
  was built in (which has no MySQL available) and is not part of the
  delivered frontend; against your real backend, apply the two fixes
  above first.
- Validation errors, 404s, and network-unreachable states were each
  triggered manually and confirmed to show a specific, readable
  message rather than a generic failure.

---

## Assumptions made

- `page`/`limit` query params are integers starting at `1` (matches
  `service.GetAll`'s `page <= 0` / `limit <= 0` validation).
- `Update`'s request body only needs `title`/`content` — confirmed
  from `handler.Update`, which decodes the body into `entity.Note` but
  then unconditionally overwrites `note.ID` from the `id` query
  parameter, so any `id`/`created_at` sent in the body is ignored by
  the server.
- The backend will be reachable at `http://localhost:8080` in
  development, per `cmd/main.go`; configurable via `VITE_API_BASE_URL`
  for any other environment.
