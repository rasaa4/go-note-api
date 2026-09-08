import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { getNoteById, deleteNote } from "../api/notes";
import { ApiError } from "../api/ApiError";
import { ErrorState, Spinner } from "../components/States";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { useToast } from "../components/ToastProvider";
import { formatDateTime } from "../lib/format";
import type { Note } from "../types/note";

export function NoteViewPage() {
  const { id } = useParams<{ id: string }>();
  const noteId = Number(id);
  const navigate = useNavigate();
  const { notify } = useToast();

  const [note, setNote] = useState<Note | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    let cancelled = false;
    async function load() {
      try {
        const result = await getNoteById(noteId);
        if (!cancelled) setNote(result);
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof ApiError ? err.message : "Could not load this note.");
        }
      }
    }
    if (Number.isFinite(noteId)) load();
    else setError("Invalid note id.");
    return () => {
      cancelled = true;
    };
  }, [noteId]);

  async function handleDelete() {
    setDeleting(true);
    try {
      await deleteNote(noteId);
      notify("success", "Note deleted.");
      navigate("/");
    } catch (err) {
      notify("error", err instanceof ApiError ? err.message : "Could not delete this note.");
    } finally {
      setDeleting(false);
    }
  }

  if (error) return <ErrorState message={error} />;

  if (!note) {
    return (
      <div className="flex items-center gap-2 text-ink-soft text-sm">
        <Spinner className="h-4 w-4" />
        Loading note…
      </div>
    );
  }

  return (
    <div className="max-w-2xl">
      <Link to="/" className="text-sm text-ink-soft hover:text-ink transition-colors">
        ← Back to notes
      </Link>

      <div className="mt-4 rounded-lg border border-paper-line bg-paper-card shadow-card p-6">
        <div className="flex items-start justify-between gap-4">
          <h1 className="font-display text-2xl text-ink">{note.title}</h1>
          <span className="shrink-0 font-mono text-xs text-ink-faint">#{note.id}</span>
        </div>
        <time className="mt-1 block font-mono text-xs text-ink-faint" dateTime={note.created_at}>
          {formatDateTime(note.created_at)}
        </time>
        <p className="mt-5 text-ink-soft leading-relaxed whitespace-pre-wrap">
          {note.content}
        </p>

        <div className="mt-6 flex items-center gap-2 border-t border-paper-line pt-4">
          <Link
            to={`/notes/${note.id}/edit`}
            className="rounded-md bg-ink px-3.5 py-2 text-sm font-medium text-white hover:bg-ink/90 transition-colors"
          >
            Edit
          </Link>
          <button
            type="button"
            onClick={() => setConfirmOpen(true)}
            className="rounded-md px-3.5 py-2 text-sm font-medium text-rust-600 hover:bg-rust-500/10 transition-colors"
          >
            Delete
          </button>
        </div>
      </div>

      <ConfirmDialog
        open={confirmOpen}
        title="Delete this note?"
        description={`"${note.title}" will be permanently deleted. This can't be undone.`}
        confirmLabel="Delete"
        danger
        busy={deleting}
        onConfirm={handleDelete}
        onCancel={() => setConfirmOpen(false)}
      />
    </div>
  );
}
