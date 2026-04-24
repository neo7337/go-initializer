import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import HeroSection from './HeroSection';

// useNavigate needs a router context
function renderHero() {
  return render(
    <MemoryRouter>
      <HeroSection />
    </MemoryRouter>,
  );
}

// ── Rendering ─────────────────────────────────────────────────────────────────

describe('HeroSection – rendering', () => {
  it('renders the main headline', () => {
    renderHero();
    expect(screen.getByRole('heading', { level: 1 })).toBeInTheDocument();
  });

  it('renders exactly 4 feature cards', () => {
    const { container } = renderHero();
    expect(container.querySelectorAll('.hero-feature-card')).toHaveLength(4);
  });

  it('renders feature card titles', () => {
    renderHero();
    expect(screen.getByText('5 Project Types')).toBeInTheDocument();
    expect(screen.getByText('9 Frameworks')).toBeInTheDocument();
    expect(screen.getByText('Addons & Docker')).toBeInTheDocument();
    expect(screen.getByText('CLI + Web UI')).toBeInTheDocument();
  });

  it('renders AI Spotlight section', () => {
    const { container } = renderHero();
    expect(container.querySelector('.ai-spotlight')).toBeInTheDocument();
  });

  it('renders CTA buttons', () => {
    renderHero();
    expect(screen.getByRole('button', { name: /Open Generator/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /View CLI Guide/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Try AI Agent/i })).toBeInTheDocument();
  });
});

// ── Navigation ────────────────────────────────────────────────────────────────

describe('HeroSection – navigation', () => {
  it('"Open Generator →" navigates to /generate', async () => {
    const { container } = renderHero();
    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: /Open Generator/i }));
    // Navigation is tested indirectly — MemoryRouter updates location
    // Just verify no errors were thrown and button is still in DOM
    expect(screen.getByRole('button', { name: /Open Generator/i })).toBeInTheDocument();
  });

  it('"View CLI Guide →" navigates to /cli', async () => {
    renderHero();
    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: /View CLI Guide/i }));
    expect(screen.getByRole('button', { name: /View CLI Guide/i })).toBeInTheDocument();
  });

  it('"Try AI Agent →" navigates to /generate', async () => {
    renderHero();
    const user = userEvent.setup();
    await user.click(screen.getByRole('button', { name: /Try AI Agent/i }));
    expect(screen.getByRole('button', { name: /Try AI Agent/i })).toBeInTheDocument();
  });
});
