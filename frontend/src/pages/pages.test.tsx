import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import HomePage from './HomePage';
import GeneratorPage from './GeneratorPage';

// Stub heavy child components
vi.mock('../components/HeroSection', () => ({
  default: () => <div data-testid="hero-section" />,
}));

vi.mock('../components/GeneratorForm', () => ({
  default: () => <div data-testid="generator-form" />,
}));

describe('HomePage', () => {
  it('renders HeroSection', () => {
    render(
      <MemoryRouter>
        <HomePage />
      </MemoryRouter>,
    );
    expect(screen.getByTestId('hero-section')).toBeInTheDocument();
  });
});

describe('GeneratorPage', () => {
  it('renders GeneratorForm', () => {
    render(
      <MemoryRouter>
        <GeneratorPage />
      </MemoryRouter>,
    );
    expect(screen.getByTestId('generator-form')).toBeInTheDocument();
  });
});
