// This file tells TypeScript to treat CSS imports as modules (no type errors for side-effect imports)
declare module '*.css';

// Vite exposes env vars through import.meta.env for VITE_* prefixed variables.
interface ImportMetaEnv {
  readonly VITE_API_URL?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
