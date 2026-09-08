/**
 * Thrown for any non-2xx response from the note-api backend, or for
 * network-level failures (backend unreachable, CORS rejection, etc).
 */
export class ApiError extends Error {
  /** HTTP status code, or 0 for a network-level failure (no response at all). */
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }

  get isNotFound(): boolean {
    return this.status === 404;
  }

  get isValidation(): boolean {
    return this.status === 400;
  }

  get isServerError(): boolean {
    return this.status >= 500;
  }

  get isNetworkError(): boolean {
    return this.status === 0;
  }
}
