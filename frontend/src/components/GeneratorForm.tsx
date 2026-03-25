import React, { useRef, useCallback } from 'react';
import { TextField, Text, Flex } from '@radix-ui/themes';
import { useGeneratorForm } from '../hooks/useGeneratorForm';
import type { AddonCategory } from '../types';

/* ── Pill selector (Go Version / Project Type / Framework) ─────── */
type PillButtonProps = {
  selected: boolean;
  onClick: () => void;
  children: React.ReactNode;
};

const PillButton: React.FC<PillButtonProps> = ({ selected, onClick, children }) => (
  <button
    type="button"
    onClick={onClick}
    aria-pressed={selected}
    style={{
      display: 'inline-flex',
      alignItems: 'center',
      padding: '0.45rem 1.1rem',
      borderRadius: '9999px',
      border: selected
        ? '1.5px solid var(--color-accent)'
        : '1.5px solid var(--color-border)',
      background: selected ? 'var(--color-accent-dim)' : 'var(--color-surface-3)',
      color: selected ? 'var(--color-accent)' : 'var(--color-text-secondary)',
      fontWeight: selected ? 600 : 400,
      fontSize: 14,
      cursor: 'pointer',
      boxShadow: selected ? '0 0 0 2px rgba(255,215,0,0.18)' : 'none',
      transition: 'border-color 0.12s, background 0.12s, color 0.12s, box-shadow 0.12s',
      fontFamily: 'inherit',
      letterSpacing: '0.01em',
    }}
  >
    {children}
  </button>
);

/* ── Tag chip (addons / docker) ─────────────────────────────────── */
type TagChipProps = {
  selected: boolean;
  onClick: () => void;
  children: React.ReactNode;
};

const TagChip: React.FC<TagChipProps> = ({ selected, onClick, children }) => (
  <button
    type="button"
    onClick={onClick}
    aria-pressed={selected}
    style={{
      display: 'inline-flex',
      alignItems: 'center',
      gap: 5,
      padding: '0.3rem 0.85rem',
      borderRadius: 'var(--radius-md)',
      border: selected
        ? '1.5px solid var(--color-accent)'
        : '1.5px solid var(--color-border)',
      background: selected ? 'var(--color-accent-dim)' : 'var(--color-surface-3)',
      color: selected ? 'var(--color-accent)' : 'var(--color-text-secondary)',
      fontWeight: selected ? 600 : 400,
      fontSize: 13,
      cursor: 'pointer',
      boxShadow: selected ? '0 0 0 2px rgba(255,215,0,0.15)' : 'none',
      transition: 'border-color 0.12s, background 0.12s, color 0.12s, box-shadow 0.12s',
      fontFamily: 'inherit',
    }}
  >
    {selected && <span style={{ fontSize: 10, lineHeight: 1, marginRight: 2 }}>✓</span>}
    {children}
  </button>
);

/* ── Numbered section label ─────────────────────────────────────── */
const SectionLabel: React.FC<{ number: string; label: string; icon?: React.ReactNode }> = ({ number, label, icon }) => (
  <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 14 }}>
    <span style={{
      fontFamily: 'monospace',
      fontSize: 11,
      fontWeight: 700,
      color: 'var(--color-accent)',
      letterSpacing: '0.06em',
      opacity: 0.75,
    }}>
      {number}
    </span>
    {icon && (
      <span style={{ display: 'flex', alignItems: 'center', color: 'var(--color-text-muted)', flexShrink: 0 }}>
        {icon}
      </span>
    )}
    <span style={{
      fontSize: 12,
      fontWeight: 700,
      color: 'var(--color-text-primary)',
      textTransform: 'uppercase',
      letterSpacing: '0.08em',
    }}>
      {label}
    </span>
  </div>
);

/* ── Section icons (inline SVG, 14×14) ─────────────────────────── */
const TagIcon = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
    <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z" />
    <line x1="7" y1="7" x2="7.01" y2="7" />
  </svg>
);
const LayersIcon2 = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
    <polygon points="12 2 2 7 12 12 22 7 12 2" />
    <polyline points="2 17 12 22 22 17" />
    <polyline points="2 12 12 17 22 12" />
  </svg>
);
const PlugIcon2 = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
    <path d="M18 6L6 18" /><path d="M7 17l-5 5" /><path d="M17 7l5-5" />
    <path d="M8 6v5a4 4 0 0 0 8 0V6" />
  </svg>
);
const PuzzleIcon2 = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
    <path d="M20.24 12.24a6 6 0 0 0-8.49-8.49L5 10.5V19h8.5l6.74-6.76z" />
    <line x1="16" y1="8" x2="2" y2="22" />
    <line x1="17" y1="15" x2="9" y2="15" />
  </svg>
);
const BoxIcon = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
    <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z" />
    <polyline points="3.27 6.96 12 12.01 20.73 6.96" />
    <line x1="12" y1="22.08" x2="12" y2="12" />
  </svg>
);
const FileIcon = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
    <polyline points="14 2 14 8 20 8" />
    <line x1="16" y1="13" x2="8" y2="13" />
    <line x1="16" y1="17" x2="8" y2="17" />
    <polyline points="10 9 9 9 8 9" />
  </svg>
);

/* ── Terminal error banner ─────────────────────────────────────── */
const ErrorBanner: React.FC<{ message: string; onDismiss: () => void }> = ({ message, onDismiss }) => (
  <div
    role="alert"
    style={{
      display: 'flex',
      alignItems: 'flex-start',
      gap: 10,
      padding: '0.75rem 1rem',
      borderRadius: 'var(--radius-md)',
      background: 'var(--color-destructive-dim)',
      border: '1px solid var(--color-destructive)',
      borderLeft: '3px solid var(--color-destructive)',
      color: 'var(--color-text-primary)',
    }}
  >
    {/* Warning triangle icon */}
    <svg style={{ width: 16, height: 16, flexShrink: 0, marginTop: 2, color: 'var(--color-destructive)' }} viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
      <path fillRule="evenodd" d="M8.485 2.495c.673-1.167 2.357-1.167 3.03 0l6.28 10.875c.673 1.167-.17 2.625-1.516 2.625H3.72c-1.347 0-2.189-1.458-1.515-2.625L8.485 2.495zM10 5a.75.75 0 01.75.75v3.5a.75.75 0 01-1.5 0v-3.5A.75.75 0 0110 5zm0 9a1 1 0 100-2 1 1 0 000 2z" clipRule="evenodd" />
    </svg>
    <div style={{ flex: 1, fontSize: 13 }}>
      <span style={{ fontFamily: 'monospace', fontWeight: 700, color: 'var(--color-destructive)', marginRight: 6 }}>ERROR</span>
      {message}
    </div>
    <button
      type="button"
      onClick={onDismiss}
      aria-label="Dismiss error"
      style={{
        background: 'none',
        border: 'none',
        color: 'var(--color-text-muted)',
        cursor: 'pointer',
        padding: '0 2px',
        fontSize: 14,
        lineHeight: 1,
        flexShrink: 0,
      }}
    >
      ✕
    </button>
  </div>
);

/* ── Animated checkmark SVG ─────────────────────────────────────── */
const AnimatedCheck: React.FC = () => (
  <svg
    style={{ width: 18, height: 18, flexShrink: 0 }}
    viewBox="0 0 24 24"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    aria-hidden="true"
  >
    <circle cx="12" cy="12" r="10" stroke="var(--color-success)" strokeWidth="1.5" opacity="0.3" />
    <path
      className="checkmark-path"
      d="M7 12.5 L10.5 16 L17 9"
      stroke="var(--color-success)"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      fill="none"
    />
  </svg>
);

/* ── Success banner with countdown progress bar ─────────────────── */
const SuccessBanner: React.FC<{ countdown: number; onDismiss: () => void }> = ({ countdown, onDismiss }) => (
  <div
    role="status"
    style={{
      borderRadius: 'var(--radius-md)',
      background: 'var(--color-success-dim)',
      border: '1px solid var(--color-success)',
      borderLeft: '3px solid var(--color-success)',
      overflow: 'hidden',
    }}
  >
    <div style={{ display: 'flex', alignItems: 'center', gap: 10, padding: '0.75rem 1rem' }}>
      <AnimatedCheck />
      <div style={{ flex: 1, fontSize: 13 }}>
        <span style={{ fontFamily: 'monospace', fontWeight: 700, color: 'var(--color-success)', marginRight: 6 }}>SUCCESS</span>
        Your project zip has been downloaded.
      </div>
      <span style={{ fontSize: 11, color: 'var(--color-text-muted)', fontFamily: 'monospace', marginRight: 6 }}>
        {countdown}s
      </span>
      <button
        type="button"
        onClick={onDismiss}
        aria-label="Dismiss success"
        style={{
          background: 'none',
          border: 'none',
          color: 'var(--color-text-muted)',
          cursor: 'pointer',
          padding: '0 2px',
          fontSize: 14,
          lineHeight: 1,
          flexShrink: 0,
        }}
      >
        ✕
      </button>
    </div>
    {/* Countdown progress bar */}
    <div style={{ height: 2, background: 'rgba(34,197,94,0.15)' }}>
      <div
        className="success-progress-bar"
        style={{ height: '100%', background: 'var(--color-success)', transformOrigin: 'left' }}
      />
    </div>
  </div>
);

/* ── Field error message ─────────────────────────────────────────── */
const FieldError: React.FC<{ message: string }> = ({ message }) => (
  <div style={{ display: 'flex', alignItems: 'center', gap: 4, marginTop: 4 }}>
    <svg style={{ width: 12, height: 12, flexShrink: 0, color: 'var(--color-destructive)' }} viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
      <path d="M8 1a7 7 0 100 14A7 7 0 008 1zM7 4a1 1 0 112 0v4a1 1 0 11-2 0V4zm1 8a1 1 0 100-2 1 1 0 000 2z" />
    </svg>
    <span style={{ fontSize: 12, color: 'var(--color-destructive)', lineHeight: 1.3 }}>{message}</span>
  </div>
);

/* ── SVG spinner ────────────────────────────────────────────────── */
const Spinner: React.FC = () => (
  <svg
    className="form-spinner"
    style={{ width: 18, height: 18, flexShrink: 0 }}
    viewBox="0 0 24 24"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    aria-hidden="true"
  >
    <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="3" opacity="0.25" />
    <path d="M12 2 A10 10 0 0 1 22 12" stroke="currentColor" strokeWidth="3" strokeLinecap="round" />
  </svg>
);


/* ── Shared field label style ───────────────────────────────────── */
const fieldLabel: React.CSSProperties = {
  color: 'var(--color-text-secondary)',
  display: 'block',
  marginBottom: 6,
};

/* ═══════════════════════════════════════════════════════════════════
   GeneratorForm
═══════════════════════════════════════════════════════════════════ */
const GeneratorForm: React.FC<{ onInteract?: () => void }> = ({ onInteract }) => {
  const {
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
  } = useGeneratorForm();

  const hasInteracted = useRef(false);
  const notifyInteract = useCallback(() => {
    if (!hasInteracted.current) {
      hasInteracted.current = true;
      onInteract?.();
    }
  }, [onInteract]);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }} role="form" aria-label="Project generator">

      {/* 01 — Go Version */}
      <div className="form-card">
        <SectionLabel number="01" label="Go Version" icon={<TagIcon />} />
        <p className="section-desc">Select the Go runtime version for your project</p>
        <div role="group" aria-label="Go Version">
          <Flex wrap="wrap" gap="2">
            {goVersionOptions.map(ver => (
              <PillButton
                key={ver.version}
                selected={goVersion === ver.version}
                onClick={() => {
                  notifyInteract();
                  setGoVersion(ver.version);
                  setTouched(t => ({ ...t, goVersion: true }));
                  setErrors(errs => ({ ...errs, goVersion: undefined }));
                }}
              >
                {ver.label}
              </PillButton>
            ))}
          </Flex>
        </div>
        {errors.goVersion && touched.goVersion && <FieldError message={errors.goVersion} />}
      </div>

      {/* 02 — Project Type */}
      <div className="form-card">
        <SectionLabel number="02" label="Project Type" icon={<LayersIcon2 />} />
        <p className="section-desc">What kind of Go application do you want to scaffold?</p>
        <div role="group" aria-label="Project Type">
          <Flex wrap="wrap" gap="2">
            {supportedProjectTypes.map(pt => {
              const value = pt.type.toLowerCase().replace(/ /g, '-');
              return (
                <PillButton
                  key={pt.type}
                  selected={projectType === value}
                  onClick={() => {
                    notifyInteract();
                    setProjectType(value);
                    setTouched(t => ({ ...t, projectType: true }));
                    setErrors(errs => ({ ...errs, projectType: undefined }));
                  }}
                >
                  {pt.label}
                </PillButton>
              );
            })}
          </Flex>
        </div>
        {errors.projectType && touched.projectType && <FieldError message={errors.projectType} />}
      </div>

      {/* 03 — Framework */}
      <div className="form-card">
        <SectionLabel number="03" label="Framework / Dependency" icon={<PlugIcon2 />} />
        <p className="section-desc">Choose the HTTP or CLI framework to use</p>
        {!currentFrameworkOptions.length ? (
          <Text size="2" style={{ color: 'var(--color-text-muted)' }}>
            Select a project type to see available frameworks.
          </Text>
        ) : (
          <div key={projectType} role="group" aria-label="Framework" className="framework-options-fade">
            <Flex wrap="wrap" gap="2">
              {currentFrameworkOptions.map(fw => (
                <PillButton
                  key={fw.value}
                  selected={framework === fw.value}
                  onClick={() => {
                    notifyInteract();
                    setFramework(fw.value);
                    setTouched(t => ({ ...t, framework: true }));
                    setErrors(errs => ({ ...errs, framework: undefined }));
                  }}
                >
                  {fw.label}
                </PillButton>
              ))}
            </Flex>
          </div>
        )}
        {errors.framework && touched.framework && <FieldError message={errors.framework} />}
      </div>

      {/* 04 — Addons */}
      <div className="form-card">
        <SectionLabel number="04" label="Addons" icon={<PuzzleIcon2 />} />
        <p className="section-desc">Optional libraries to include (cache, database, logging)</p>
        <Flex direction="column" gap="4">
          {Object.entries(addonOptions).map(([category, opts]) => {
            const categoryLabel = category === 'other' ? 'Other Libraries' : category.charAt(0).toUpperCase() + category.slice(1);
            return (
            <div key={category} role="group" aria-label={categoryLabel}>
              <Text
                as="p"
                size="1"
                weight="bold"
                style={{
                  color: 'var(--color-text-muted)',
                  textTransform: 'uppercase',
                  letterSpacing: '0.07em',
                  marginBottom: 8,
                }}
              >
                {categoryLabel}
              </Text>
              <Flex wrap="wrap" gap="2">
                {opts.map(opt => (
                  <TagChip
                    key={opt.value}
                    selected={selectedAddons[category]?.includes(opt.value) ?? false}
                    onClick={() => { notifyInteract(); handleAddonChange(category as AddonCategory, opt.value); }}
                  >
                    {opt.label}
                  </TagChip>
                ))}
              </Flex>
            </div>
            );
          })}
        </Flex>
      </div>

      {/* 05 — Docker Support */}
      <div className="form-card">
        <SectionLabel number="05" label="Docker Support" icon={<BoxIcon />} />
        <p className="section-desc">Generate a multi-stage production Dockerfile</p>
        <Flex align="center" gap="3">
          <TagChip
            selected={dockerSupport}
            onClick={() => { notifyInteract(); setDockerSupport(v => !v); }}
          >
            Generate Dockerfile
          </TagChip>
          <Text size="2" style={{ color: 'var(--color-text-muted)' }}>
            Multi-stage production Dockerfile
          </Text>
        </Flex>
      </div>

      {/* 06 — Project Details */}
      <div className="form-card">
        <SectionLabel number="06" label="Project Details" icon={<FileIcon />} />
        <p className="section-desc">Name and module path for your new project</p>
        <Flex direction="column" gap="4">
          <div>
            <Text as="label" htmlFor="field-module-name" size="2" weight="bold" style={fieldLabel}>
              Module Name
            </Text>
            <TextField.Root
              id="field-module-name"
              placeholder="github.com/your/module"
              value={moduleName}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                notifyInteract();
                setModuleName(e.target.value);
                setTouched(t => ({ ...t, moduleName: true }));
                setErrors(errs => ({ ...errs, moduleName: e.target.value.trim() ? undefined : 'Module Name is required.' }));
              }}
              onBlur={() => setTouched(t => ({ ...t, moduleName: true }))}
              color={errors.moduleName && touched.moduleName ? 'red' : undefined}
            />
            {errors.moduleName && touched.moduleName && <FieldError message={errors.moduleName} />}
          </div>
          <Flex gap="4" className="form-name-desc-row">
            <div style={{ flex: 1 }}>
              <Text as="label" htmlFor="field-project-name" size="2" weight="bold" style={fieldLabel}>
                Name
              </Text>
              <TextField.Root
                id="field-project-name"
                placeholder="my-app"
                value={name}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                  notifyInteract();
                  setName(e.target.value);
                  setTouched(t => ({ ...t, name: true }));
                  setErrors(errs => ({ ...errs, name: e.target.value.trim() ? undefined : 'Name is required.' }));
                }}
                onBlur={() => setTouched(t => ({ ...t, name: true }))}
                color={errors.name && touched.name ? 'red' : undefined}
              />
              {errors.name && touched.name && <FieldError message={errors.name} />}
            </div>
            <div style={{ flex: 1 }}>
              <Text as="label" htmlFor="field-description" size="2" weight="bold" style={fieldLabel}>
                Description
              </Text>
              <TextField.Root
                id="field-description"
                placeholder="Short project description"
                value={description}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                  notifyInteract();
                  setDescription(e.target.value);
                  setTouched(t => ({ ...t, description: true }));
                  setErrors(errs => ({ ...errs, description: e.target.value.trim() ? undefined : 'Description is required.' }));
                }}
                onBlur={() => setTouched(t => ({ ...t, description: true }))}
                color={errors.description && touched.description ? 'red' : undefined}
              />
              {errors.description && touched.description && <FieldError message={errors.description} />}
            </div>
          </Flex>
        </Flex>
      </div>

      {/* Generate Button */}
      <button
        type="button"
        className="generate-btn"
        disabled={isGenerating}
        onClick={handleGenerate}
        aria-label={`Generate project${isMac ? ' (⌘↵)' : ' (Ctrl+↵)'}`}
        style={{
          width: '100%',
          padding: '0.9rem 1.5rem',
          borderRadius: 'var(--radius-md)',
          border: 'none',
          background: isGenerating ? 'var(--color-surface-3)' : 'var(--color-accent)',
          color: isGenerating ? 'var(--color-text-muted)' : '#000000',
          fontSize: 15,
          fontWeight: 700,
          letterSpacing: '0.08em',
          cursor: isGenerating ? 'not-allowed' : 'pointer',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          gap: 10,
          transition: 'background 0.15s, color 0.15s',
          fontFamily: 'inherit',
          boxShadow: isGenerating ? 'none' : '0 0 0 1px rgba(255,215,0,0.3), var(--shadow-sm)',
          position: 'relative',
          overflow: 'hidden',
        }}
      >
        {isGenerating ? (
          <>
            <Spinner />
            GENERATING…
          </>
        ) : (
          <>
            GENERATE
            <span style={{ opacity: 0.55, fontSize: 12, fontWeight: 500, letterSpacing: '0.04em' }}>
              {isMac ? '⌘↵' : 'Ctrl+↵'}
            </span>
          </>
        )}
      </button>

      {/* Terminal-style error banner */}
      {generateError && (
        <ErrorBanner message={generateError} onDismiss={() => setGenerateError(null)} />
      )}

      {/* Success banner with animated checkmark + countdown progress bar */}
      {generateSuccess && (
        <SuccessBanner countdown={successCountdown} onDismiss={() => setGenerateSuccess(false)} />
      )}
    </div>
  );
};

export default GeneratorForm;

