import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Explore from './Explore';

// ── Mock docsConfig to keep test data small and predictable ───────────────────

vi.mock('../docsConfig', () => {
  const docsConfig = [
    {
      group: 'Getting Started',
      pages: [
        { id: 'introduction', title: 'Introduction', file: '/docs/introduction.md' },
        { id: 'quick-start', title: 'Quick Start', file: '/docs/quick-start.md' },
      ],
    },
    {
      group: 'Reference',
      pages: [
        { id: 'api-reference', title: 'API Reference', file: '/docs/api-reference.md' },
      ],
    },
  ];
  return { default: docsConfig };
});

// ── Mock react-syntax-highlighter (heavy dep, not needed in tests) ────────────

vi.mock('react-syntax-highlighter', () => ({
  Prism: ({ children }: { children: string }) => <pre>{children}</pre>,
}));
vi.mock('react-syntax-highlighter/dist/esm/styles/prism', () => ({
  oneDark: {},
}));

// ── Helper: stub fetch to return markdown ─────────────────────────────────────

const FIXTURE_MD = `# Introduction\n\n## Overview\n\nWelcome to go-initializer.\n\n## Installation\n\nFoo bar.`;

function stubFetch(text = FIXTURE_MD) {
  vi.stubGlobal('fetch', vi.fn(() =>
    Promise.resolve({ ok: true, text: () => Promise.resolve(text) } as Response),
  ));
}

// ── Rendering ─────────────────────────────────────────────────────────────────

describe('Explore – rendering', () => {
  beforeEach(() => stubFetch());
  afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); });

  it('renders both sidebar groups', async () => {
    const { container } = render(<Explore onBack={vi.fn()} />);
    const sidebar = container.querySelector('nav[aria-label="Documentation navigation"]')!;
    await waitFor(() => expect(within(sidebar as HTMLElement).getByText('Getting Started')).toBeInTheDocument());
    expect(within(sidebar as HTMLElement).getByText('Reference')).toBeInTheDocument();
  });

  it('renders all sidebar page links', async () => {
    const { container } = render(<Explore onBack={vi.fn()} />);
    const sidebar = container.querySelector('nav[aria-label="Documentation navigation"]')!;
    await waitFor(() => expect(within(sidebar as HTMLElement).getByText('Introduction')).toBeInTheDocument());
    expect(within(sidebar as HTMLElement).getByText('Quick Start')).toBeInTheDocument();
    expect(within(sidebar as HTMLElement).getByText('API Reference')).toBeInTheDocument();
  });

  it('shows skeleton while content is loading', () => {
    render(<Explore onBack={vi.fn()} />);
    expect(document.querySelector('[aria-busy="true"]')).toBeInTheDocument();
  });

  it('renders fetched markdown content after load', async () => {
    render(<Explore onBack={vi.fn()} />);
    await waitFor(() => expect(screen.getByText(/Welcome to go-initializer/)).toBeInTheDocument());
  });
});

// ── Page navigation via sidebar ───────────────────────────────────────────────

// Helper to scope queries to the sidebar only (avoids breadcrumb conflicts)
function getSidebar(container: HTMLElement) {
  return container.querySelector('nav[aria-label="Documentation navigation"]') as HTMLElement;
}

describe('Explore – sidebar navigation', () => {
  afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); vi.clearAllMocks(); });

  it('fetches new page when a different sidebar link is clicked', async () => {
    stubFetch('# Quick Start\n\nGet going fast.');
    const { container } = render(<Explore onBack={vi.fn()} />);
    const sidebar = getSidebar(container);
    await waitFor(() => expect(within(sidebar).getByText('Quick Start')).toBeInTheDocument());

    const user = userEvent.setup();
    await user.click(within(sidebar).getByText('Quick Start'));

    await waitFor(() => expect(vi.mocked(fetch)).toHaveBeenCalledWith('/docs/quick-start.md'));
  });
});

// ── Group collapse / expand ───────────────────────────────────────────────────

describe('Explore – group collapse', () => {
  beforeEach(() => stubFetch());
  afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); });

  it('collapses a group when its header is clicked', async () => {
    const { container } = render(<Explore onBack={vi.fn()} />);
    const sidebar = getSidebar(container);
    await waitFor(() => expect(within(sidebar).getByText('Getting Started')).toBeInTheDocument());

    const user = userEvent.setup();
    // Click the group header button to collapse
    await user.click(within(sidebar).getByRole('button', { name: /Getting Started/i }));

    // After collapse, the group pages div should have maxHeight:0
    await waitFor(() => {
      const groupPages = sidebar.querySelector('.explore-group:first-child .explore-group-pages') as HTMLElement;
      expect(groupPages?.style.maxHeight).toBe('0px');
    });
  });

  it('expands a collapsed group when its header is clicked again', async () => {
    const { container } = render(<Explore onBack={vi.fn()} />);
    const sidebar = getSidebar(container);
    await waitFor(() => expect(within(sidebar).getByText('Getting Started')).toBeInTheDocument());

    const user = userEvent.setup();
    const groupBtn = within(sidebar).getByRole('button', { name: /Getting Started/i });
    await user.click(groupBtn); // collapse
    await user.click(groupBtn); // expand

    await waitFor(() => expect(within(sidebar).getByText('Quick Start')).toBeInTheDocument());
  });
});

// ── Sidebar search filter ─────────────────────────────────────────────────────

describe('Explore – sidebar search', () => {
  beforeEach(() => stubFetch());
  afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); });

  it('filters pages by query (case-insensitive)', async () => {
    const { container } = render(<Explore onBack={vi.fn()} />);
    const sidebar = getSidebar(container);
    await waitFor(() => expect(within(sidebar).getByText('Introduction')).toBeInTheDocument());

    const user = userEvent.setup();
    const searchInput = screen.getByPlaceholderText(/Filter pages/i);
    await user.type(searchInput, 'api');

    await waitFor(() => {
      expect(within(sidebar).queryByText('Introduction')).not.toBeInTheDocument();
      expect(within(sidebar).getByText('API Reference')).toBeInTheDocument();
    });
  });

  it('shows all pages when search is cleared', async () => {
    const { container } = render(<Explore onBack={vi.fn()} />);
    const sidebar = getSidebar(container);
    await waitFor(() => expect(within(sidebar).getByText('Introduction')).toBeInTheDocument());

    const user = userEvent.setup();
    const searchInput2 = screen.getByPlaceholderText(/Filter pages/i);
    await user.type(searchInput2, 'api');
    await user.clear(searchInput2);

    await waitFor(() => {
      expect(within(sidebar).getByText('Introduction')).toBeInTheDocument();
      expect(within(sidebar).getByText('Quick Start')).toBeInTheDocument();
    });
  });
});

// ── Prev / Next navigation ────────────────────────────────────────────────────

describe('Explore – prev/next navigation', () => {
  beforeEach(() => stubFetch());
  afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); });

  it('does not render a Prev button on the first page', async () => {
    const { container } = render(<Explore onBack={vi.fn()} />);
    // Wait for content to load
    await waitFor(() => expect(container.querySelector('.docs-pager-btn--next')).toBeInTheDocument());
    // "Introduction" is first page, no prev button
    expect(container.querySelector('.docs-pager-btn--prev')).not.toBeInTheDocument();
  });

  it('renders Next button on the first page', async () => {
    const { container } = render(<Explore onBack={vi.fn()} />);
    await waitFor(() => expect(container.querySelector('.docs-pager-btn--next')).toBeInTheDocument());
  });

  it('navigates to next page when Next is clicked', async () => {
    const { container } = render(<Explore onBack={vi.fn()} />);
    const user = userEvent.setup();
    const nextBtn = await waitFor(() => {
      const btn = container.querySelector('.docs-pager-btn--next') as HTMLElement;
      expect(btn).toBeInTheDocument();
      return btn;
    });
    await user.click(nextBtn);
    await waitFor(() => expect(vi.mocked(fetch)).toHaveBeenCalledWith('/docs/quick-start.md'));
  });
});

// ── onBack prop ───────────────────────────────────────────────────────────────

describe('Explore – onBack', () => {
  beforeEach(() => stubFetch());
  afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); });

  it('calls onBack when the back button is clicked', async () => {
    const onBack = vi.fn();
    render(<Explore onBack={onBack} />);
    await waitFor(() => expect(screen.getByText(/Welcome to go-initializer/)).toBeInTheDocument());

    const user = userEvent.setup();
    // The back button has an aria-label containing "back" or similar
    const backBtn = screen.getByRole('button', { name: /back/i });
    await user.click(backBtn);
    expect(onBack).toHaveBeenCalledTimes(1);
  });
});

// ── ToC parsing ───────────────────────────────────────────────────────────────

describe('Explore – table of contents', () => {
  afterEach(() => { vi.restoreAllMocks(); vi.unstubAllGlobals(); });

  it('populates ToC from H2 headings in markdown', async () => {
    stubFetch('# Title\n\n## Section One\n\nContent.\n\n## Section Two\n\nMore content.');
    render(<Explore onBack={vi.fn()} />);
    // ToC links are <a> elements with class explore-toc-link
    await waitFor(() => expect(screen.getAllByText('Section One').length).toBeGreaterThan(0));
    expect(screen.getAllByText('Section Two').length).toBeGreaterThan(0);
  });
});
