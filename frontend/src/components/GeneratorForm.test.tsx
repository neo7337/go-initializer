import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import GeneratorForm from './GeneratorForm';

// ── Mock the hook so tests control state deterministically ────────────────────

vi.mock('../hooks/useGeneratorForm', () => ({ useGeneratorForm: vi.fn() }));

import { useGeneratorForm } from '../hooks/useGeneratorForm';

function buildHookState(overrides: Partial<ReturnType<typeof useGeneratorForm>> = {}): ReturnType<typeof useGeneratorForm> {
  return {
    dockerSupport: false,
    setDockerSupport: vi.fn(),
    selectedAddons: { cache: [], database: [], other: [] },
    handleAddonChange: vi.fn(),
    projectType: '',
    setProjectType: vi.fn(),
    goVersion: '',
    setGoVersion: vi.fn(),
    framework: '',
    setFramework: vi.fn(),
    moduleName: '',
    setModuleName: vi.fn(),
    name: '',
    setName: vi.fn(),
    description: '',
    setDescription: vi.fn(),
    touched: { moduleName: false, name: false, description: false, goVersion: false, projectType: false, framework: false },
    setTouched: vi.fn(),
    errors: {},
    setErrors: vi.fn(),
    goVersionOptions: [
      { version: '1.23', label: '1.23 (latest stable)' },
      { version: '1.22', label: '1.22' },
    ],
    supportedProjectTypes: [
      { type: 'microservice', label: 'Microservice' },
      { type: 'ai-agent', label: 'AI Agent' },
    ],
    currentFrameworkOptions: [],
    addonOptions: {
      cache: [{ value: 'redis', label: 'Redis' }],
      database: [{ value: 'gorm', label: 'Gorm' }],
    },
    generateError: null,
    setGenerateError: vi.fn(),
    generateSuccess: false,
    setGenerateSuccess: vi.fn(),
    successCountdown: 5,
    isGenerating: false,
    isMac: false,
    handleGenerate: vi.fn(),
    ...overrides,
  };
}

// ── Rendering ─────────────────────────────────────────────────────────────────

describe('GeneratorForm – rendering', () => {
  beforeEach(() => vi.mocked(useGeneratorForm).mockReturnValue(buildHookState()));
  afterEach(() => vi.clearAllMocks());

  it('renders all 6 numbered section labels', () => {
    render(<GeneratorForm />);
    ['01', '02', '03', '04', '05', '06'].forEach(n =>
      expect(screen.getByText(n)).toBeInTheDocument(),
    );
  });

  it('renders go version pill buttons', () => {
    render(<GeneratorForm />);
    expect(screen.getByText('1.23 (latest stable)')).toBeInTheDocument();
    expect(screen.getByText('1.22')).toBeInTheDocument();
  });

  it('renders project type pill buttons', () => {
    render(<GeneratorForm />);
    expect(screen.getByText('Microservice')).toBeInTheDocument();
    expect(screen.getByText('AI Agent')).toBeInTheDocument();
  });

  it('renders addon chip buttons', () => {
    render(<GeneratorForm />);
    expect(screen.getByText('Redis')).toBeInTheDocument();
    expect(screen.getByText('Gorm')).toBeInTheDocument();
  });

  it('renders project detail text inputs', () => {
    render(<GeneratorForm />);
    expect(screen.getByLabelText(/Module Name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Project Name/i)).toBeInTheDocument();
  });
});

// ── aria-pressed on PillButton ────────────────────────────────────────────────

describe('GeneratorForm – PillButton aria-pressed', () => {
  afterEach(() => vi.clearAllMocks());

  it('selected version pill has aria-pressed="true"', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ goVersion: '1.23' }));
    render(<GeneratorForm />);
    const btn = screen.getByRole('button', { name: '1.23 (latest stable)' });
    expect(btn).toHaveAttribute('aria-pressed', 'true');
  });

  it('unselected version pill has aria-pressed="false"', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ goVersion: '1.23' }));
    render(<GeneratorForm />);
    const btn = screen.getByRole('button', { name: '1.22' });
    expect(btn).toHaveAttribute('aria-pressed', 'false');
  });
});

// ── onInteract callback ───────────────────────────────────────────────────────

describe('GeneratorForm – onInteract callback', () => {
  afterEach(() => vi.clearAllMocks());

  it('fires onInteract once on the first pill click', async () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState());
    const onInteract = vi.fn();
    render(<GeneratorForm onInteract={onInteract} />);
    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: '1.22' }));
    expect(onInteract).toHaveBeenCalledTimes(1);
  });

  it('does NOT fire onInteract a second time on subsequent interaction', async () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState());
    const onInteract = vi.fn();
    render(<GeneratorForm onInteract={onInteract} />);
    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: '1.22' }));
    await user.click(screen.getByRole('button', { name: '1.23 (latest stable)' }));
    expect(onInteract).toHaveBeenCalledTimes(1);
  });
});

// ── ai-agent type banner ──────────────────────────────────────────────────────

describe('GeneratorForm – ai-agent project type', () => {
  afterEach(() => vi.clearAllMocks());

  it('shows ai-agent-form-banner when projectType is ai-agent', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ projectType: 'ai-agent' }));
    const { container } = render(<GeneratorForm />);
    expect(container.querySelector('.ai-agent-form-banner')).toBeInTheDocument();
  });

  it('does NOT show ai-agent-form-banner for other project types', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ projectType: 'microservice' }));
    const { container } = render(<GeneratorForm />);
    expect(container.querySelector('.ai-agent-form-banner')).not.toBeInTheDocument();
  });

  it('relabels section 03 to "LLM Provider" when projectType is ai-agent', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ projectType: 'ai-agent' }));
    render(<GeneratorForm />);
    expect(screen.getByText('LLM Provider')).toBeInTheDocument();
  });
});

// ── Field error display ───────────────────────────────────────────────────────

describe('GeneratorForm – field errors', () => {
  afterEach(() => vi.clearAllMocks());

  it('shows field error when errors.goVersion is set and touched', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(
      buildHookState({
        errors: { goVersion: 'Go Version is required.' },
        touched: { moduleName: true, name: true, description: true, goVersion: true, projectType: true, framework: true },
      }),
    );
    render(<GeneratorForm />);
    expect(screen.getByText('Go Version is required.')).toBeInTheDocument();
  });

  it('shows field error when errors.projectType is set and touched', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(
      buildHookState({
        errors: { projectType: 'Project Type is required.' },
        touched: { moduleName: true, name: true, description: true, goVersion: true, projectType: true, framework: true },
      }),
    );
    render(<GeneratorForm />);
    expect(screen.getByText('Project Type is required.')).toBeInTheDocument();
  });

  it('does NOT show field error when field is not touched', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(
      buildHookState({
        errors: { goVersion: 'Go Version is required.' },
        touched: { moduleName: false, name: false, description: false, goVersion: false, projectType: false, framework: false },
      }),
    );
    render(<GeneratorForm />);
    expect(screen.queryByText('Go Version is required.')).not.toBeInTheDocument();
  });
});

// ── ErrorBanner ───────────────────────────────────────────────────────────────

describe('GeneratorForm – ErrorBanner', () => {
  afterEach(() => vi.clearAllMocks());

  it('renders ErrorBanner with role="alert" when generateError is set', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ generateError: 'API failed' }));
    render(<GeneratorForm />);
    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.getByText('API failed')).toBeInTheDocument();
  });

  it('calls setGenerateError(null) when dismiss button is clicked', async () => {
    const setGenerateError = vi.fn();
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ generateError: 'oops', setGenerateError }));
    render(<GeneratorForm />);
    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: /dismiss error/i }));
    expect(setGenerateError).toHaveBeenCalledWith(null);
  });
});

// ── SuccessBanner ─────────────────────────────────────────────────────────────

describe('GeneratorForm – SuccessBanner', () => {
  afterEach(() => vi.clearAllMocks());

  it('renders SuccessBanner with role="status" when generateSuccess is true', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ generateSuccess: true, successCountdown: 3 }));
    render(<GeneratorForm />);
    expect(screen.getByRole('status')).toBeInTheDocument();
  });

  it('displays the countdown value', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ generateSuccess: true, successCountdown: 3 }));
    render(<GeneratorForm />);
    expect(screen.getByText('3s')).toBeInTheDocument();
  });

  it('calls setGenerateSuccess(false) when dismiss button clicked', async () => {
    const setGenerateSuccess = vi.fn();
    vi.mocked(useGeneratorForm).mockReturnValue(
      buildHookState({ generateSuccess: true, successCountdown: 5, setGenerateSuccess }),
    );
    render(<GeneratorForm />);
    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: /dismiss success/i }));
    expect(setGenerateSuccess).toHaveBeenCalledWith(false);
  });
});

// ── Generate button ───────────────────────────────────────────────────────────

describe('GeneratorForm – Generate button', () => {
  afterEach(() => vi.clearAllMocks());

  it('is disabled while isGenerating is true', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ isGenerating: true, isMac: false }));
    render(<GeneratorForm />);
    // aria-label is always "Generate project (Ctrl+↵)" or "Generate project (⌘↵)"
    const btn = screen.getByRole('button', { name: /Generate project/i });
    expect(btn).toBeDisabled();
  });

  it('is enabled when isGenerating is false', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ isGenerating: false, isMac: false }));
    render(<GeneratorForm />);
    const btn = screen.getByRole('button', { name: /Generate project/i });
    expect(btn).not.toBeDisabled();
  });

  it('shows Ctrl+↵ hint on non-Mac', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ isMac: false }));
    render(<GeneratorForm />);
    expect(screen.getByText(/Ctrl/)).toBeInTheDocument();
  });

  it('shows ⌘↵ hint on Mac', () => {
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ isMac: true }));
    render(<GeneratorForm />);
    expect(screen.getByText(/⌘/)).toBeInTheDocument();
  });

  it('calls handleGenerate when button is clicked', async () => {
    const handleGenerate = vi.fn();
    vi.mocked(useGeneratorForm).mockReturnValue(buildHookState({ handleGenerate, isMac: false }));
    render(<GeneratorForm />);
    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: /Generate project/i }));
    expect(handleGenerate).toHaveBeenCalledTimes(1);
  });
});
