import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import TerminalDemo from '../components/TerminalDemo';

const INSTALL_CMD = 'go install github.com/neo7337/go-initializer/cmd/goini@latest';

const DEMO_STEPS = [
  {
    command: 'goini list types',
    output: [
      '  microservice    — HTTP microservice with a chosen framework',
      '  simple-project  — Minimal Go project with a main package',
      '  cli-app         — Command-line application using Cobra',
      '  ai-agent        — ReAct-style LLM agent scaffold',
    ],
  },
  {
    command: 'goini list frameworks --type microservice',
    output: [
      '  gin     — High-performance HTTP framework',
      '  echo    — Minimalist, flexible web framework',
      '  fiber   — Express-inspired framework built on Fasthttp',
      '  chi     — Lightweight idiomatic router',
      '  gorilla — Classic toolkit-based router',
    ],
  },
  {
    command: 'goini new --type microservice --name my-svc --module github.com/acme/my-svc --framework gin --docker',
    output: [
      '  Generating microservice "my-svc"…',
      '  ✓  go.mod',
      '  ✓  main.go',
      '  ✓  internal/server/server.go',
      '  ✓  internal/server/handler.go',
      '  ✓  Makefile',
      '  ✓  Dockerfile',
      '  ✓  .gitignore',
      '',
      '  ✔ Project created at ./my-svc',
      '  → cd my-svc && go mod tidy && make run',
    ],
  },
];

const FLAGS = [
  { flag: '--name',        type: 'string',   desc: 'Project name' },
  { flag: '--module',      type: 'string',   desc: 'Go module path (e.g. github.com/acme/myapp)' },
  { flag: '--description', type: 'string',   desc: 'Short project description' },
  { flag: '--go-version',  type: 'string',   desc: 'Go version (e.g. 1.23)' },
  { flag: '--type',        type: 'string',   desc: 'Project type: microservice | simple-project | cli-app | ai-agent' },
  { flag: '--framework',   type: 'string',   desc: 'Framework name — must be valid for the selected project type' },
  { flag: '--addon',       type: 'string[]', desc: 'Repeatable: category=value (e.g. --addon cache=redis --addon database=gorm)' },
  { flag: '--docker',      type: 'bool',     desc: 'Generate a multi-stage production Dockerfile' },
  { flag: '--output',      type: 'string',   desc: 'Output directory (default: ./<name>)' },
];

const CLIPage: React.FC = () => {
  const navigate = useNavigate();
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(INSTALL_CMD).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  };

  return (
    <div className="cli-page">

      {/* ── Section 1: Hero strip ─────────────────────────────── */}
      <section className="cli-page-section">
        <h1 style={{
          fontSize: 'clamp(1.75rem, 4vw, 2.5rem)',
          fontWeight: 800,
          letterSpacing: '-0.02em',
          color: 'var(--color-text-primary)',
          marginBottom: '0.5rem',
        }}>
          goini CLI
        </h1>
        <p style={{
          fontSize: '1.0625rem',
          color: 'var(--color-text-secondary)',
          maxWidth: 580,
          lineHeight: 1.65,
          marginBottom: '1.75rem',
        }}>
          Scaffold production-ready Go projects from your terminal in seconds.
          Every option available in the web UI is also available as a flag — making{' '}
          <code style={{ fontFamily: 'monospace', color: 'var(--color-accent)', fontSize: '0.95em' }}>goini</code>{' '}
          fully scriptable in CI pipelines.
        </p>

        {/* Install command block */}
        <div className="terminal-window" style={{ maxWidth: 640 }}>
          <div className="terminal-titlebar">
            <span className="terminal-dot" style={{ background: '#ff5f57' }} />
            <span className="terminal-dot" style={{ background: '#febc2e' }} />
            <span className="terminal-dot" style={{ background: '#28c840' }} />
            <span style={{
              marginLeft: 10,
              fontSize: '0.75rem',
              color: 'var(--color-text-muted)',
              fontFamily: 'monospace',
            }}>
              install
            </span>
          </div>
          <div className="terminal-body" style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: '1rem' }}>
            <div>
              <span className="terminal-prompt">$ </span>
              <span className="terminal-cmd">{INSTALL_CMD}</span>
            </div>
            <button
              type="button"
              onClick={handleCopy}
              aria-label="Copy install command"
              title="Copy to clipboard"
              style={{
                flexShrink: 0,
                background: 'none',
                border: '1px solid var(--color-border)',
                borderRadius: 'var(--radius-sm)',
                color: copied ? 'var(--color-success)' : 'var(--color-text-muted)',
                cursor: 'pointer',
                padding: '0.25rem 0.6rem',
                fontSize: '0.75rem',
                fontFamily: 'monospace',
                transition: 'border-color 0.12s, color 0.12s',
                whiteSpace: 'nowrap',
              }}
            >
              {copied ? '✓ Copied' : 'Copy'}
            </button>
          </div>
        </div>
      </section>

      {/* ── Section 2: Terminal demo ───────────────────────────── */}
      <section className="cli-page-section">
        <h2>See it in action</h2>
        <TerminalDemo steps={DEMO_STEPS} title="goini — demo" />
      </section>

      {/* ── Section 3: Flags reference ────────────────────────── */}
      <section className="cli-page-section">
        <h2>goini new — flags</h2>
        <p style={{ color: 'var(--color-text-muted)', fontSize: '0.9rem', marginBottom: '1.25rem' }}>
          Any flag provided skips its corresponding interactive prompt, making the command fully scriptable.
        </p>
        <div style={{ overflowX: 'auto' }}>
          <table className="cli-flags-table">
            <thead>
              <tr>
                <th>Flag</th>
                <th>Type</th>
                <th>Description</th>
              </tr>
            </thead>
            <tbody>
              {FLAGS.map(({ flag, type, desc }) => (
                <tr key={flag}>
                  <td><code>{flag}</code></td>
                  <td><span style={{ color: 'var(--color-text-muted)', fontFamily: 'monospace', fontSize: '0.8125rem' }}>{type}</span></td>
                  <td>{desc}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* ── Section 4: CTA strip ──────────────────────────────── */}
      <section className="cli-page-section" style={{ textAlign: 'center', paddingTop: '1rem' }}>
        <p style={{ color: 'var(--color-text-secondary)', marginBottom: '1.25rem', fontSize: '1rem' }}>
          Prefer a visual interface? The web UI offers the same options with live previews.
        </p>
        <button
          type="button"
          onClick={() => navigate('/')}
          className="hero-btn-primary"
          style={{ fontSize: '0.9375rem' }}
        >
          Try it in the Web UI →
        </button>
      </section>

    </div>
  );
};

export default CLIPage;

