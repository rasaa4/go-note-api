import type { ReactNode } from "react";

export function LoadingGrid() {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4" aria-hidden="true">
      {Array.from({ length: 6 }).map((_, i) => (
        <div
          key={i}
          className="rounded-lg border border-paper-line bg-paper-card p-5 h-40 animate-pulse"
        >
          <div className="h-4 w-2/3 bg-paper rounded" />
          <div className="mt-3 h-3 w-full bg-paper rounded" />
          <div className="mt-2 h-3 w-5/6 bg-paper rounded" />
          <div className="mt-6 h-3 w-1/4 bg-paper rounded" />
        </div>
      ))}
      <span className="sr-only">Loading notes…</span>
    </div>
  );
}

export function EmptyState({
  title,
  description,
  action,
}: {
  title: string;
  description: string;
  action?: ReactNode;
}) {
  return (
    <div className="rounded-lg border border-dashed border-paper-line bg-paper-card/60 px-6 py-16 text-center">
      <p className="font-display text-xl text-ink">{title}</p>
      <p className="mt-2 text-sm text-ink-soft max-w-sm mx-auto">{description}</p>
      {action && <div className="mt-5">{action}</div>}
    </div>
  );
}

export function ErrorState({
  message,
  onRetry,
}: {
  message: string;
  onRetry?: () => void;
}) {
  return (
    <div className="rounded-lg border border-rust-500/30 bg-rust-500/5 px-6 py-10 text-center">
      <p className="font-display text-lg text-rust-600">Something went wrong</p>
      <p className="mt-2 text-sm text-ink-soft max-w-md mx-auto">{message}</p>
      {onRetry && (
        <button
          type="button"
          onClick={onRetry}
          className="mt-4 rounded-md bg-ink px-4 py-2 text-sm font-medium text-white hover:bg-ink/90 transition-colors"
        >
          Try again
        </button>
      )}
    </div>
  );
}

export function Spinner({ className = "" }: { className?: string }) {
  return (
    <svg
      className={`animate-spin ${className}`}
      viewBox="0 0 24 24"
      fill="none"
      aria-hidden="true"
    >
      <circle cx="12" cy="12" r="9" stroke="currentColor" strokeWidth="3" opacity="0.25" />
      <path d="M21 12a9 9 0 0 0-9-9" stroke="currentColor" strokeWidth="3" strokeLinecap="round" />
    </svg>
  );
}
