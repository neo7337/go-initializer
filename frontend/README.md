# go-initializer frontend

This package contains the React frontend for go-initializer. It provides the browser UI for generating Go projects, browsing the built-in documentation, and using the CLI installation guide.

## What it includes

- Landing page for the project
- Project generator flow at `/generate`
- Documentation explorer at `/docs`
- CLI installation and usage guide at `/cli`
- Light and dark theme support with local persistence
- API integration for project generation and metadata lookup

## Stack

- React 18
- TypeScript
- Vite 8
- React Router
- Radix Themes
- Tailwind CSS 4 and PostCSS
- Vitest + Testing Library

## Prerequisites

- Node.js 18 or newer
- npm
- The go-initializer backend running locally if you want to use the generator UI end-to-end

## Local development

Install dependencies:

```bash
npm install
```

Start the dev server:

```bash
npm start
```

The frontend runs on port `3000` by default.

## API configuration

The frontend uses `VITE_API_URL` to locate the backend API.

- Default API base URL: `http://localhost:8182`
- Override it by setting `VITE_API_URL` before starting Vite

Example:

```bash
VITE_API_URL=http://localhost:8182 npm start
```

## Available scripts

- `npm start` - start the Vite development server
- `npm run build` - build the production bundle
- `npm run preview` - preview the production build locally
- `npm test` - run the Vitest suite once
- `npm run test:watch` - run Vitest in watch mode
- `npm run coverage` - run tests with coverage reporting

## Testing

The frontend test setup uses:

- `vitest` with the `jsdom` environment
- `@testing-library/react`
- `@testing-library/jest-dom`
- coverage via `@vitest/coverage-v8`

Current global coverage thresholds in [vite.config.ts](/Users/adityakumar/code/Projects/go-initializer/frontend/vite.config.ts):

- lines: `80%`
- branches: `70%`

Run the full suite:

```bash
npm test
```

Run coverage:

```bash
npm run coverage
```

## Project structure

- `src/App.tsx` - application shell and routing
- `src/pages` - top-level route pages
- `src/components` - reusable UI components
- `src/hooks` - generator state and behavior hooks
- `src/service.ts` - API client and typed request helpers
- `src/test/setup.ts` - Vitest setup for the test environment

## Notes

- Unknown routes are redirected to `/`
- Theme preference is stored in `localStorage`
- The generator downloads the generated project as a ZIP from the backend
