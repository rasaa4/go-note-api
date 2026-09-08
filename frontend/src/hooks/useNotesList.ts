import { useCallback, useEffect, useState } from "react";
import { getNotes } from "../api/notes";
import { ApiError } from "../api/ApiError";
import type { Note } from "../types/note";

interface UseNotesListResult {
  notes: Note[];
  page: number;
  limit: number;
  status: "loading" | "error" | "success";
  errorMessage: string | null;
  setPage: (page: number) => void;
  setLimit: (limit: number) => void;
  refetch: () => void;
}

const DEFAULT_LIMIT = 10;

export function useNotesList(): UseNotesListResult {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(DEFAULT_LIMIT);
  const [notes, setNotes] = useState<Note[]>([]);
  const [status, setStatus] = useState<"loading" | "error" | "success">("loading");
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [reloadToken, setReloadToken] = useState(0);

  const load = useCallback(async () => {
    setStatus("loading");
    setErrorMessage(null);
    try {
      const result = await getNotes(page, limit);
      setNotes(result ?? []);
      setStatus("success");
    } catch (err) {
      const message =
        err instanceof ApiError ? err.message : "Unexpected error loading notes.";
      setErrorMessage(message);
      setStatus("error");
    }
  }, [page, limit]);

  useEffect(() => {
    load();
  }, [load, reloadToken]);

  const refetch = useCallback(() => setReloadToken((t) => t + 1), []);

  function handleSetLimit(next: number) {
    setLimit(next);
    setPage(1);
  }

  return {
    notes,
    page,
    limit,
    status,
    errorMessage,
    setPage: (p) => setPage(Math.max(1, p)),
    setLimit: handleSetLimit,
    refetch,
  };
}
