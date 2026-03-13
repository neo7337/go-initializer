// This file tells TypeScript to treat CSS imports as modules (no type errors for side-effect imports)
declare module '*.css';

// CRA injects REACT_APP_* variables into process.env at build time.
// This declaration allows TypeScript to resolve them without @types/node.
declare const process: {
  env: {
    NODE_ENV: 'development' | 'production' | 'test';
    REACT_APP_API_URL?: string;
    [key: string]: string | undefined;
  };
};
