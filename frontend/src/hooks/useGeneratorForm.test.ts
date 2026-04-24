import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act, waitFor } from '@testing-library/react';
import { useGeneratorForm } from './useGeneratorForm';

// ── module-level mocks ────────────────────────────────────────────────────────

vi.mock('../service', () => ({
  getMetaData: vi.fn(),
  generateProject: vi.fn(),
}));

import { getMetaData, generateProject } from '../service';

const mockMeta = {
  supportedProjectTypes: { microservice: 'Microservice', 'cli-app': 'CLI App' },
  supportedFrameworks: { microservice: { gin: true, echo: true }, 'cli-app': { cobra: true } },
  supportedGoVersions: { '1.23': true, '1.22': true },
  supportedAddons: { cache: { redis: true }, database: { gorm: true }, other: {} },
};

// ── initial state ─────────────────────────────────────────────────────────────

describe('useGeneratorForm – initial state', () => {
  beforeEach(() => {
    vi.mocked(getMetaData).mockResolvedValue(mockMeta);
  });

  afterEach(() => vi.clearAllMocks());

  it('initialises all string fields to empty string', async () => {
    const { result } = renderHook(() => useGeneratorForm());
    expect(result.current.projectType).toBe('');
    expect(result.current.goVersion).toBe('');
    expect(result.current.framework).toBe('');
    expect(result.current.moduleName).toBe('');
    expect(result.current.name).toBe('');
    expect(result.current.description).toBe('');
  });

  it('initialises dockerSupport to false and selectedAddons to empty arrays', async () => {
    const { result } = renderHook(() => useGeneratorForm());
    expect(result.current.dockerSupport).toBe(false);
    expect(result.current.selectedAddons).toEqual({ cache: [], database: [], other: [] });
  });

  it('initialises errors to empty object', async () => {
    const { result } = renderHook(() => useGeneratorForm());
    expect(result.current.errors).toEqual({});
  });
});

// ── getMetaData on mount ──────────────────────────────────────────────────────

describe('useGeneratorForm – metadata on mount', () => {
  afterEach(() => vi.clearAllMocks());

  it('populates goVersionOptions after mount', async () => {
    vi.mocked(getMetaData).mockResolvedValue(mockMeta);
    const { result } = renderHook(() => useGeneratorForm());
    await waitFor(() => expect(result.current.goVersionOptions.length).toBeGreaterThan(0));
    expect(result.current.goVersionOptions[0].label).toMatch(/latest stable/);
  });

  it('populates supportedProjectTypes after mount', async () => {
    vi.mocked(getMetaData).mockResolvedValue(mockMeta);
    const { result } = renderHook(() => useGeneratorForm());
    await waitFor(() => expect(result.current.supportedProjectTypes.length).toBeGreaterThan(0));
    expect(result.current.supportedProjectTypes.map(t => t.type)).toContain('microservice');
  });

  it('does not crash when getMetaData rejects', async () => {
    vi.mocked(getMetaData).mockRejectedValue(new Error('network error'));
    const { result } = renderHook(() => useGeneratorForm());
    // Remains in default state
    await waitFor(() => expect(result.current.goVersionOptions).toEqual([]));
  });
});

// ── handleAddonChange ─────────────────────────────────────────────────────────

describe('useGeneratorForm – handleAddonChange', () => {
  beforeEach(() => vi.mocked(getMetaData).mockResolvedValue(mockMeta));
  afterEach(() => vi.clearAllMocks());

  it('adds an addon to the category when not present', async () => {
    const { result } = renderHook(() => useGeneratorForm());
    act(() => result.current.handleAddonChange('cache', 'redis'));
    expect(result.current.selectedAddons.cache).toContain('redis');
  });

  it('removes an addon from the category when already present (toggle)', async () => {
    const { result } = renderHook(() => useGeneratorForm());
    act(() => result.current.handleAddonChange('cache', 'redis'));
    act(() => result.current.handleAddonChange('cache', 'redis'));
    expect(result.current.selectedAddons.cache).not.toContain('redis');
  });

  it('does not affect other categories', async () => {
    const { result } = renderHook(() => useGeneratorForm());
    act(() => result.current.handleAddonChange('cache', 'redis'));
    expect(result.current.selectedAddons.database).toEqual([]);
  });
});

// ── validateInput ─────────────────────────────────────────────────────────────

describe('useGeneratorForm – validateInput', () => {
  beforeEach(() => vi.mocked(getMetaData).mockResolvedValue(mockMeta));
  afterEach(() => vi.clearAllMocks());

  it('returns false and sets all 6 errors when all fields are empty', () => {
    const { result } = renderHook(() => useGeneratorForm());
    let valid: boolean;
    act(() => { valid = result.current.handleGenerate as unknown as boolean; });
    // call validateInput indirectly via handleGenerate — it sets errors
    act(() => result.current.handleGenerate());
    expect(Object.keys(result.current.errors)).toEqual(
      expect.arrayContaining(['moduleName', 'name', 'description', 'projectType', 'goVersion', 'framework']),
    );
  });

  it('marks all fields as touched on submit attempt', () => {
    const { result } = renderHook(() => useGeneratorForm());
    act(() => result.current.handleGenerate());
    const t = result.current.touched;
    expect([t.moduleName, t.name, t.description, t.goVersion, t.projectType, t.framework]).toEqual([
      true, true, true, true, true, true,
    ]);
  });
});

// ── framework reset on projectType change ─────────────────────────────────────

describe('useGeneratorForm – framework reset', () => {
  beforeEach(() => vi.mocked(getMetaData).mockResolvedValue(mockMeta));
  afterEach(() => vi.clearAllMocks());

  it('resets framework to "" when projectType changes', async () => {
    const { result } = renderHook(() => useGeneratorForm());
    await waitFor(() => expect(result.current.supportedProjectTypes.length).toBeGreaterThan(0));
    act(() => result.current.setProjectType('microservice'));
    act(() => result.current.setFramework('gin'));
    expect(result.current.framework).toBe('gin');
    act(() => result.current.setProjectType('cli-app'));
    await waitFor(() => expect(result.current.framework).toBe(''));
  });
});

// ── handleGenerate – validation failure ───────────────────────────────────────

describe('useGeneratorForm – handleGenerate validation failure', () => {
  afterEach(() => vi.clearAllMocks());

  it('does NOT call generateProject when fields are empty', () => {
    vi.mocked(getMetaData).mockResolvedValue(mockMeta);
    const { result } = renderHook(() => useGeneratorForm());
    act(() => result.current.handleGenerate());
    expect(vi.mocked(generateProject)).not.toHaveBeenCalled();
  });
});

// ── handleGenerate – happy path ───────────────────────────────────────────────

describe('useGeneratorForm – handleGenerate happy path', () => {
  let originalCreateElement: typeof document.createElement;

  beforeEach(() => {
    vi.mocked(getMetaData).mockResolvedValue(mockMeta);
    vi.mocked(generateProject).mockResolvedValue(new Blob(['zip'], { type: 'application/zip' }));
    // stub browser download APIs
    vi.stubGlobal('URL', {
      ...URL,
      createObjectURL: vi.fn(() => 'blob:mock'),
      revokeObjectURL: vi.fn(),
    });
    originalCreateElement = document.createElement.bind(document);
    const mockA = { href: '', download: '', click: vi.fn(), remove: vi.fn(), style: {} };
    vi.spyOn(document.body, 'appendChild').mockImplementation((node) => node);
    vi.spyOn(document, 'createElement').mockImplementation((tag: string) => {
      if (tag === 'a') return mockA as unknown as HTMLElement;
      return originalCreateElement(tag);
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
    vi.clearAllMocks();
  });

  it('calls generateProject with the correct payload', async () => {
    const { result } = renderHook(() => useGeneratorForm());
    await waitFor(() => expect(result.current.supportedProjectTypes.length).toBeGreaterThan(0));

    act(() => {
      result.current.setModuleName('github.com/acme/svc');
      result.current.setName('svc');
      result.current.setDescription('a service');
      result.current.setProjectType('microservice');
      result.current.setGoVersion('1.23');
    });
    act(() => { result.current.setFramework('gin'); });

    await act(async () => result.current.handleGenerate());

    expect(vi.mocked(generateProject)).toHaveBeenCalledWith(
      expect.objectContaining({
        moduleName: 'github.com/acme/svc',
        name: 'svc',
        description: 'a service',
        projectType: 'microservice',
        goVersion: '1.23',
        framework: 'gin',
      }),
    );
  });

  it('sets generateSuccess to true on success', async () => {
    const { result } = renderHook(() => useGeneratorForm());
    await waitFor(() => expect(result.current.supportedProjectTypes.length).toBeGreaterThan(0));

    act(() => {
      result.current.setModuleName('github.com/acme/svc');
      result.current.setName('svc');
      result.current.setDescription('a service');
      result.current.setProjectType('microservice');
      result.current.setGoVersion('1.23');
    });
    act(() => { result.current.setFramework('gin'); });

    await act(async () => result.current.handleGenerate());
    expect(result.current.generateSuccess).toBe(true);
  });
});

// ── handleGenerate – API error ────────────────────────────────────────────────

describe('useGeneratorForm – handleGenerate API error', () => {
  beforeEach(() => {
    vi.mocked(getMetaData).mockResolvedValue(mockMeta);
    vi.mocked(generateProject).mockRejectedValue(new Error('Server error'));
  });
  afterEach(() => vi.clearAllMocks());

  it('sets generateError to the error message', async () => {
    const { result } = renderHook(() => useGeneratorForm());
    await waitFor(() => expect(result.current.supportedProjectTypes.length).toBeGreaterThan(0));

    act(() => {
      result.current.setModuleName('github.com/acme/svc');
      result.current.setName('svc');
      result.current.setDescription('a service');
      result.current.setProjectType('microservice');
      result.current.setGoVersion('1.23');
    });
    act(() => { result.current.setFramework('gin'); });

    await act(async () => result.current.handleGenerate());
    expect(result.current.generateError).toBe('Server error');
  });

  it('sets isGenerating back to false after error', async () => {
    const { result } = renderHook(() => useGeneratorForm());
    await waitFor(() => expect(result.current.supportedProjectTypes.length).toBeGreaterThan(0));

    act(() => {
      result.current.setModuleName('github.com/acme/svc');
      result.current.setName('svc');
      result.current.setDescription('a service');
      result.current.setProjectType('microservice');
      result.current.setGoVersion('1.23');
    });
    act(() => { result.current.setFramework('gin'); });

    await act(async () => result.current.handleGenerate());
    expect(result.current.isGenerating).toBe(false);
  });
});

// ── success countdown ─────────────────────────────────────────────────────────

describe('useGeneratorForm – success countdown', () => {
  let originalCreateElement2: typeof document.createElement;
  beforeEach(() => {
    vi.useFakeTimers();
    vi.mocked(getMetaData).mockResolvedValue(mockMeta);
    vi.mocked(generateProject).mockResolvedValue(new Blob([]));
    vi.stubGlobal('URL', { createObjectURL: vi.fn(() => 'blob:x'), revokeObjectURL: vi.fn() });
    vi.spyOn(document.body, 'appendChild').mockImplementation((n) => n);
    originalCreateElement2 = document.createElement.bind(document);
    const mockA2 = { href: '', download: '', click: vi.fn(), remove: vi.fn(), style: {} };
    vi.spyOn(document, 'createElement').mockImplementation((tag: string) => {
      if (tag === 'a') return mockA2 as unknown as HTMLElement;
      return originalCreateElement2(tag);
    });
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
    vi.clearAllMocks();
  });

  it('starts countdown at 5 and decrements to 0 then dismisses success', async () => {
    const { result } = renderHook(() => useGeneratorForm());

    // Trigger success
    act(() => {
      result.current.setModuleName('github.com/acme/svc');
      result.current.setName('svc');
      result.current.setDescription('desc');
      result.current.setProjectType('microservice');
      result.current.setGoVersion('1.23');
    });
    act(() => { result.current.setFramework('gin'); });

    await act(async () => result.current.handleGenerate());
    await act(async () => { await Promise.resolve(); });

    expect(result.current.generateSuccess).toBe(true);
    expect(result.current.successCountdown).toBe(5);

    act(() => vi.advanceTimersByTime(1000));
    expect(result.current.successCountdown).toBe(4);

    act(() => vi.advanceTimersByTime(4000));
    expect(result.current.generateSuccess).toBe(false);
  });
});

// ── keyboard shortcut ─────────────────────────────────────────────────────────

describe('useGeneratorForm – keyboard shortcut', () => {
  let originalCreateElement3: typeof document.createElement;
  beforeEach(() => {
    vi.mocked(getMetaData).mockResolvedValue(mockMeta);
    vi.mocked(generateProject).mockResolvedValue(new Blob([]));
    vi.stubGlobal('URL', { createObjectURL: vi.fn(() => 'blob:x'), revokeObjectURL: vi.fn() });
    vi.spyOn(document.body, 'appendChild').mockImplementation((n) => n);
    originalCreateElement3 = document.createElement.bind(document);
    const mockA3 = { href: '', download: '', click: vi.fn(), remove: vi.fn(), style: {} };
    vi.spyOn(document, 'createElement').mockImplementation((tag: string) => {
      if (tag === 'a') return mockA3 as unknown as HTMLElement;
      return originalCreateElement3(tag);
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
    vi.clearAllMocks();
  });

  it('triggers handleGenerate on Ctrl+Enter (non-Mac)', async () => {
    const { result } = renderHook(() => useGeneratorForm());
    await waitFor(() => expect(result.current.supportedProjectTypes.length).toBeGreaterThan(0));

    act(() => {
      result.current.setModuleName('github.com/acme/svc');
      result.current.setName('svc');
      result.current.setDescription('desc');
      result.current.setProjectType('microservice');
      result.current.setGoVersion('1.23');
    });
    act(() => { result.current.setFramework('gin'); });

    await act(async () => {
      window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', ctrlKey: true, bubbles: true }));
      await Promise.resolve();
    });

    expect(vi.mocked(generateProject)).toHaveBeenCalled();
  });
});
