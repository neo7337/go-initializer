import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, act, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import CLIPage from './CLIPage';

// TerminalDemo runs animation timers — stub it to keep tests fast
vi.mock('../components/TerminalDemo', () => ({
  default: ({ steps, title }: { steps: unknown[]; title?: string }) => (
    <div data-testid="terminal-demo" data-steps={steps.length} data-title={title} />
  ),
}));

function renderCLI() {
  return render(
    <MemoryRouter>
      <CLIPage />
    </MemoryRouter>,
  );
}

// ── Rendering ─────────────────────────────────────────────────────────────────

describe('CLIPage – rendering', () => {
  it('renders the install command', () => {
    renderCLI();
    expect(
      screen.getByText(/go install github\.com\/neo7337\/go-initializer/),
    ).toBeInTheDocument();
  });

  it('renders the TerminalDemo component', () => {
    renderCLI();
    expect(screen.getByTestId('terminal-demo')).toBeInTheDocument();
  });

  it('renders TerminalDemo with 3 demo steps', () => {
    renderCLI();
    expect(screen.getByTestId('terminal-demo')).toHaveAttribute('data-steps', '3');
  });

  it('renders the flags table with 9 rows', () => {
    const { container } = renderCLI();
    const rows = container.querySelectorAll('tbody tr');
    expect(rows).toHaveLength(9);
  });

  it('renders the copy button', () => {
    renderCLI();
    expect(screen.getByRole('button', { name: /copy/i })).toBeInTheDocument();
  });
});

// ── Copy button ───────────────────────────────────────────────────────────────

describe('CLIPage – copy button', () => {
  let writeText: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    vi.useFakeTimers();
    writeText = vi.fn().mockResolvedValue(undefined);
    Object.defineProperty(navigator, 'clipboard', {
      value: { writeText },
      configurable: true,
      writable: true,
    });
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.restoreAllMocks();
  });

  it('calls clipboard.writeText with the install command', async () => {
    renderCLI();
    await act(async () => { fireEvent.click(screen.getByRole('button', { name: /copy/i })); });
    expect(writeText).toHaveBeenCalledWith(
      'go install github.com/neo7337/go-initializer/cmd/goini@latest',
    );
  });

  it('shows "Copied!" feedback after click', async () => {
    renderCLI();
    await act(async () => { fireEvent.click(screen.getByRole('button', { name: /copy/i })); });
    await act(async () => { await Promise.resolve(); });
    expect(screen.getByText(/Copied/i)).toBeInTheDocument();
  });

  it('reverts back to "Copy" after 2 s', async () => {
    renderCLI();
    await act(async () => { fireEvent.click(screen.getByRole('button', { name: /copy/i })); });
    await act(async () => { await Promise.resolve(); });
    act(() => vi.advanceTimersByTime(2000));
    expect(screen.queryByText(/Copied/i)).not.toBeInTheDocument();
  });
});

// ── Navigation ────────────────────────────────────────────────────────────────

describe('CLIPage – navigation', () => {
  it('renders a "Try it in the Web UI" link or button', () => {
    renderCLI();
    expect(screen.getByRole('button', { name: /Try it in the Web UI/i })).toBeInTheDocument();
  });
});
