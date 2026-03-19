import { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import { generateProject, getMetaData } from '../service';
import { toGoVersionOptions, toSupportedFrameworkOptionsMap, toSupportedProjectTypes, toAddonOptions } from '../utils';
import type { AddonCategory, AddonState, FormErrors, FormTouched } from '../types';

const isMac = typeof navigator !== 'undefined' && navigator.platform.toLowerCase().includes('mac');

export function useGeneratorForm() {
  const [dockerSupport, setDockerSupport] = useState(false);
  const [selectedAddons, setSelectedAddons] = useState<AddonState>({
    cache: [],
    database: [],
    other: [],
  });
  const [projectType, setProjectType] = useState('');
  const [goVersion, setGoVersion] = useState('');
  const [framework, setFramework] = useState('');
  const [moduleName, setModuleName] = useState('');
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [touched, setTouched] = useState<FormTouched>({
    moduleName: false,
    name: false,
    description: false,
    goVersion: false,
    projectType: false,
    framework: false,
  });
  const [errors, setErrors] = useState<FormErrors>({});
  const [goVersionOptions, setGoVersionOptions] = useState<{ version: string; label: string }[]>([]);
  const [supportedProjectTypes, setSupportedProjectTypes] = useState<{ type: string; label: string }[]>([]);
  const [supportedFrameworkOptions, setSupportedFrameworkOptions] = useState<Record<string, { label: string; value: string }[]>>({});
  const [addonOptions, setAddonOptions] = useState<Record<string, { value: string; label: string }[]>>({});
  const [generateError, setGenerateError] = useState<string | null>(null);
  const [generateSuccess, setGenerateSuccess] = useState(false);
  const [successCountdown, setSuccessCountdown] = useState(5);
  const [isGenerating, setIsGenerating] = useState(false);

  const currentFrameworkOptions = useMemo(
    () => supportedFrameworkOptions[projectType] || [],
    [projectType, supportedFrameworkOptions],
  );

  useEffect(() => {
    setFramework('');
  }, [projectType, currentFrameworkOptions]);

  useEffect(() => {
    getMetaData()
      .then(data => {
        setGoVersionOptions(toGoVersionOptions(data.supportedGoVersions || []));
        setSupportedProjectTypes(toSupportedProjectTypes(data.supportedProjectTypes || []));
        setSupportedFrameworkOptions(toSupportedFrameworkOptionsMap(data.supportedFrameworks || []));
        setAddonOptions(toAddonOptions(data.supportedAddons || {}));
      })
      .catch(error => {
        console.error('Error fetching metadata:', error);
      });
  }, []);

  // Auto-dismiss success banner after 5 s with countdown
  useEffect(() => {
    if (!generateSuccess) {
      setSuccessCountdown(5);
      return;
    }
    setSuccessCountdown(5);
    const interval = setInterval(() => {
      setSuccessCountdown(prev => {
        if (prev <= 1) {
          clearInterval(interval);
          setGenerateSuccess(false);
          return 5;
        }
        return prev - 1;
      });
    }, 1000);
    return () => clearInterval(interval);
  }, [generateSuccess]);

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

  const validateInput = useCallback(() => {
    const newErrors: FormErrors = {};
    if (!moduleName.trim()) newErrors.moduleName = 'Module Name is required.';
    if (!name.trim()) newErrors.name = 'Name is required.';
    if (!description.trim()) newErrors.description = 'Description is required.';
    if (!projectType) newErrors.projectType = 'Project Type is required.';
    if (!goVersion) newErrors.goVersion = 'Go Version is required.';
    if (!framework) newErrors.framework = 'Framework/Dependency is required.';
    setErrors(newErrors);
    setTouched({ moduleName: true, name: true, description: true, goVersion: true, projectType: true, framework: true });
    return Object.keys(newErrors).length === 0;
  }, [moduleName, name, description, projectType, goVersion, framework]);

  const handleGenerate = useCallback(() => {
    if (!validateInput()) return;
    setGenerateError(null);
    setGenerateSuccess(false);
    setIsGenerating(true);
    generateProject({ projectType, goVersion, framework, moduleName, name, description, selectedAddons: selectedAddons as Record<string, string[]>, dockerSupport })
      .then(blob => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'project.zip';
        document.body.appendChild(a);
        a.click();
        a.remove();
        window.URL.revokeObjectURL(url);
        setGenerateSuccess(true);
      })
      .catch(error => {
        setGenerateError(error.message);
      })
      .finally(() => {
        setIsGenerating(false);
      });
  }, [validateInput, projectType, goVersion, framework, moduleName, name, description, selectedAddons, dockerSupport]);

  // Stable ref so the keydown listener never needs to be re-registered
  const handleGenerateRef = useRef(handleGenerate);
  useEffect(() => { handleGenerateRef.current = handleGenerate; }, [handleGenerate]);

  useEffect(() => {
    const listener = (e: KeyboardEvent) => {
      if ((isMac && e.metaKey && e.key === 'Enter') || (!isMac && e.ctrlKey && e.key === 'Enter')) {
        e.preventDefault();
        handleGenerateRef.current();
      }
    };
    window.addEventListener('keydown', listener);
    return () => window.removeEventListener('keydown', listener);
  }, []);

  return {
    dockerSupport, setDockerSupport,
    selectedAddons, handleAddonChange,
    projectType, setProjectType,
    goVersion, setGoVersion,
    framework, setFramework,
    moduleName, setModuleName,
    name, setName,
    description, setDescription,
    touched, setTouched,
    errors, setErrors,
    goVersionOptions,
    supportedProjectTypes,
    currentFrameworkOptions,
    addonOptions,
    generateError, setGenerateError,
    generateSuccess, setGenerateSuccess,
    successCountdown,
    isGenerating,
    isMac,
    handleGenerate,
  };
}
