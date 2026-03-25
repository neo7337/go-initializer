import React from 'react';
import { useNavigate } from 'react-router-dom';

/* ── Feature card icons ─────────────────────────────────────── */
const LayersIcon = () => (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
        <polygon points="12 2 2 7 12 12 22 7 12 2" />
        <polyline points="2 17 12 22 22 17" />
        <polyline points="2 12 12 17 22 12" />
    </svg>
);

const PlugIcon = () => (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
        <path d="M12 20v-6M9 20h6" />
        <path d="M7 10V4l2 2 3-3 3 3 2-2v6" />
        <rect x="5" y="10" width="14" height="4" rx="1" />
    </svg>
);

const PuzzleIcon = () => (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
        <path d="M19.439 7.85c-.049.322.059.648.289.878l1.568 1.568c.47.47.706 1.087.706 1.704s-.235 1.233-.706 1.704l-1.611 1.611a.98.98 0 0 1-.837.276c-.47-.07-.802-.48-.968-.925a2.501 2.501 0 1 0-3.214 3.214c.446.166.855.497.925.968a.979.979 0 0 1-.276.837l-1.61 1.61a2.404 2.404 0 0 1-3.408 0l-1.569-1.567a.877.877 0 0 0-.877-.29c-.493.074-.84.504-1.02.968a2.5 2.5 0 1 1-3.237-3.237c.464-.18.894-.527.967-1.02a.877.877 0 0 0-.289-.877l-1.568-1.568A2.402 2.402 0 0 1 1.998 12c0-.617.236-1.234.706-1.704L4.23 8.77c.24-.24.581-.353.917-.303.515.077.877.528 1.073 1.01a2.5 2.5 0 1 0 3.259-3.259c-.482-.196-.933-.558-1.01-1.073-.05-.336.062-.676.303-.917l1.525-1.525A2.402 2.402 0 0 1 12 2c.617 0 1.234.236 1.704.706l1.568 1.568c.23.23.556.338.877.29.493-.074.84-.504 1.02-.968a2.5 2.5 0 1 1 3.237 3.237c-.464.18-.894.527-.967 1.02Z" />
    </svg>
);

const TerminalIcon = () => (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
        <polyline points="4 17 10 11 4 5" />
        <line x1="12" y1="19" x2="20" y2="19" />
    </svg>
);

const BrainIcon = () => (
    <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
        <path d="M12 5a3 3 0 1 0-5.997.125 4 4 0 0 0-2.526 5.77 4 4 0 0 0 .556 6.588A4 4 0 1 0 12 18Z" />
        <path d="M12 5a3 3 0 1 1 5.997.125 4 4 0 0 1 2.526 5.77 4 4 0 0 1-.556 6.588A4 4 0 1 1 12 18Z" />
        <path d="M15 13a4.5 4.5 0 0 1-3-4 4.5 4.5 0 0 1-3 4" />
        <path d="M17.599 6.5a3 3 0 0 0 .399-1.375" />
        <path d="M6.003 5.125A3 3 0 0 0 6.401 6.5" />
        <path d="M3.477 10.896a4 4 0 0 1 .585-.396" />
        <path d="M19.938 10.5a4 4 0 0 1 .585.396" />
        <path d="M6 18a4 4 0 0 1-1.967-.516" />
        <path d="M19.967 17.484A4 4 0 0 1 18 18" />
    </svg>
);

const ZapIcon = () => (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
        <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
    </svg>
);

const WorkflowIcon = () => (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
        <rect x="3" y="3" width="6" height="6" rx="1" />
        <rect x="15" y="3" width="6" height="6" rx="1" />
        <rect x="9" y="15" width="6" height="6" rx="1" />
        <path d="M6 9v3a3 3 0 0 0 3 3h6a3 3 0 0 0 3-3V9" />
        <line x1="12" y1="12" x2="12" y2="15" />
    </svg>
);

const VectorIcon = () => (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
        <circle cx="6" cy="6" r="2" />
        <circle cx="18" cy="6" r="2" />
        <circle cx="12" cy="18" r="2" />
        <line x1="8" y1="6" x2="16" y2="6" />
        <line x1="7" y1="7.5" x2="11" y2="16.5" />
        <line x1="17" y1="7.5" x2="13" y2="16.5" />
    </svg>
);

/* ── Feature data ───────────────────────────────────────────── */
const FEATURES = [
    {
        icon: <LayersIcon />,
        title: '5 Project Types',
        desc: 'Simple, Microservice, API Server, CLI App & AI Agent — all production-ready.',
    },
    {
        icon: <PlugIcon />,
        title: '9 Frameworks',
        desc: 'Gin, Echo, Fiber, Chi, Golly, GoKit, Cobra, urfave/cli and more.',
    },
    {
        icon: <PuzzleIcon />,
        title: 'Addons & Docker',
        desc: 'Redis, Postgres (GORM / Ent), Zap, Logrus, pgvector — opt-in per project.',
    },
    {
        icon: <TerminalIcon />,
        title: 'CLI + Web UI',
        desc: 'Use the browser form or the goini CLI — same engine powers both.',
    },
];

/* ── Component ──────────────────────────────────────────────── */

const HeroSection: React.FC = () => {
    const navigate = useNavigate();

    return (
        <section className="hero-section" aria-label="Introduction">
            <p className="hero-eyebrow">Go project scaffolding, done right</p>

            <h1 className="hero-headline">
                Scaffold production-ready{' '}
                <span className="hero-headline-accent">Go</span>{' '}
                projects in seconds
            </h1>

            <p className="hero-sub">
                Choose a type, pick a framework, generate.{' '}
                From REST APIs to production-ready{' '}
                <span className="hero-sub-accent">AI agent workflows</span>{' '}
                — zero boilerplate.
            </p>

            {/* AI Agent Spotlight */}
            <div className="ai-spotlight">
                {/* Glow orb */}
                <div className="ai-spotlight-glow" aria-hidden="true" />

                <div className="ai-spotlight-header">
                    <div className="ai-spotlight-header-left">
                        <div className="ai-spotlight-icon-wrap">
                            <BrainIcon />
                        </div>
                        <div>
                            <div className="ai-spotlight-badge">
                                <span className="ai-spotlight-badge-dot" aria-hidden="true" />
                                New &nbsp;&middot;&nbsp; AI Agent
                            </div>
                            <strong className="ai-spotlight-title">
                                Build LLM-powered Go agents in seconds
                            </strong>
                            <p className="ai-spotlight-desc">
                                Generate a fully-wired AI agent project — provider client, tool loop, vector store —
                                ready to integrate with any orchestration framework or workflow automation platform.
                            </p>
                        </div>
                    </div>
                    <button className="ai-spotlight-cta" onClick={() => navigate('/generate')}>
                        Try AI Agent →
                    </button>
                </div>

                {/* 3-column capability grid */}
                <div className="ai-capability-grid">
                    <div className="ai-capability-item">
                        <span className="ai-capability-icon"><ZapIcon /></span>
                        <div>
                            <strong className="ai-capability-title">4 LLM Providers</strong>
                            <p className="ai-capability-desc">LangChainGo, OpenAI, Gemini, Ollama — swap with one line.</p>
                        </div>
                    </div>
                    <div className="ai-capability-item">
                        <span className="ai-capability-icon"><VectorIcon /></span>
                        <div>
                            <strong className="ai-capability-title">Vector Store RAG</strong>
                            <p className="ai-capability-desc">pgvector, Qdrant, or chromem — pre-wired retrieval pipeline.</p>
                        </div>
                    </div>
                    <div className="ai-capability-item">
                        <span className="ai-capability-icon"><WorkflowIcon /></span>
                        <div>
                            <strong className="ai-capability-title">Workflow Ready</strong>
                            <p className="ai-capability-desc">Drop into Temporal, Asynq, n8n, or any HTTP pipeline.</p>
                        </div>
                    </div>
                </div>

                {/* Provider chips */}
                <div className="ai-spotlight-chips">
                    <span className="ai-chip">LangChainGo</span>
                    <span className="ai-chip">OpenAI</span>
                    <span className="ai-chip">Gemini</span>
                    <span className="ai-chip">Ollama</span>
                    <span className="ai-chip ai-chip--dim">pgvector</span>
                    <span className="ai-chip ai-chip--dim">Qdrant</span>
                    <span className="ai-chip ai-chip--dim">chromem</span>
                </div>
            </div>

            <div className="hero-features">
                {FEATURES.map(f => (
                    <div className="hero-feature-card" key={f.title}>
                        <span className="hero-feature-icon">{f.icon}</span>
                        <div>
                            <strong className="hero-feature-title">{f.title}</strong>
                            <p className="hero-feature-desc">{f.desc}</p>
                        </div>
                    </div>
                ))}
            </div>

            <div className="hero-actions">
                <button className="hero-btn-primary" onClick={() => navigate('/generate')}>
                    Open Generator →
                </button>
                <button className="hero-btn-secondary" onClick={() => navigate('/cli')}>
                    View CLI Guide →
                </button>
            </div>
        </section>
    );
};

export default HeroSection;
