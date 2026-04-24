import '@testing-library/jest-dom';

// jsdom does not implement IntersectionObserver — provide a no-op stub
(window as unknown as Record<string, unknown>).IntersectionObserver = class IntersectionObserver {
  root = null;
  rootMargin = '';
  thresholds = [];
  constructor(_: IntersectionObserverCallback, __?: IntersectionObserverInit) {}
  observe() {}
  unobserve() {}
  disconnect() {}
  takeRecords(): IntersectionObserverEntry[] { return []; }
} as unknown as typeof IntersectionObserver;
