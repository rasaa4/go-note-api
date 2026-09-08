import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { NoteForm, type NoteFormValues } from "../components/NoteForm";
import { getNoteById, updateNote } from "../api/notes";
import { ApiError } from "../api/ApiError";
import { useToast } from "../components/ToastProvider";
import { ErrorState, Spinner } from "../components/States";

export function EditNotePage() {
  const { id } = useParams<{ id: string }>();
  const noteId = Number(id);
  const navigate = useNavigate();
  const { notify } = useToast();

  const [initialValues, setInitialValues] = useState<NoteFormValues | null>(null);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const [serverError, setServerError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    async function load() {
      try {
        const note = await getNoteById(noteId);
        if (!cancelled) setInitialValues({ title: note.title, content: note.content });
      } catch (err) {
        if (!cancelled) {
          setLoadError(
            err instanceof ApiError ? err.message : "Could not load this note."
          );
        }
      }
    }
    if (Number.isFinite(noteId)) load();
    else setLoadError("Invalid note id.");
    return () => {
      cancelled = true;
    };
  }, [noteId]);

  async function handleSubmit(values: NoteFormValues) {
    setBusy(true);
    setServerError(null);
    try {
      await updateNote(noteId, values);
      notify("success", "Note updated.");
      navigate("/");
    } catch (err) {
      setServerError(
        err instanceof ApiError ? err.message : "Could not update this note."
      );
    } finally {
      setBusy(false);
    }
  }

  if (loadError) {
    return <ErrorState message={loadError} />;
  }

  if (!initialValues) {
    return (
      <div className="flex items-center gap-2 text-ink-soft text-sm">
        <Spinner className="h-4 w-4" />
        Loading note…
      </div>
    );
  }

  return (
    <div className="max-w-xl">
      <h1 className="font-display text-2xl text-ink mb-6">Edit note</h1>
      <NoteForm
        initialValues={initialValues}
        submitLabel="Save changes"
        busy={busy}
        serverError={serverError}
        onSubmit={handleSubmit}
        onCancel={() => navigate("/")}
      />
    </div>
  );
}
