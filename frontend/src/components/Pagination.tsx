interface PaginationProps {
  page: number;
  limit: number;
  resultCount: number;
  onPageChange: (page: number) => void;
  onLimitChange: (limit: number) => void;
}

const LIMIT_OPTIONS = [5, 10, 20, 50];

/**
 * The backend (note-api/repository.GetAll) computes offset = (page-1)*limit
 * and returns a plain Note[] with no total count or "has more" flag —
 * see README.md "Backend limitations". Without that metadata there is no
 * reliable way to know how many pages exist or to disable "Next" on the
 * true last page, so this is intentionally Previous/Next only, and "Next"
 * is disabled only when the current page returned fewer notes than the
 * requested limit (a reasonable but not perfectly certain signal).
 */
export function Pagination({
  page,
  limit,
  resultCount,
  onPageChange,
  onLimitChange,
}: PaginationProps) {
  const looksLikeLastPage = resultCount < limit;

  return (
    <div className="flex flex-col sm:flex-row items-center justify-between gap-3 mt-6 pt-4 border-t border-paper-line">
      <div className="flex items-center gap-2 text-sm text-ink-soft">
        <label htmlFor="page-size" className="font-mono text-xs text-ink-faint">
          per page
        </label>
        <select
          id="page-size"
          value={limit}
          onChange={(e) => onLimitChange(Number(e.target.value))}
          className="rounded-md border border-paper-line bg-paper-card px-2 py-1 text-sm text-ink focus:outline-none focus:ring-2 focus:ring-moss-300"
        >
          {LIMIT_OPTIONS.map((opt) => (
            <option key={opt} value={opt}>
              {opt}
            </option>
          ))}
        </select>
      </div>

      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={() => onPageChange(page - 1)}
          disabled={page <= 1}
          className="rounded-md border border-paper-line px-3 py-1.5 text-sm font-medium text-ink-soft hover:bg-paper disabled:opacity-40 disabled:hover:bg-transparent transition-colors"
        >
          Previous
        </button>
        <span className="font-mono text-sm text-ink-faint tabular-nums">page {page}</span>
        <button
          type="button"
          onClick={() => onPageChange(page + 1)}
          disabled={looksLikeLastPage}
          className="rounded-md border border-paper-line px-3 py-1.5 text-sm font-medium text-ink-soft hover:bg-paper disabled:opacity-40 disabled:hover:bg-transparent transition-colors"
        >
          Next
        </button>
      </div>
    </div>
  );
}
