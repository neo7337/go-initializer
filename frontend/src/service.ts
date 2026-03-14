// src/service.ts

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8182';

/** Default request timeout in milliseconds. */
const DEFAULT_TIMEOUT_MS = 15_000;

// ── Error type ────────────────────────────────────────────────────────────────

/** Structured error thrown for non-2xx responses. */
export class ApiError extends Error {
  constructor(
    public readonly status: number,
    public readonly statusText: string,
    message: string,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

// ── Core request ──────────────────────────────────────────────────────────────

type RequestOptions = Omit<RequestInit, 'method' | 'body'> & {
  /** Override the default timeout (ms). Pass 0 to disable. */
  timeoutMs?: number;
};

async function request<T>(
  method: string,
  endpoint: string,
  body?: unknown,
  { timeoutMs = DEFAULT_TIMEOUT_MS, ...init }: RequestOptions = {},
): Promise<T> {
  const controller = new AbortController();
  const timerId =
    timeoutMs > 0 ? window.setTimeout(() => controller.abort(), timeoutMs) : undefined;

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...init.headers,
  };

  let response: Response;
  try {
    response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...init,
      method,
      headers,
      body: body !== undefined ? JSON.stringify(body) : undefined,
      signal: controller.signal,
    });
  } catch (err) {
    if ((err as Error).name === 'AbortError') {
      throw new ApiError(0, 'Timeout', `Request to ${endpoint} timed out after ${timeoutMs}ms`);
    }
    throw err;
  } finally {
    if (timerId !== undefined) window.clearTimeout(timerId);
  }

  if (!response.ok) {
    const message = await response.text().catch(() => response.statusText);
    throw new ApiError(response.status, response.statusText, message);
  }

  // 204 No Content — return undefined cast to T
  if (response.status === 204) return undefined as T;

  return response.json() as Promise<T>;
}

// ── Public HTTP helpers ───────────────────────────────────────────────────────

export const get = <T>(endpoint: string, options?: RequestOptions): Promise<T> =>
  request<T>('GET', endpoint, undefined, options);

export const post = <T>(endpoint: string, data: unknown, options?: RequestOptions): Promise<T> =>
  request<T>('POST', endpoint, data, options);

export const put = <T>(endpoint: string, data: unknown, options?: RequestOptions): Promise<T> =>
  request<T>('PUT', endpoint, data, options);

export const patch = <T>(endpoint: string, data: unknown, options?: RequestOptions): Promise<T> =>
  request<T>('PATCH', endpoint, data, options);

export const del = <T>(endpoint: string, options?: RequestOptions): Promise<T> =>
  request<T>('DELETE', endpoint, undefined, options);

// ── Domain helpers ────────────────────────────────────────────────────────────

export interface GenerateProjectPayload {
  projectType: string;
  name: string;
  moduleName: string;
  description: string;
  goVersion: string;
  framework: string;
  dockerSupport: boolean;
  /** Mapped to `selectedAddons` in the backend JSON. */
  selectedAddons: Record<string, string[]>;
}

/** Downloads the generated project as a zip Blob. */
export async function generateProject(data: GenerateProjectPayload): Promise<Blob> {
  const controller = new AbortController();
  const timerId = window.setTimeout(() => controller.abort(), DEFAULT_TIMEOUT_MS);

  let response: Response;
  try {
    response = await fetch(`${API_BASE_URL}/generate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
      signal: controller.signal,
    });
  } catch (err) {
    if ((err as Error).name === 'AbortError') {
      throw new ApiError(0, 'Timeout', `generateProject timed out after ${DEFAULT_TIMEOUT_MS}ms`);
    }
    throw err;
  } finally {
    window.clearTimeout(timerId);
  }

  if (!response.ok) {
    const message = await response.text().catch(() => response.statusText);
    throw new ApiError(response.status, response.statusText, message);
  }

  return response.blob();
}

export interface MetaData {
  supportedProjectTypes: Record<string, string>;
  supportedFrameworks: Record<string, Record<string, boolean>>;
  supportedGoVersions: Record<string, boolean>;
  supportedAddons: Record<string, Record<string, boolean>>;
}

/** Fetches supported project metadata from the API. */
export async function getMetaData(): Promise<MetaData> {
  return get<MetaData>('/meta');
}