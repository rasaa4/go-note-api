import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { NoteForm, type NoteFormValues } from "../components/NoteForm";
import { createNote } from "../api/notes";
import { ApiError } from "../api/ApiError";
import { useToast } from "../components/ToastProvider";

export function CreateNotePage() {
  const navigate = useNavigate();
  const { notify } = useToast();
  const [busy, setBusy] = useState(false);
  const [serverError, setServerError] = useState<string | null>(null);

  async function handleSubmit(values: NoteFormValues) {
    setBusy(true);
    setServerError(null);
    try {
      await createNote(values);
      notify("success", "Note created.");
      navigate("/");
    } catch (err) {
      setServerError(
        err instanceof ApiError ? err.message : "Could not create this note."
      );
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="max-w-xl">
      <h1 className="font-display text-2xl text-ink mb-6">New note</h1>
      <NoteForm
        submitLabel="Create note"
        busy={busy}
        serverError={serverError}
        onSubmit={handleSubmit}
        onCancel={() => navigate("/")}
      />
    </div>
  );
}
