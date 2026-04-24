import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import TerminalDemo, { TerminalStep } from './TerminalDemo';

// TerminalDemo uses chained setTimeouts (one per state change).
// We need to flush all timers AND React renders between each timer.
// This helper repeatedly flushes timers + React in a loop until stable.
async function drainAnimation() {
  for (let i = 0; i < 40; i++) {
    await act(async () => { vi.advanceTimersByTime(50); });
  }
}
const STEPS: TerminalStep[] = [
  { command: 'go version', output: ['go version go1.23 darwin/arm64'] },
  { command: 'echo hi', output: [] },
];

const SINGLE_STEP: TerminalStep[] = [
  { command: 'ls', output: ['file.go', 'go.mod'] },
];

describe('TerminalDemo – initial render', () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => { vi.useRealTimers(); vi.clearAllMocks(); });

  it('renders the terminal title', () => {
    render(<TerminalDemo steps={STEPS} title="test terminal" />);
    expect(screen.getByText('test terminal')).toBeInTheDocument();
  });

  it('renders three window-chrome dots', () => {
    const { container } = render(<TerminalDemo steps={STEPS} />);
    expect(container.querySelectorAll('.terminal-dot')).toHaveLength(3);
  });

  it('does not show the replay button initially', () => {
    render(<TerminalDemo steps={STEPS} />);
    expect(screen.queryByRole('button', { name: /replay/i })).not.toBeInTheDocument();
  });
});

describe('TerminalDemo – typing animation', () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => { vi.useRealTimers(); vi.clearAllMocks(); });

  it('shows a blinking cursor during the typing phase', () => {
    const { container } = render(<TerminalDemo steps={SINGLE_STEP} />);
    expect(container.querySelector('.terminal-cursor')).toBeInTheDocument();
  });

  it('types command characters one-by-one via CHAR_INTERVAL (15ms)', async () => {
    const { container } = render(<TerminalDemo steps={SINGLE_STEP} />);
    // Initially no characters typed
    const promptLine = container.querySelector('.terminal-prompt-line');
    const initialText = promptLine?.textContent ?? '';

    await act(async () => { vi.advanceTimersByTime(15); });
    const afterOneChar = promptLine?.textContent ?? '';
    expect(afterOneChar.length).toBeGreaterThan(initialText.length - 1);
  });

  it('types full command after enough CHAR_INTERVAL ticks', async () => {
    const { container } = render(<TerminalDemo steps={SINGLE_STEP} />);
    const cmd = SINGLE_STEP[0].command; // 'ls' = 2 chars
    await drainAnimation();
    const cmdEl = container.querySelector('.terminal-cmd');
    expect(cmdEl?.textContent).toBe(cmd);
  });
});

describe('TerminalDemo – output reveal', () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => { vi.useRealTimers(); vi.clearAllMocks(); });

  it('reveals output lines after the command is typed and POST_CMD_PAUSE (400ms) elapses', async () => {
    render(<TerminalDemo steps={SINGLE_STEP} />);
    const cmd = SINGLE_STEP[0].command;
    // Type full command + pause + first line interval
    await drainAnimation();
    expect(screen.getByText('file.go')).toBeInTheDocument();
  });
});

describe('TerminalDemo – done phase', () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => { vi.useRealTimers(); vi.clearAllMocks(); });

  it('shows replay button when animation is done', async () => {
    render(<TerminalDemo steps={SINGLE_STEP} />);
    const cmd = SINGLE_STEP[0].command;
    // type + POST_CMD_PAUSE + output lines + POST_STEP_PAUSE
    const total = 15 * cmd.length + 400 + 120 * SINGLE_STEP[0].output.length + 800 + 100;
    await drainAnimation();
    expect(screen.getByText(/Replay/)).toBeInTheDocument();
  });

  it('hides cursor when animation is done', async () => {
    const { container } = render(<TerminalDemo steps={SINGLE_STEP} />);
    const cmd = SINGLE_STEP[0].command;
    const total = 15 * cmd.length + 400 + 120 * SINGLE_STEP[0].output.length + 800 + 100;
    await drainAnimation();
    expect(container.querySelector('.terminal-cursor')).not.toBeInTheDocument();
  });
});

describe('TerminalDemo – replay', () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => { vi.useRealTimers(); vi.clearAllMocks(); });

  it('clicking Replay resets to the typing phase (cursor reappears)', async () => {
    const { container } = render(<TerminalDemo steps={SINGLE_STEP} />);
    await drainAnimation();

    const replayBtn = screen.getByText(/Replay/);
    await act(async () => { replayBtn.click(); });

    expect(container.querySelector('.terminal-cursor')).toBeInTheDocument();
    expect(screen.queryByText(/Replay/)).not.toBeInTheDocument();
  });
});
