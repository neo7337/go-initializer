import React, { useState, useEffect, useRef, useCallback } from 'react';
import './App.css';
import { generateProject, getMetaData } from './service';
import { toGoVersionOptions, toSupportedFrameworkOptionsMap, toSupportedProjectTypes } from './utils';
import Explore from './components/Explore';

function App() {
    const [theme, setTheme] = useState('dark');
    const [projectType, setProjectType] = useState('');
    const [goVersion, setGoVersion] = useState('');
    const [framework, setFramework] = useState('');
    const [moduleName, setModuleName] = useState('');
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [touched, setTouched] = useState<{
        moduleName: boolean;
        name: boolean;
        description: boolean;
        goVersion: boolean;
        projectType: boolean;
        framework: boolean;
    }>({
        moduleName: false,
        name: false,
        description: false,
        goVersion: false,
        projectType: false,
        framework: false,
    });
    const [errors, setErrors] = useState<{
        moduleName?: string;
        name?: string;
        description?: string;
        goVersion?: string;
        projectType?: string;
        framework?: string;
    }>({});
    const [goVersionOptions, setGoVersionOptions] = useState<{ version: string; label: string }[]>([]);
    const [supportedProjectTypes, setSupportedProjectTypes] = useState<{ type: string; label: string }[]>([]);
    const [supportedFrameworkOptions, setSupportedFrameworkOptions] = useState<Record<string, string[]>>({});
    const [showExplore, setShowExplore] = useState(false);

    // Framework options based on project type
    const currentFrameworkOptions = React.useMemo(() => {
        return supportedFrameworkOptions[projectType] || ['None'];
    }, [projectType, supportedFrameworkOptions]);

    useEffect(() => {
        // Reset framework if project type or options change, only if not already set or not in options
        setFramework('');
    }, [projectType, currentFrameworkOptions]);

    useEffect(() => {
        document.body.setAttribute('data-theme', theme);
    }, [theme]);

    const toggleTheme = () => {
        setTheme((prev) => (prev === 'light' ? 'dark' : 'light'));
    };

    const handleProjectTypeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setProjectType(e.target.value);
        setTouched(t => ({...t, projectType: true}));
        setErrors(errs => ({...errs, projectType: e.target.value.trim() ? undefined : 'Project Type is required.'}));
    };

    const handleGoVersionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setGoVersion(e.target.value);
        setTouched(t => ({...t, goVersion: true}));
        setErrors(errs => ({...errs, goVersion: e.target.value.trim() ? undefined : 'Go Version is required.'}));
    };

    const handleFrameworkChange = (fw: string) => {
        setFramework(fw);
        setTouched(t => ({...t, framework: true}));
        setErrors(errs => ({...errs, framework: fw.trim() ? undefined : 'Framework/Dependency is required.'}));
    };

    const validateInput = useCallback(() => {
        const newErrors: {moduleName?: string; name?: string; description?: string, projectType?: string, goVersion?: string, framework?: string} = {};
        if (!moduleName.trim()) newErrors.moduleName = 'Module Name is required.';
        if (!name.trim()) newErrors.name = 'Name is required.';
        if (!description.trim()) newErrors.description = 'Description is required.';
        if (!projectType) newErrors.projectType = 'Project Type is required.';
        if (!goVersion) newErrors.goVersion = 'Go Version is required.';
        if (!framework) newErrors.framework = 'Framework/Dependency is required.';
        setErrors(newErrors);
        setTouched({
            moduleName: true,
            name: true,
            description: true,
            goVersion: true,
            projectType: true,
            framework: true,
        });
        return Object.keys(newErrors).length === 0;
    }, [moduleName, name, description, projectType, goVersion, framework]);

    // Handler for Generate action
    const handleGenerate = useCallback(() => {
        if (!validateInput()) return;
        const requestBody = {
            projectType,
            goVersion,
            framework,
            moduleName,
            name,
            description,
        };
        generateProject(requestBody)
            .then(blob => {
                const filename = 'project.zip';
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = filename;
                document.body.appendChild(a);
                a.click();
                a.remove();
                window.URL.revokeObjectURL(url);
            })
            .catch(error => {
                alert('Error:\n' + error.message);
            });

    }, [validateInput, projectType, goVersion, framework, moduleName, name, description]);

    // Keep framework and other state always up to date for hotkey
    const frameworkRef = useRef(framework);
    const moduleNameRef = useRef(moduleName);
    const nameRef = useRef(name);
    const descriptionRef = useRef(description);
    useEffect(() => { frameworkRef.current = framework; }, [framework]);
    useEffect(() => { moduleNameRef.current = moduleName; }, [moduleName]);
    useEffect(() => { nameRef.current = name; }, [name]);
    useEffect(() => { descriptionRef.current = description; }, [description]);

    // Platform detection for hotkey icon
    const isMac = typeof navigator !== 'undefined' && navigator.platform.toLowerCase().includes('mac');

    // Hotkey: Cmd+Enter (macOS) or Ctrl+Enter (Windows/Linux)
    useEffect(() => {
        const listener = (e: KeyboardEvent) => {
            if (
                (isMac && e.metaKey && e.key === 'Enter') ||
                (!isMac && e.ctrlKey && e.key === 'Enter')
            ) {
                e.preventDefault();
                // Validation logic for hotkey
                handleGenerate();
            }
        };
        window.addEventListener('keydown', listener);
        return () => window.removeEventListener('keydown', listener);
    }, [handleGenerate, isMac]);

    useEffect(() => {
        getMetaData()
            .then(data => {
                // Handle metadata
                var formattedGoVersions = toGoVersionOptions(data.supportedGoVersions || []);
                console.log(formattedGoVersions)
                setGoVersionOptions(formattedGoVersions);
                var formattedSupportedProjectTypes = toSupportedProjectTypes(data.supportedProjectTypes || []);
                console.log(formattedSupportedProjectTypes)
                setSupportedProjectTypes(formattedSupportedProjectTypes);
                var formattedSupportedFrameworkOptions = toSupportedFrameworkOptionsMap(data.supportedFrameworks || []);
                console.log(formattedSupportedFrameworkOptions)
                setSupportedFrameworkOptions(formattedSupportedFrameworkOptions);
            })
            .catch(error => {
                console.error('Error fetching metadata:', error);
            });
    }, []);

    // const navigate = useNavigate();
    return (
        <div className="App" style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column', background: 'var(--background)', color: 'var(--text)', transition: 'background 0.3s, color 0.3s' }}>
            {/* Header */}
            <header style={{ background: 'var(--navbar-bg)', boxShadow: '0 2px 8px rgba(0,0,0,0.04)', padding: '1rem 2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', color: 'var(--navbar-text)' }}>
                <div style={{ display: 'flex', alignItems: 'center' }}>
                    {/* <img src={logo} alt="logo" style={{ height: 40, marginRight: 16 }} /> */}
                    <h1 style={{ fontSize: 28, fontWeight: 800, color: 'var(--navbar-text)', letterSpacing: 0.5 }}>go <span style={{ color: '#ffd700' }}>initializer</span></h1>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                    <button style={{ padding: 8, borderRadius: '50%', background: 'none', border: 'none', color: '#666', cursor: 'pointer', fontSize: 22 }} title="Toggle theme" onClick={toggleTheme}>
                        {theme === 'light' ? '🌙' : '☀️'}
                    </button>
                </div>
            </header>
            
            {/* Main Content */}
            <main style={{ flex: 1, padding: '2.5rem 1rem', background: 'var(--content-bg)', color: 'var(--text)', display: 'grid', gridTemplateColumns: '1fr', gap: 32, maxWidth: 1200, width: '100%', margin: '0 auto' }}>
                {showExplore ? (
                    <Explore onBack={() => setShowExplore(false)} />
                ) : (
                <div style={{ display: 'flex', flexDirection: 'column', gap: 32 }}>
                    {/* Go Version Card */}
                    <section style={{ background: 'var(--card-bg)', borderRadius: 16, boxShadow: '0 2px 12px 0 rgba(0,0,0,0.08)', padding: '2rem', marginBottom: 0, color: 'var(--text)' }}>
                        <h2 style={{ fontSize: 22, fontWeight: 700, marginBottom: 18, color: 'var(--text)' }}>Go Version</h2>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 16 }}>
                            {goVersionOptions.map((ver) => {
                                const value = ver.version;
                                const checked = goVersion === value;
                                return (
                                    <label
                                        key={ver.version}
                                        style={{
                                            display: 'flex',
                                            alignItems: 'center',
                                            cursor: 'pointer',
                                            border: checked ? '2px solid #ffd700' : '1.5px solid #e3e8f0',
                                            borderRadius: 999,
                                            padding: '0.5rem 1.5rem',
                                            fontWeight: 600,
                                            fontSize: 16,
                                            color: 'var(--text)',
                                            background: 'var(--card-bg)',
                                            transition: 'all 0.2s',
                                            marginRight: 0,
                                            boxShadow: checked ? '0 0 0 2px #ffd70033' : undefined,
                                        }}
                                    >
                                        <input
                                            type="radio"
                                            name="go-version"
                                            value={value}
                                            style={{ display: 'none' }}
                                            checked={checked}
                                            onChange={handleGoVersionChange}
                                        />
                                        {ver.label}
                                    </label>
                                );
                            })}
                        </div>
                        {errors.goVersion && touched.goVersion && (
                            <span style={{ color: '#ff4d4f', fontSize: 13, marginTop: 10, display: 'block' }}>{errors.goVersion}</span>
                        )}
                    </section>
                    {/* Project Type Card */}
                    <section style={{ background: 'var(--card-bg)', borderRadius: 16, boxShadow: '0 2px 12px 0 rgba(0,0,0,0.08)', padding: '2rem', marginBottom: 0, color: 'var(--text)' }}>
                        <h2 style={{ fontSize: 22, fontWeight: 700, marginBottom: 18, color: 'var(--text)' }}>Project Type</h2>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 16 }}>
                            {supportedProjectTypes.map((project_type) => {
                                const value = project_type.type.toLowerCase().replace(/ /g, '-');
                                const checked = projectType === value;
                                return (
                                    <label
                                        key={project_type.type}
                                        style={{
                                            display: 'flex',
                                            alignItems: 'center',
                                            cursor: 'pointer',
                                            border: checked ? '2px solid #ffd700' : '1.5px solid #e3e8f0',
                                            borderRadius: 999,
                                            padding: '0.5rem 1.5rem',
                                            fontWeight: 600,
                                            fontSize: 16,
                                            color: 'var(--text)',
                                            background: 'var(--card-bg)',
                                            transition: 'all 0.2s',
                                            marginRight: 0,
                                            boxShadow: checked ? '0 0 0 2px #ffd70033' : undefined,
                                        }}
                                    >
                                        <input
                                            type="radio"
                                            name="project-type"
                                            value={value}
                                            style={{ display: 'none' }}
                                            checked={checked}
                                            onChange={handleProjectTypeChange}
                                        />
                                        {project_type.label}
                                    </label>
                                );
                            })}
                        </div>
                        {errors.projectType && touched.projectType && (
                            <span style={{ color: '#ff4d4f', fontSize: 13, marginTop: 10, display: 'block' }}>{errors.projectType}</span>
                        )}
                    </section>

                    {/* Framework/Dependency Card */}
                    <section style={{ background: 'var(--card-bg)', borderRadius: 16, boxShadow: '0 2px 12px 0 rgba(0,0,0,0.08)', padding: '2rem', marginBottom: 0, color: 'var(--text)' }}>
                        <h2 style={{ fontSize: 22, fontWeight: 700, marginBottom: 18, color: 'var(--text)' }}>Select Framework/Dependency</h2>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 16 }}>
                            {currentFrameworkOptions.map((fw) => (
                                <label
                                    key={fw}
                                    style={{
                                        display: 'flex',
                                        alignItems: 'center',
                                        cursor: 'pointer',
                                        border: framework === fw ? '2px solid #ffd700' : '1.5px solid #e3e8f0',
                                        borderRadius: 999,
                                        padding: '0.5rem 1.5rem',
                                        fontWeight: 600,
                                        fontSize: 16,
                                        color: 'var(--text)',
                                        background: 'var(--card-bg)',
                                        transition: 'all 0.2s',
                                        marginRight: 0,
                                        boxShadow: framework === fw ? '0 0 0 2px #ffd70033' : undefined,
                                    }}
                                >
                                    <input
                                        type="radio"
                                        name="framework"
                                        value={fw}
                                        style={{ display: 'none' }}
                                        checked={framework === fw}
                                        onChange={() => handleFrameworkChange(fw)}
                                    />
                                    {fw}
                                </label>
                            ))}
                        </div>
                        {errors.framework && touched.framework && (
                            <span style={{ color: '#ff4d4f', fontSize: 13, marginTop: 10, display: 'block' }}>{errors.framework}</span>
                        )}
                    </section>
                    {/* Project Metadata Card */}
                    <section style={{ background: 'var(--card-bg)', borderRadius: 16, boxShadow: '0 2px 12px 0 rgba(0,0,0,0.08)', padding: '2rem', marginBottom: 0, color: 'var(--text)' }}>
                        <h2 style={{ fontSize: 22, fontWeight: 700, marginBottom: 18, color: 'var(--text)' }}>Project Metadata</h2>
                        <form style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 24 }}>
                            <div style={{ gridColumn: 'span 2' }}>
                                <label style={{ fontWeight: 600, color: 'var(--text)', marginBottom: 4, display: 'block' }}>Module Name</label>
                                <input
                                    type="text"
                                    placeholder="github.com/your/module"
                                    style={{
                                        width: '100%',
                                        padding: '0.7rem',
                                        fontSize: 16,
                                        borderRadius: 8,
                                        border: errors.moduleName && touched.moduleName ? '2px solid #ff4d4f' : '1.5px solid #e3e8f0',
                                        background: 'var(--card-bg)',
                                        color: 'var(--text)',
                                        fontWeight: 500,
                                        outline: errors.moduleName && touched.moduleName ? '2px solid #ff4d4f' : 'none',
                                        marginTop: 4,
                                    }}
                                    value={moduleName}
                                    onChange={e => {
                                        setModuleName(e.target.value);
                                        setTouched(t => ({...t, moduleName: true}));
                                        setErrors(errs => ({...errs, moduleName: e.target.value.trim() ? undefined : 'Module Name is required.'}));
                                    }}
                                    required
                                    onBlur={() => setTouched(t => ({...t, moduleName: true}))}
                                />
                                {errors.moduleName && touched.moduleName && (
                                    <span style={{ color: '#ff4d4f', fontSize: 13, marginTop: 2, display: 'block' }}>{errors.moduleName}</span>
                                )}
                            </div>
                            <div>
                                <label style={{ fontWeight: 600, color: 'var(--text)', marginBottom: 4, display: 'block' }}>Name</label>
                                <input
                                    type="text"
                                    placeholder="my-app"
                                    style={{
                                        width: '100%',
                                        padding: '0.7rem',
                                        fontSize: 16,
                                        borderRadius: 8,
                                        border: errors.name && touched.name ? '2px solid #ff4d4f' : '1.5px solid #e3e8f0',
                                        background: 'var(--card-bg)',
                                        color: 'var(--text)',
                                        fontWeight: 500,
                                        outline: errors.name && touched.name ? '2px solid #ff4d4f' : 'none',
                                        marginTop: 4,
                                    }}
                                    value={name}
                                    onChange={e => {
                                        setName(e.target.value);
                                        setTouched(t => ({...t, name: true}));
                                        setErrors(errs => ({...errs, name: e.target.value.trim() ? undefined : 'Name is required.'}));
                                    }}
                                    required
                                    onBlur={() => setTouched(t => ({...t, name: true}))}
                                />
                                {errors.name && touched.name && (
                                    <span style={{ color: '#ff4d4f', fontSize: 13, marginTop: 2, display: 'block' }}>{errors.name}</span>
                                )}
                            </div>
                            <div>
                                <label style={{ fontWeight: 600, color: 'var(--text)', marginBottom: 4, display: 'block' }}>Description</label>
                                <input
                                    type="text"
                                    placeholder="Short project description"
                                    style={{
                                        width: '100%',
                                        padding: '0.7rem',
                                        fontSize: 16,
                                        borderRadius: 8,
                                        border: errors.description && touched.description ? '2px solid #ff4d4f' : '1.5px solid #e3e8f0',
                                        background: 'var(--card-bg)',
                                        color: 'var(--text)',
                                        fontWeight: 500,
                                        outline: errors.description && touched.description ? '2px solid #ff4d4f' : 'none',
                                        marginTop: 4,
                                    }}
                                    value={description}
                                    onChange={e => {
                                        setDescription(e.target.value);
                                        setTouched(t => ({...t, description: true}));
                                        setErrors(errs => ({...errs, description: e.target.value.trim() ? undefined : 'Description is required.'}));
                                    }}
                                    required
                                    onBlur={() => setTouched(t => ({...t, description: true}))}
                                />
                                {errors.description && touched.description && (
                                    <span style={{ color: '#ff4d4f', fontSize: 13, marginTop: 2, display: 'block' }}>{errors.description}</span>
                                )}
                            </div>
                        </form>
                    </section>
                    {/* Generate Button */}
                    <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 16 }}>
                        <button
                            style={{ background: '#ffd700', color: theme === 'dark' ? '#23272f' : '#222', fontWeight: 700, fontSize: 18, padding: '0.9rem 2.5rem', borderRadius: 10, border: 'none', boxShadow: '0 2px 8px 0 rgba(34,34,34,0.10)', letterSpacing: 0.5, cursor: 'pointer', transition: 'background 0.2s, color 0.2s', display: 'flex', alignItems: 'center', gap: 10 }}
                            onClick={handleGenerate}
                            title={isMac ? 'Cmd+Enter (macOS)' : 'Ctrl+Enter (Windows/Linux)'}
                        >
                            <span style={{ display: 'flex', alignItems: 'center', gap: 4, marginRight: 10, fontSize: 15 }}>
                                {isMac ? (
                                    <>
                                        {/* macOS icon */}
                                        <span style={{ fontWeight: 600, fontFamily: 'monospace', fontSize: 22 }}>⌘</span>
                                        <span style={{ fontWeight: 600, fontFamily: 'monospace', fontSize: 22 }}>↵</span>
                                    </>
                                ) : (
                                    <>
                                        {/* Windows/Linux icon */}
                                        <svg style={{ height: 24, width: 24 }} viewBox="0 0 20 20" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                                            <rect x="2" y="4" width="16" height="12" rx="2" />
                                        </svg>
                                        <span style={{ fontWeight: 600, fontFamily: 'monospace', fontSize: 20 }}>Ctrl</span>
                                        <span style={{ fontWeight: 600, fontFamily: 'monospace', fontSize: 22 }}>↵</span>
                                    </>
                                )}
                            </span>
                            GENERATE
                        </button>
                        <button
                            style={{ background: theme === 'dark' ? '#23272f' : '#f8fafc', color: 'var(--text)', fontWeight: 700, fontSize: 18, padding: '0.9rem 2.5rem', borderRadius: 10, border: '1.5px solid #e3e8f0', boxShadow: '0 2px 8px 0 rgba(34,34,34,0.06)', letterSpacing: 0.5, cursor: 'pointer', transition: 'background 0.2s, color 0.2s' }}
                            onClick={() => setShowExplore(true)}
                        >
                            EXPLORE
                        </button>
                    </div>
                </div>
                )}
            </main>
            {/* Footer */}
            <footer style={{ background: 'var(--footer-bg)', color: 'var(--footer-text)', boxShadow: '0 -2px 8px rgba(0,0,0,0.04)', padding: '1rem 2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 'auto' }}>
                <span style={{ display: 'flex', alignItems: 'center' }}>
                    <svg style={{ height: 24, width: 24, marginRight: 6 }} fill="currentColor" viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
                        <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.012 8.012 0 0 0 16 8c0-4.42-3.58-8-8-8z" />
                    </svg>
                </span>
                <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
                    <span style={{ color: 'var(--footer-text)', fontWeight: 500, fontSize: 15 }}>&copy; {new Date().getFullYear()} Go Initializer. All rights reserved.</span>
                </div>
            </footer>
        </div>
    );
}

export default App;
