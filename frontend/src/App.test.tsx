import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import App from './App';

// ── Heavy child components — render minimal stubs ─────────────────────────────

vi.mock('./pages/HomePage', () => ({ default: () => <div data-testid="home-page">Home</div> }));
vi.mock('./pages/GeneratorPage', () => ({ default: () => <div data-testid="generator-page">Generator</div> }));
vi.mock('./pages/CLIPage', () => ({ default: () => <div data-testid="cli-page">CLI</div> }));
vi.mock('./components/Explore', () => ({ default: ({ onBack }: { onBack: () => void }) => <div data-testid="explore-page"><button onClick={onBack}>back</button></div> }));

// App uses useLocation/useNavigate so it must be wrapped in a Router.
// We render inside MemoryRouter here instead of BrowserRouter.

function renderAt(path: string) {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <App />
    </MemoryRouter>,
  );
}

// ── Routing ───────────────────────────────────────────────────────────────────

describe('App – routing', () => {
  beforeEach(() => {
    // Avoid localStorage errors in jsdom
    vi.spyOn(Storage.prototype, 'getItem').mockReturnValue(null);
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {});
    // matchMedia stub
    vi.stubGlobal('matchMedia', (query: string) => ({
      matches: false,
      media: query,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    }));
  });

  afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); });

  it('renders HomePage at "/"', () => {
    renderAt('/');
    expect(screen.getByTestId('home-page')).toBeInTheDocument();
  });

  it('renders GeneratorPage at "/generate"', () => {
    renderAt('/generate');
    expect(screen.getByTestId('generator-page')).toBeInTheDocument();
  });

  it('renders Explore at "/docs"', () => {
    renderAt('/docs');
    expect(screen.getByTestId('explore-page')).toBeInTheDocument();
  });

  it('renders CLIPage at "/cli"', () => {
    renderAt('/cli');
    expect(screen.getByTestId('cli-page')).toBeInTheDocument();
  });

  it('redirects unknown paths to "/"', () => {
    renderAt('/does-not-exist');
    expect(screen.getByTestId('home-page')).toBeInTheDocument();
  });
});

// ── Theme toggle ──────────────────────────────────────────────────────────────

describe('App – theme toggle', () => {
  beforeEach(() => {
    vi.spyOn(Storage.prototype, 'getItem').mockReturnValue(null);
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {});
    vi.stubGlobal('matchMedia', (query: string) => ({
      matches: false,
      media: query,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    }));
  });

  afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); });

  it('applies data-theme attribute to document.body', () => {
    renderAt('/');
    // Default theme is 'dark' (no localStorage, no prefers-color-scheme)
    expect(document.body.getAttribute('data-theme')).toBe('dark');
  });

  it('toggles theme when the theme button is clicked', async () => {
    renderAt('/');
    expect(document.body.getAttribute('data-theme')).toBe('dark');

    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: /switch to light mode/i }));
    expect(document.body.getAttribute('data-theme')).toBe('light');
  });

  it('persists the new theme to localStorage', async () => {
    // The setItem mock from beforeEach is a no-op. We need to track real calls.
    // Replace the localStorage instance on window to intercept setItem.
    const realLS = window.localStorage;
    const setItemFn = vi.fn();
    Object.defineProperty(window, 'localStorage', {
      value: { getItem: vi.fn(() => null), setItem: setItemFn, removeItem: vi.fn(), clear: vi.fn(), length: 0 },
      configurable: true, writable: true,
    });
    renderAt('/');
    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: /switch to light mode/i }));
    await act(async () => {});
    expect(setItemFn).toHaveBeenCalledWith('theme', 'light');
    Object.defineProperty(window, 'localStorage', { value: realLS, configurable: true, writable: true });
  });
});

// ── getInitialTheme ───────────────────────────────────────────────────────────
// getInitialTheme is a lazy useState initializer — it runs on render inside the
// component, so mocking localStorage before rendering works correctly.

describe('App – getInitialTheme', () => {
  afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); document.body.removeAttribute('data-theme'); });

  it('uses "light" when localStorage.theme is "light"', () => {
    vi.stubGlobal('matchMedia', () => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn(), dispatchEvent: vi.fn() }));
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {});
    const realLocalStorage = window.localStorage;
    const getItemFn = vi.fn((key: string) => key === 'theme' ? 'light' : null);
    Object.defineProperty(window, 'localStorage', { value: { getItem: getItemFn, setItem: vi.fn(), removeItem: vi.fn(), clear: vi.fn(), length: 0 }, configurable: true, writable: true });
    renderAt('/');
    expect(document.body.getAttribute('data-theme')).toBe('light');
    Object.defineProperty(window, 'localStorage', { value: realLocalStorage, configurable: true, writable: true });
  });

  it('uses "dark" when localStorage.theme is "dark"', () => {
    vi.stubGlobal('matchMedia', () => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn(), dispatchEvent: vi.fn() }));
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation((key) => key === 'theme' ? 'dark' : null);
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {});
    renderAt('/');
    // Default when getItem returns null (or 'dark') — result should be 'dark'
    expect(document.body.getAttribute('data-theme')).toBe('dark');
  });

  it('defaults to "dark" when localStorage is empty and prefers-color-scheme does not match light', () => {
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => null);
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {});
    vi.stubGlobal('matchMedia', () => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn(), dispatchEvent: vi.fn() }));
    renderAt('/');
    expect(document.body.getAttribute('data-theme')).toBe('dark');
  });
});

// ── Header nav links ──────────────────────────────────────────────────────────

describe('App – header navigation', () => {
  beforeEach(() => {
    vi.spyOn(Storage.prototype, 'getItem').mockReturnValue(null);
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {});
    vi.stubGlobal('matchMedia', () => ({ matches: false, addEventListener: vi.fn(), removeEventListener: vi.fn(), dispatchEvent: vi.fn() }));
  });

  afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); });

  it('renders nav links: Generate, Docs, CLI', () => {
    renderAt('/');
    expect(screen.getByRole('link', { name: 'Generate' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Docs' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'CLI' })).toBeInTheDocument();
  });
});
