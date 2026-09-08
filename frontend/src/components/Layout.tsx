import { Link, Outlet } from "react-router-dom";

export function Layout() {
  return (
    <div className="min-h-screen bg-paper">
      <header className="border-b border-paper-line bg-paper-card">
        <div className="mx-auto max-w-5xl px-5 py-4 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-2.5">
            <span className="inline-flex h-7 w-7 items-center justify-center rounded-md bg-moss-500 text-white font-display text-sm">
              N
            </span>
            <span className="font-display text-lg text-ink">Notes</span>
          </Link>
          <Link
            to="/notes/new"
            className="rounded-md bg-ink px-3.5 py-2 text-sm font-medium text-white hover:bg-ink/90 transition-colors"
          >
            New note
          </Link>
        </div>
      </header>
      <main className="mx-auto max-w-5xl px-5 py-8">
        <Outlet />
      </main>
    </div>
  );
}
