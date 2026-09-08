import { useState, type FormEvent } from "react";
import { Spinner } from "./States";

export interface NoteFormValues {
  title: string;
  content: string;
}

interface NoteFormProps {
  initialValues?: NoteFormValues;
  submitLabel: string;
  busy?: boolean;
  /** Field-level errors surfaced by the backend, keyed loosely by message content. */
  serverError?: string | null;
  onSubmit: (values: NoteFormValues) => void;
  onCancel: () => void;
}

const TITLE_MAX = 100;

/**
 * Mirrors note-api/service.Update's validation (service/service.go):
 * title/content required (after trim), title <= 100 chars. CreateNote
 * only checks required-ness, but applying the stricter rule up front
 * avoids a round trip that would fail anyway on update.
 */
function validate(values: NoteFormValues): Partial<Record<keyof NoteFormValues, string>> {
  const errors: Partial<Record<keyof NoteFormValues, string>> = {};
  if (values.title.trim() === "") errors.title = "Title is required.";
  else if (values.title.length > TITLE_MAX)
    errors.title = `Title must be less than ${TITLE_MAX} characters.`;
  if (values.content.trim() === "") errors.content = "Content is required.";
  return errors;
}

export function NoteForm({
  initialValues,
  submitLabel,
  busy = false,
  serverError,
  onSubmit,
  onCancel,
}: NoteFormProps) {
  const [values, setValues] = useState<NoteFormValues>(
    initialValues ?? { title: "", content: "" }
  );
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const errors = validate(values);

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setTouched({ title: true, content: true });
    if (Object.keys(errors).length > 0) return;
    onSubmit(values);
  }

  return (
    <form onSubmit={handleSubmit} noValidate className="space-y-5">
      {serverError && (
        <div
          role="alert"
          className="rounded-md border border-rust-500/30 bg-rust-500/5 px-3.5 py-2.5 text-sm text-rust-600"
        >
          {serverError}
        </div>
      )}

      <div>
        <div className="flex items-baseline justify-between">
          <label htmlFor="title" className="text-sm font-medium text-ink">
            Title
          </label>
          <span className="font-mono text-xs text-ink-faint">
            {values.title.length}/{TITLE_MAX}
          </span>
        </div>
        <input
          id="title"
          type="text"
          value={values.title}
          onChange={(e) => setValues((v) => ({ ...v, title: e.target.value }))}
          onBlur={() => setTouched((t) => ({ ...t, title: true }))}
          aria-invalid={touched.title && !!errors.title}
          aria-describedby={errors.title ? "title-error" : undefined}
          className="mt-1.5 w-full rounded-md border border-paper-line bg-paper-card px-3.5 py-2.5 text-ink placeholder:text-ink-faint focus:outline-none focus:ring-2 focus:ring-moss-300"
          placeholder="Give your note a title"
        />
        {touched.title && errors.title && (
          <p id="title-error" className="mt-1.5 text-sm text-rust-600">
            {errors.title}
          </p>
        )}
      </div>

      <div>
        <label htmlFor="content" className="text-sm font-medium text-ink">
          Content
        </label>
        <textarea
          id="content"
          value={values.content}
          onChange={(e) => setValues((v) => ({ ...v, content: e.target.value }))}
          onBlur={() => setTouched((t) => ({ ...t, content: true }))}
          aria-invalid={touched.content && !!errors.content}
          aria-describedby={errors.content ? "content-error" : undefined}
          rows={8}
          className="mt-1.5 w-full rounded-md border border-paper-line bg-paper-card px-3.5 py-2.5 text-ink placeholder:text-ink-faint focus:outline-none focus:ring-2 focus:ring-moss-300 resize-y"
          placeholder="Write what's on your mind…"
        />
        {touched.content && errors.content && (
          <p id="content-error" className="mt-1.5 text-sm text-rust-600">
            {errors.content}
          </p>
        )}
      </div>

      <div className="flex items-center gap-2 pt-1">
        <button
          type="submit"
          disabled={busy}
          className="inline-flex items-center gap-2 rounded-md bg-moss-500 px-4 py-2.5 text-sm font-medium text-white hover:bg-moss-600 disabled:opacity-60 transition-colors"
        >
          {busy && <Spinner className="h-4 w-4" />}
          {busy ? "Saving…" : submitLabel}
        </button>
        <button
          type="button"
          onClick={onCancel}
          disabled={busy}
          className="rounded-md px-4 py-2.5 text-sm font-medium text-ink-soft hover:bg-paper transition-colors"
        >
          Cancel
        </button>
      </div>
    </form>
  );
}
