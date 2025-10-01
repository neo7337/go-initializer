import React, { useState, useEffect, useRef, useCallback } from 'react';
import './App.css';
import { generateProject, getMetaData } from './service';
import { toGoVersionOptions, toSupportedFrameworkOptionsMap, toSupportedProjectTypes } from './utils';
import Explore from './components/Explore';
import { Theme, Card, Button, Text, Flex, Heading, RadioGroup, Checkbox } from '@radix-ui/themes';


function App() {
    // Docker support state
    const [dockerSupport, setDockerSupport] = useState(false);
    // Addons state and types
    type AddonCategory = 'cache' | 'database' | 'other';
    type AddonState = {
        cache: string[];
        database: string[];
        other: string[];
    };
    const [selectedAddons, setSelectedAddons] = useState<AddonState>({
        cache: [],
        database: [],
        other: [],
    });
    const addonOptions: Record<AddonCategory, { value: string; label: string; description: string }[]> = {
        cache: [
            { value: 'redis', label: 'Redis', description: 'In-memory cache and message broker' },
            { value: 'memcached', label: 'Memcached', description: 'High-performance distributed memory caching system' },
        ],
        database: [
            { value: 'gorm', label: 'Gorm', description: 'Popular ORM library for Golang' },
            { value: 'ent', label: 'Ent', description: 'Entity framework for Go' },
        ],
        other: [
            { value: 'zap', label: 'Zap', description: 'Blazing fast, structured, leveled logging in Go' },
            { value: 'cobra', label: 'Cobra', description: 'Commander for modern Go CLI interactions' },
        ],
    };
    const handleAddonChange = (category: AddonCategory, value: string) => {
        setSelectedAddons(prev => {
            const alreadySelected = prev[category].includes(value);
            return {
                ...prev,
                [category]: alreadySelected
                    ? prev[category].filter((v: string) => v !== value)
                    : [...prev[category], value],
            };
        });
    };
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
        <Theme appearance={theme === 'dark' ? 'dark' : 'light'}>
            <div className="App" style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column', background: 'var(--background)', color: 'var(--text)', transition: 'background 0.3s, color 0.3s' }}>
            {/* Header */}
            <header style={{ background: 'var(--navbar-bg)', boxShadow: '0 2px 8px rgba(0,0,0,0.04)', padding: '1rem 2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center', color: 'var(--navbar-text)' }}>
                <Flex align="center" justify="between" style={{ width: '100%' }}>
                    <Flex align="center">
                        <Heading size="6" weight="bold" style={{ color: 'var(--navbar-text)', letterSpacing: 0.5 }}>
                            go<span style={{ color: '#ffd700' }}>initializer</span>
                        </Heading>
                    </Flex>
                    <Button variant="ghost" color="gray" radius="full" size="3" onClick={toggleTheme} title="Toggle theme">
                        {theme === 'light' ? '🌙' : '☀️'}
                    </Button>
                </Flex>
            </header>
            
            {/* Main Content */}
            <main style={{ flex: 1, padding: '2.5rem 1rem', background: 'var(--content-bg)', color: 'var(--text)', display: 'grid', gridTemplateColumns: '1fr', gap: 32, maxWidth: 1200, width: '100%', margin: '0 auto' }}>
                {showExplore ? (
                    <Explore onBack={() => setShowExplore(false)} />
                ) : (
                <div style={{ display: 'flex', flexDirection: 'column', gap: 32 }}>
                    {/* Go Version Card */}
                    <Card style={{ marginBottom: 0 }}>
                        <Heading size="4" mb="3">Go Version</Heading>
                        <RadioGroup.Root
                            value={goVersion}
                            onValueChange={val => {
                                setGoVersion(val);
                                setTouched(t => ({...t, goVersion: true}));
                                setErrors(errs => ({...errs, goVersion: val.trim() ? undefined : 'Go Version is required.'}));
                            }}
                            orientation="horizontal"
                            style={{ display: 'flex', flexDirection: 'row', flexWrap: 'wrap', gap: 16 }}
                        >
                            {goVersionOptions.map((ver) => (
                                <RadioGroup.Item key={ver.version} value={ver.version} style={{ marginRight: 0 }}>
                                    <Text as="span" size="3" weight={goVersion === ver.version ? 'bold' : 'regular'}>{ver.label}</Text>
                                </RadioGroup.Item>
                            ))}
                        </RadioGroup.Root>
                        {errors.goVersion && touched.goVersion && (
                            <Text color="red" size="2" mt="2" as="span">{errors.goVersion}</Text>
                        )}
                    </Card>
                    {/* Project Type Card */}
                    <Card style={{ marginBottom: 0 }}>
                        <Heading size="4" mb="3">Project Type</Heading>
                        <RadioGroup.Root
                            value={projectType}
                            onValueChange={val => {
                                setProjectType(val);
                                setTouched(t => ({...t, projectType: true}));
                                setErrors(errs => ({...errs, projectType: val.trim() ? undefined : 'Project Type is required.'}));
                            }}
                            orientation="horizontal"
                            style={{ display: 'flex', flexDirection: 'row', flexWrap: 'wrap', gap: 16 }}
                        >
                            {supportedProjectTypes.map((project_type) => {
                                const value = project_type.type.toLowerCase().replace(/ /g, '-');
                                return (
                                    <RadioGroup.Item key={project_type.type} value={value} style={{ marginRight: 0 }}>
                                        <Text as="span" size="3" weight={projectType === value ? 'bold' : 'regular'}>{project_type.label}</Text>
                                    </RadioGroup.Item>
                                );
                            })}
                        </RadioGroup.Root>
                        {errors.projectType && touched.projectType && (
                            <Text color="red" size="2" mt="2" as="span">{errors.projectType}</Text>
                        )}
                    </Card>

                    {/* Framework/Dependency Card */}
                    <Card style={{ marginBottom: 0 }}>
                        <Heading size="4" mb="3">Select Framework/Dependency</Heading>
                        <RadioGroup.Root
                            value={framework}
                            onValueChange={val => {
                                setFramework(val);
                                setTouched(t => ({...t, framework: true}));
                                setErrors(errs => ({...errs, framework: val.trim() ? undefined : 'Framework/Dependency is required.'}));
                            }}
                            orientation="horizontal"
                            style={{ display: 'flex', flexDirection: 'row', flexWrap: 'wrap', gap: 16 }}
                        >
                            {currentFrameworkOptions.map((fw) => (
                                <RadioGroup.Item key={fw} value={fw} style={{ marginRight: 0 }}>
                                    <Text as="span" size="3" weight={framework === fw ? 'bold' : 'regular'}>{fw}</Text>
                                </RadioGroup.Item>
                            ))}
                        </RadioGroup.Root>
                        {errors.framework && touched.framework && (
                            <Text color="red" size="2" mt="2" as="span">{errors.framework}</Text>
                        )}
                    </Card>
                    {/* Addons Card */}
                    <Card style={{ marginBottom: 0 }}>
                        <Heading size="4" mb="3">Addons</Heading>
                        <Flex direction="row" gap="6" wrap="wrap">
                            {/* Cache Addons */}
                            <div style={{ minWidth: 180 }}>
                                <Text as="div" size="3" weight="bold" mb="2">Cache</Text>
                                <Flex gap="4" direction="column">
                                    {addonOptions.cache.map(opt => (
                                        <label key={opt.value} style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer' }}>
                                            <Checkbox
                                                checked={selectedAddons.cache.includes(opt.value)}
                                                onCheckedChange={() => handleAddonChange('cache', opt.value)}
                                            />
                                            <Text as="span" size="2">{opt.label}</Text>
                                        </label>
                                    ))}
                                </Flex>
                            </div>
                            {/* Database Addons */}
                            <div style={{ minWidth: 180 }}>
                                <Text as="div" size="3" weight="bold" mb="2">Database</Text>
                                <Flex gap="4" direction="column">
                                    {addonOptions.database.map(opt => (
                                        <label key={opt.value} style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer' }}>
                                            <Checkbox
                                                checked={selectedAddons.database.includes(opt.value)}
                                                onCheckedChange={() => handleAddonChange('database', opt.value)}
                                            />
                                            <Text as="span" size="2">{opt.label}</Text>
                                        </label>
                                    ))}
                                </Flex>
                            </div>
                            {/* Other Libraries */}
                            <div style={{ minWidth: 180 }}>
                                <Text as="div" size="3" weight="bold" mb="2">Other Libraries</Text>
                                <Flex gap="4" direction="column">
                                    {addonOptions.other.map(opt => (
                                        <label key={opt.value} style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer' }}>
                                            <Checkbox
                                                checked={selectedAddons.other.includes(opt.value)}
                                                onCheckedChange={() => handleAddonChange('other', opt.value)}
                                            />
                                            <Text as="span" size="2">{opt.label}</Text>
                                        </label>
                                    ))}
                                </Flex>
                            </div>
                        </Flex>
                    </Card>
                    {/* Docker Support */}
                    <Card style={{ marginBottom: 0 }}>
                        <Flex align="center" gap="3">
                            <Checkbox
                                checked={dockerSupport}
                                onCheckedChange={() => setDockerSupport(v => !v)}
                                id="docker-support"
                            />
                            <label htmlFor="docker-support" style={{ cursor: 'pointer' }}>
                                <Text as="span" size="3" weight="bold">Add Docker support</Text>
                            </label>
                        </Flex>
                    </Card>
                    {/* Project Metadata Card */}
                    <Card style={{ marginBottom: 0 }}>
                        <Heading size="4" mb="3">Project Metadata</Heading>
                        <Flex direction="column">
                            <Text as="label" size="2" weight="bold" mb="1">Module Name</Text>
                            <input
                                type="text"
                                placeholder="github.com/your/module"
                                value={moduleName}
                                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                                    setModuleName(e.target.value);
                                    setTouched(t => ({...t, moduleName: true}));
                                    setErrors(errs => ({...errs, moduleName: e.target.value.trim() ? undefined : 'Module Name is required.'}));
                                }}
                                onBlur={() => setTouched(t => ({...t, moduleName: true}))}
                                required
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
                            />
                            {errors.moduleName && touched.moduleName && (
                                <Text color="red" size="2" mt="1" as="span">{errors.moduleName}</Text>
                            )}
                            <Flex gap="4" mt="4">
                                <div style={{ flex: 1 }}>
                                    <Text as="label" size="2" weight="bold" mb="1">Name</Text>
                                    <input
                                        type="text"
                                        placeholder="my-app"
                                        value={name}
                                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                                            setName(e.target.value);
                                            setTouched(t => ({...t, name: true}));
                                            setErrors(errs => ({...errs, name: e.target.value.trim() ? undefined : 'Name is required.'}));
                                        }}
                                        onBlur={() => setTouched(t => ({...t, name: true}))}
                                        required
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
                                    />
                                    {errors.name && touched.name && (
                                        <Text color="red" size="2" mt="1" as="span">{errors.name}</Text>
                                    )}
                                </div>
                                <div style={{ flex: 1 }}>
                                    <Text as="label" size="2" weight="bold" mb="1">Description</Text>
                                    <input
                                        type="text"
                                        placeholder="Short project description"
                                        value={description}
                                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                                            setDescription(e.target.value);
                                            setTouched(t => ({...t, description: true}));
                                            setErrors(errs => ({...errs, description: e.target.value.trim() ? undefined : 'Description is required.'}));
                                        }}
                                        onBlur={() => setTouched(t => ({...t, description: true}))}
                                        required
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
                                    />
                                    {errors.description && touched.description && (
                                        <Text color="red" size="2" mt="1" as="span">{errors.description}</Text>
                                    )}
                                </div>
                            </Flex>
                        </Flex>
                    </Card>
                    {/* Generate Button */}
                    <Flex justify="end" gap="3" mt="4">
                        <Button
                            size="4"
                            color="yellow"
                            variant="solid"
                            onClick={handleGenerate}
                            title={isMac ? 'Cmd+Enter (macOS)' : 'Ctrl+Enter (Windows/Linux)'}
                        >
                            <span style={{ display: 'flex', alignItems: 'center', gap: 4, marginRight: 10, fontSize: 15 }}>
                                {isMac ? (
                                    <>
                                        <span style={{ fontWeight: 600, fontFamily: 'monospace', fontSize: 22 }}>⌘</span>
                                        <span style={{ fontWeight: 600, fontFamily: 'monospace', fontSize: 22 }}>↵</span>
                                    </>
                                ) : (
                                    <>
                                        <svg style={{ height: 24, width: 24 }} viewBox="0 0 20 20" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                                            <rect x="2" y="4" width="16" height="12" rx="2" />
                                        </svg>
                                        <span style={{ fontWeight: 600, fontFamily: 'monospace', fontSize: 20 }}>Ctrl</span>
                                        <span style={{ fontWeight: 600, fontFamily: 'monospace', fontSize: 22 }}>↵</span>
                                    </>
                                )}
                            </span>
                            GENERATE
                        </Button>
                        <Button
                            size="4"
                            color={theme === 'dark' ? 'gray' : 'gray'}
                            variant="soft"
                            onClick={() => setShowExplore(true)}
                        >
                            EXPLORE
                        </Button>
                    </Flex>
                </div>
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
