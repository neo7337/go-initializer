import React, { useState, useEffect } from 'react';
import { Theme } from '@radix-ui/themes';
import Explore from './components/Explore';
import GeneratorForm from './components/GeneratorForm';

const SunIcon = () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
        <circle cx="12" cy="12" r="5" />
        <line x1="12" y1="1" x2="12" y2="3" />
        <line x1="12" y1="21" x2="12" y2="23" />
        <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
        <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
        <line x1="1" y1="12" x2="3" y2="12" />
        <line x1="21" y1="12" x2="23" y2="12" />
        <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
        <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
    </svg>
);

const MoonIcon = () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
        <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
    </svg>
);

const GitHubIcon = () => (
    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
        <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.012 8.012 0 0 0 16 8c0-4.42-3.58-8-8-8z" />
    </svg>
);

function getInitialTheme(): string {
    try {
        const saved = localStorage.getItem('theme');
        if (saved === 'light' || saved === 'dark') return saved;
    } catch { /* localStorage unavailable */ }
    if (window.matchMedia('(prefers-color-scheme: light)').matches) return 'light';
    return 'dark';
}

function App() {
    const [theme, setTheme] = useState(getInitialTheme);
    const [showExplore, setShowExplore] = useState(false);

    useEffect(() => {
        document.body.setAttribute('data-theme', theme);
        try { localStorage.setItem('theme', theme); } catch { /* ignore */ }
    }, [theme]);

    const toggleTheme = () => setTheme(prev => prev === 'light' ? 'dark' : 'light');

    return (
        <Theme appearance={theme === 'dark' ? 'dark' : 'light'}>
            <div className={`app-shell${showExplore ? ' app-shell--locked' : ''}`}>
                {/* Sticky header with backdrop-blur */}
                <header className="app-header">
                    <div className="header-inner">
                        <button
                            className="logo-wordmark"
                            onClick={() => setShowExplore(false)}
                            title="Go home"
                        >
                            <span className="logo-go">go</span><span className="logo-accent">initializer</span>
                        </button>

                        <div className="header-actions">
                            <button
                                className="icon-btn"
                                onClick={toggleTheme}
                                aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
                            >
                                {theme === 'dark' ? <SunIcon /> : <MoonIcon />}
                            </button>
                            <a
                                href="https://github.com/neo7337/go-initializer"
                                target="_blank"
                                rel="noopener noreferrer"
                                className="icon-btn"
                                aria-label="View source on GitHub"
                            >
                                <GitHubIcon />
                            </a>
                        </div>
                    </div>
                </header>

                {/* Main content */}
                <main className={`app-main${showExplore ? ' app-main--explore' : ''}`}>
                    {showExplore ? (
                        <Explore onBack={() => setShowExplore(false)} />
                    ) : (
                        <div className="main-centered">
                            <GeneratorForm onExplore={() => setShowExplore(true)} />
                        </div>
                    )}
                </main>

                {/* Slim footer */}
                <footer className="app-footer">
                    <div className="footer-inner">
                        <span className="footer-copy">&copy; {new Date().getFullYear()} Go Initializer</span>
                        <nav className="footer-links" aria-label="Footer navigation">
                            <a href="https://github.com/neo7337/go-initializer" target="_blank" rel="noopener noreferrer">GitHub</a>
                            <a href="https://github.com/neo7337/go-initializer/blob/main/LICENSE" target="_blank" rel="noopener noreferrer">License</a>
                        </nav>
                    </div>
                </footer>
            </div>
        </Theme>
    );
}

export default App;

