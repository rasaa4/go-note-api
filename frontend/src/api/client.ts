import { ApiError } from "./ApiError";

/**
 * Base URL of the Go backend. Configure via VITE_API_BASE_URL in .env
 * (see .env.example). Falls back to the backend's default port (8080)
 * for local development, per note-api/cmd/main.go.
 */
export const API_BASE_URL: string =
  import.meta.env.VITE_API_BASE_URL?.replace(/\/$/, "") ||
  "http://localhost:8080";

interface RequestOptions {
  method?: "GET" | "POST" | "PUT" | "DELETE";
  body?: unknown;
  query?: Record<string, string | number | undefined>;
}

function buildUrl(path: string, query?: RequestOptions["query"]): string {
  const url = new URL(`${API_BASE_URL}${path}`);
  if (query) {
    for (const [key, value] of Object.entries(query)) {
      if (value !== undefined) url.searchParams.set(key, String(value));
    }
  }
  return url.toString();
}

/**
 * The Go backend is inconsistent in how it reports errors (this is a
 * property of the backend, not something the frontend can normalize
 * away without changing note-api itself):
 *
 *  - Errors raised via handler.handleError (validation, not-found,
 *    internal) come back as JSON: {"error": "message"}.
 *  - Errors raised via the standard library's http.Error — used for
 *    malformed JSON bodies and unparsable id/page/limit query params,
 *    before the request ever reaches the service layer — come back as
 *    plain text with a trailing newline.
 *
 * This function handles both so the rest of the app only ever sees a
 * clean ApiError with a readable message.
 */
async function parseErrorBody(res: Response): Promise<string> {
  const raw = await res.text();
  try {
    const parsed = JSON.parse(raw) as { error?: string };
    if (parsed && typeof parsed.error === "string") return parsed.error;
  } catch {
    // Not JSON — this is the plain-text http.Error() path.
  }
  return raw.trim() || res.statusText || "Request failed";
}

/**
 * Performs a request against the note-api backend and returns the
 * parsed JSON body. Throws ApiError for any non-2xx response or for
 * network-level failures (backend down, CORS blocked, offline, etc).
 */
export async function apiRequest<T>(
  path: string,
  options: RequestOptions = {}
): Promise<T> {
  const { method = "GET", body, query } = options;
  const url = buildUrl(path, query);

  let res: Response;
  try {
    res = await fetch(url, {
      method,
      headers: body !== undefined ? { "Content-Type": "application/json" } : undefined,
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });
  } catch {
    throw new ApiError(
      "Could not reach the server. Check that the Go backend is running and that CORS is enabled.",
      0
    );
  }

  if (!res.ok) {
    const message = await parseErrorBody(res);
    throw new ApiError(message, res.status);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return (await res.json()) as T;
}
