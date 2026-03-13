import React, { useState, useEffect } from 'react';
import './App.css';
import { Theme, Button, Flex, Heading, Text } from '@radix-ui/themes';
import Explore from './components/Explore';
import GeneratorForm from './components/GeneratorForm';


function App() {
    const [theme, setTheme] = useState('dark');
    const [showExplore, setShowExplore] = useState(false);

    useEffect(() => {
        document.body.setAttribute('data-theme', theme);
    }, [theme]);

    const toggleTheme = () => setTheme(prev => prev === 'light' ? 'dark' : 'light');

    return (
        <Theme appearance={theme === 'dark' ? 'dark' : 'light'}>
            <div className="App" style={{ height: '100vh', display: 'flex', flexDirection: 'column', background: 'var(--background)', color: 'var(--text)', transition: 'background 0.3s, color 0.3s', overflow: showExplore ? 'hidden' : 'auto' }}>
                {/* Header */}
                <header style={{ background: 'var(--navbar-bg)', boxShadow: '0 2px 8px rgba(0,0,0,0.04)', padding: '1rem 2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', color: 'var(--navbar-text)' }}>
                    <Flex align="center" justify="between" style={{ width: '100%' }}>
                        <Flex align="center">
                            <Heading
                                size="6"
                                weight="bold"
                                style={{ color: 'var(--navbar-text)', letterSpacing: 0.5, cursor: 'pointer', userSelect: 'none' }}
                                onClick={() => setShowExplore(false)}
                                title="Go home"
                            >
                                go<span style={{ color: '#ffd700' }}>initializer</span>
                            </Heading>
                        </Flex>
                        <Button variant="ghost" color="gray" radius="full" size="3" onClick={toggleTheme} title="Toggle theme">
                            {theme === 'light' ? '🌙' : '☀️'}
                        </Button>
                    </Flex>
                </header>

                {/* Main Content */}
                <main style={showExplore ? { flex: 1, background: 'var(--content-bg)', color: 'var(--text)', display: 'flex', flexDirection: 'column', overflow: 'hidden', minHeight: 0 } : { flex: 1, padding: '2.5rem 1rem', background: 'var(--content-bg)', color: 'var(--text)', display: 'grid', gridTemplateColumns: '1fr', gap: 32, maxWidth: 1200, width: '100%', margin: '0 auto' }}>
                    {showExplore ? (
                        <Explore onBack={() => setShowExplore(false)} />
                    ) : (
                        <GeneratorForm onExplore={() => setShowExplore(true)} />
                    )}
                </main>

                {/* Footer */}
                <footer style={{ background: 'var(--footer-bg)', color: 'var(--footer-text)', boxShadow: '0 -2px 8px rgba(0,0,0,0.04)', padding: '1rem 2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 'auto' }}>
                    <Flex align="center" justify="between" style={{ width: '100%' }}>
                        <span style={{ display: 'flex', alignItems: 'center' }}>
                            <svg style={{ height: 24, width: 24, marginRight: 6 }} fill="currentColor" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
                                <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.012 8.012 0 0 0 16 8c0-4.42-3.58-8-8-8z" />
                            </svg>
                        </span>
                        <Text color="gray" size="3">&copy; {new Date().getFullYear()} Go Initializer. All rights reserved.</Text>
                    </Flex>
                </footer>
            </div>
        </Theme>
    );
}

export default App;

