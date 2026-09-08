import { Link } from "react-router-dom";
import type { Note } from "../types/note";
import { formatDate, truncate } from "../lib/format";

interface NoteCardProps {
  note: Note;
  onDelete: (note: Note) => void;
}

export function NoteCard({ note, onDelete }: NoteCardProps) {
  return (
    <article className="group relative rounded-lg border border-paper-line bg-paper-card shadow-card overflow-hidden">
      <div className="absolute inset-x-0 top-0 h-1 bg-moss-500" />
      <div className="p-5">
        <div className="flex items-start justify-between gap-3">
          <h3 className="font-display text-xl leading-snug text-ink line-clamp-2">
            {note.title}
          </h3>
          <span className="shrink-0 font-mono text-[11px] tracking-tight text-ink-faint mt-1">
            #{note.id}
          </span>
        </div>

        <p className="mt-2 text-sm text-ink-soft leading-relaxed line-clamp-3">
          {truncate(note.content, 160)}
        </p>

        <div className="mt-4 flex items-center justify-between border-t border-dashed border-paper-line pt-3">
          <time className="font-mono text-xs text-ink-faint" dateTime={note.created_at}>
            {formatDate(note.created_at)}
          </time>

          <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 group-focus-within:opacity-100 transition-opacity">
            <Link
              to={`/notes/${note.id}`}
              className="rounded px-2 py-1 text-xs font-medium text-ink-soft hover:text-moss-700 hover:bg-moss-50 transition-colors"
            >
              View
            </Link>
            <Link
              to={`/notes/${note.id}/edit`}
              className="rounded px-2 py-1 text-xs font-medium text-ink-soft hover:text-moss-700 hover:bg-moss-50 transition-colors"
            >
              Edit
            </Link>
            <button
              type="button"
              onClick={() => onDelete(note)}
              className="rounded px-2 py-1 text-xs font-medium text-ink-soft hover:text-rust-600 hover:bg-rust-500/10 transition-colors"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </article>
  );
}
