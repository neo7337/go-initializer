import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { ApiError, get, post, generateProject, getMetaData } from './service';

// ── helpers ───────────────────────────────────────────────────────────────────

function makeResponse(status: number, body: unknown, ok = status >= 200 && status < 300): Response {
  const bodyStr = typeof body === 'string' ? body : JSON.stringify(body);
  return {
    ok,
    status,
    statusText: status === 200 ? 'OK' : status === 204 ? 'No Content' : 'Error',
    json: () => Promise.resolve(body),
    text: () => Promise.resolve(bodyStr),
    blob: () => Promise.resolve(new Blob([bodyStr])),
    headers: new Headers(),
  } as unknown as Response;
}

// ── ApiError ──────────────────────────────────────────────────────────────────

describe('ApiError', () => {
  it('is an instance of Error', () => {
    const err = new ApiError(404, 'Not Found', 'resource not found');
    expect(err).toBeInstanceOf(Error);
  });

  it('exposes status, statusText, and message', () => {
    const err = new ApiError(500, 'Internal Server Error', 'something broke');
    expect(err.status).toBe(500);
    expect(err.statusText).toBe('Internal Server Error');
    expect(err.message).toBe('something broke');
  });

  it('has name "ApiError"', () => {
    const err = new ApiError(400, 'Bad Request', 'invalid input');
    expect(err.name).toBe('ApiError');
  });
});

// ── request (via get / post) ──────────────────────────────────────────────────

describe('request', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn());
    // stub window.setTimeout/clearTimeout so AbortController tests work in jsdom
    vi.stubGlobal('window', {
      ...window,
      setTimeout: globalThis.setTimeout,
      clearTimeout: globalThis.clearTimeout,
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  it('returns parsed JSON for a 200 response', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(makeResponse(200, { hello: 'world' }));
    const result = await get<{ hello: string }>('/test');
    expect(result).toEqual({ hello: 'world' });
  });

  it('returns undefined for 204 No Content', async () => {
    const noContent = { ...makeResponse(200, null), ok: true, status: 204 } as unknown as Response;
    vi.mocked(fetch).mockResolvedValueOnce(noContent);
    const result = await get<undefined>('/test');
    expect(result).toBeUndefined();
  });

  it('throws ApiError for non-2xx responses', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(makeResponse(404, 'not found', false));
    await expect(get('/missing')).rejects.toBeInstanceOf(ApiError);
  });

  it('ApiError carries the correct status from a non-2xx response', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(makeResponse(422, 'validation error', false));
      let err: unknown;
      try {
        await get('/bad');
      } catch (error) {
        err = error;
      }
      expect(err).toBeInstanceOf(ApiError);
      if (!(err instanceof ApiError)) {
        throw err;
      }
    expect(err.status).toBe(422);
  });

  it('throws ApiError(0, Timeout) when the request is aborted', async () => {
    vi.mocked(fetch).mockRejectedValueOnce(Object.assign(new Error('aborted'), { name: 'AbortError' }));
      let err: unknown;
      try {
        await get('/slow', { timeoutMs: 1 });
      } catch (error) {
        err = error;
      }
      expect(err).toBeInstanceOf(ApiError);
      if (!(err instanceof ApiError)) {
        throw err;
      }
    expect(err).toBeInstanceOf(ApiError);
    expect(err.status).toBe(0);
    expect(err.statusText).toBe('Timeout');
  });

  it('sends POST with JSON body and Content-Type header', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(makeResponse(200, { ok: true }));
    await post('/submit', { name: 'test' });
    const [, init] = vi.mocked(fetch).mock.calls[0];
    expect((init as RequestInit).method).toBe('POST');
    expect(JSON.parse((init as RequestInit).body as string)).toEqual({ name: 'test' });
    expect(((init as RequestInit).headers as Record<string, string>)['Content-Type']).toBe(
      'application/json',
    );
  });
});

// ── generateProject ───────────────────────────────────────────────────────────

describe('generateProject', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn());
    vi.stubGlobal('window', {
      ...window,
      setTimeout: globalThis.setTimeout,
      clearTimeout: globalThis.clearTimeout,
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  const payload = {
    projectType: 'microservice',
    name: 'my-svc',
    moduleName: 'github.com/acme/my-svc',
    description: 'test service',
    goVersion: '1.23',
    framework: 'gin',
    dockerSupport: false,
    selectedAddons: {},
  };

  it('POSTs to /api/generate and returns a Blob', async () => {
    const mockBlob = new Blob(['zip content'], { type: 'application/zip' });
    const mockResponse = {
      ok: true,
      status: 200,
      statusText: 'OK',
      blob: () => Promise.resolve(mockBlob),
      text: () => Promise.resolve(''),
    } as unknown as Response;
    vi.mocked(fetch).mockResolvedValueOnce(mockResponse);

    const result = await generateProject(payload);
    expect(result).toBeInstanceOf(Blob);

    const [url, init] = vi.mocked(fetch).mock.calls[0];
    expect(String(url)).toMatch(/\/api\/generate$/);
    expect((init as RequestInit).method).toBe('POST');
  });

  it('throws ApiError on non-2xx from /api/generate', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(makeResponse(500, 'server error', false));
    await expect(generateProject(payload)).rejects.toBeInstanceOf(ApiError);
  });

  it('throws ApiError(0, Timeout) when generateProject request times out', async () => {
    vi.mocked(fetch).mockRejectedValueOnce(
      Object.assign(new Error('aborted'), { name: 'AbortError' }),
    );
    const err = await generateProject(payload).catch(e => e);
    expect(err.status).toBe(0);
    expect(err.statusText).toBe('Timeout');
  });
});

// ── getMetaData ───────────────────────────────────────────────────────────────

describe('getMetaData', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn());
    vi.stubGlobal('window', {
      ...window,
      setTimeout: globalThis.setTimeout,
      clearTimeout: globalThis.clearTimeout,
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  const mockMeta = {
    supportedProjectTypes: { microservice: 'Microservice' },
    supportedFrameworks: { microservice: { gin: true } },
    supportedGoVersions: { '1.23': true },
    supportedAddons: { cache: { redis: true } },
  };

  it('GETs /api/meta and returns MetaData', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(makeResponse(200, mockMeta));
    const result = await getMetaData();
    expect(result).toEqual(mockMeta);
    const [url] = vi.mocked(fetch).mock.calls[0];
    expect(String(url)).toMatch(/\/api\/meta$/);
  });

  it('throws ApiError on network failure', async () => {
    vi.mocked(fetch).mockResolvedValueOnce(makeResponse(503, 'unavailable', false));
    await expect(getMetaData()).rejects.toBeInstanceOf(ApiError);
  });
});
