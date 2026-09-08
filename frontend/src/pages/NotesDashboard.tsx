import { useState } from "react";
import { Link } from "react-router-dom";
import { useNotesList } from "../hooks/useNotesList";
import { NoteCard } from "../components/NoteCard";
import { LoadingGrid, EmptyState, ErrorState } from "../components/States";
import { Pagination } from "../components/Pagination";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { useToast } from "../components/ToastProvider";
import { deleteNote } from "../api/notes";
import { ApiError } from "../api/ApiError";
import type { Note } from "../types/note";

export function NotesDashboard() {
  const { notes, page, limit, status, errorMessage, setPage, setLimit, refetch } =
    useNotesList();
  const { notify } = useToast();
  const [pendingDelete, setPendingDelete] = useState<Note | null>(null);
  const [deleting, setDeleting] = useState(false);

  async function handleConfirmDelete() {
    if (!pendingDelete) return;
    setDeleting(true);
    try {
      await deleteNote(pendingDelete.id);
      notify("success", `"${pendingDelete.title}" was deleted.`);
      setPendingDelete(null);
      refetch();
    } catch (err) {
      const message =
        err instanceof ApiError ? err.message : "Could not delete this note.";
      notify("error", message);
    } finally {
      setDeleting(false);
    }
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="font-display text-3xl text-ink">Your notes</h1>
        <p className="mt-1 text-sm text-ink-soft">
          Everything you've written down, in one place.
        </p>
      </div>

      {status === "loading" && <LoadingGrid />}

      {status === "error" && <ErrorState message={errorMessage ?? "Failed to load notes."} onRetry={refetch} />}

      {status === "success" && notes.length === 0 && (
        <EmptyState
          title="No notes yet"
          description="Create your first note to see it appear here."
          action={
            <Link
              to="/notes/new"
              className="inline-flex rounded-md bg-moss-500 px-4 py-2 text-sm font-medium text-white hover:bg-moss-600 transition-colors"
            >
              New note
            </Link>
          }
        />
      )}

      {status === "success" && notes.length > 0 && (
        <>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {notes.map((note) => (
              <NoteCard key={note.id} note={note} onDelete={setPendingDelete} />
            ))}
          </div>

          <Pagination
            page={page}
            limit={limit}
            resultCount={notes.length}
            onPageChange={setPage}
            onLimitChange={setLimit}
          />
        </>
      )}

      <ConfirmDialog
        open={pendingDelete !== null}
        title="Delete this note?"
        description={
          pendingDelete
            ? `"${pendingDelete.title}" will be permanently deleted. This can't be undone.`
            : ""
        }
        confirmLabel="Delete"
        danger
        busy={deleting}
        onConfirm={handleConfirmDelete}
        onCancel={() => setPendingDelete(null)}
      />
    </div>
  );
}
