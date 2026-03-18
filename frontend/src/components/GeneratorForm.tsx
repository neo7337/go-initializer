import React from 'react';
import { Card, Button, Text, Flex, Heading, RadioGroup, Checkbox, Callout } from '@radix-ui/themes';
import { useGeneratorForm } from '../hooks/useGeneratorForm';
import type { AddonCategory } from '../types';

const inputStyle = (hasError: boolean): React.CSSProperties => ({
  width: '100%',
  padding: '0.7rem',
  fontSize: 16,
  borderRadius: 'var(--radius-md)',
  border: hasError ? '2px solid var(--color-destructive)' : '1.5px solid var(--color-border)',
  background: 'var(--color-surface-3)',
  color: 'var(--color-text-primary)',
  fontWeight: 500,
  outline: hasError ? '2px solid var(--color-destructive)' : 'none',
  marginTop: 4,
});

const GeneratorForm: React.FC<{ onExplore: () => void }> = ({ onExplore }) => {
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
    isMac,
    handleGenerate,
  } = useGeneratorForm();

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 32 }}>
      {/* Go Version Card */}
      <Card style={{ marginBottom: 0 }}>
        <Heading size="4" mb="3">Go Version</Heading>
        <RadioGroup.Root
          value={goVersion}
          onValueChange={val => {
            setGoVersion(val);
            setTouched(t => ({ ...t, goVersion: true }));
            setErrors(errs => ({ ...errs, goVersion: val.trim() ? undefined : 'Go Version is required.' }));
          }}
          orientation="horizontal"
          style={{ display: 'flex', flexDirection: 'row', flexWrap: 'wrap', gap: 16 }}
        >
          {goVersionOptions.map(ver => (
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
            setTouched(t => ({ ...t, projectType: true }));
            setErrors(errs => ({ ...errs, projectType: val.trim() ? undefined : 'Project Type is required.' }));
          }}
          orientation="horizontal"
          style={{ display: 'flex', flexDirection: 'row', flexWrap: 'wrap', gap: 16 }}
        >
          {supportedProjectTypes.map(pt => {
            const value = pt.type.toLowerCase().replace(/ /g, '-');
            return (
              <RadioGroup.Item key={pt.type} value={value} style={{ marginRight: 0 }}>
                <Text as="span" size="3" weight={projectType === value ? 'bold' : 'regular'}>{pt.label}</Text>
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
            setTouched(t => ({ ...t, framework: true }));
            setErrors(errs => ({ ...errs, framework: val.trim() ? undefined : 'Framework/Dependency is required.' }));
          }}
          orientation="horizontal"
          style={{ display: 'flex', flexDirection: 'row', flexWrap: 'wrap', gap: 16 }}
        >
          {!currentFrameworkOptions.length ? (
            <Text as="span" size="3" weight="regular">No frameworks available for the selected project type.</Text>
          ) : currentFrameworkOptions.map(fw => (
            <RadioGroup.Item key={fw.value} value={fw.value} style={{ marginRight: 0 }}>
              <Text as="span" size="3" weight={framework === fw.value ? 'bold' : 'regular'}>{fw.label}</Text>
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
          {Object.entries(addonOptions).map(([category, opts]) => (
            <div key={category} style={{ minWidth: 180 }}>
              <Text as="div" size="3" weight="bold" mb="2" style={{ textTransform: 'capitalize' }}>
                {category === 'other' ? 'Other Libraries' : category.charAt(0).toUpperCase() + category.slice(1)}
              </Text>
              <Flex gap="4" direction="column">
                {opts.map(opt => (
                  <label key={opt.value} style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer' }}>
                    <Checkbox
                      checked={selectedAddons[category]?.includes(opt.value) ?? false}
                      onCheckedChange={() => handleAddonChange(category as AddonCategory, opt.value)}
                    />
                    <Text as="span" size="2">{opt.label}</Text>
                  </label>
                ))}
              </Flex>
            </div>
          ))}
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
              setTouched(t => ({ ...t, moduleName: true }));
              setErrors(errs => ({ ...errs, moduleName: e.target.value.trim() ? undefined : 'Module Name is required.' }));
            }}
            onBlur={() => setTouched(t => ({ ...t, moduleName: true }))}
            required
            style={inputStyle(!!(errors.moduleName && touched.moduleName))}
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
                  setTouched(t => ({ ...t, name: true }));
                  setErrors(errs => ({ ...errs, name: e.target.value.trim() ? undefined : 'Name is required.' }));
                }}
                onBlur={() => setTouched(t => ({ ...t, name: true }))}
                required
                style={inputStyle(!!(errors.name && touched.name))}
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
                  setTouched(t => ({ ...t, description: true }));
                  setErrors(errs => ({ ...errs, description: e.target.value.trim() ? undefined : 'Description is required.' }));
                }}
                onBlur={() => setTouched(t => ({ ...t, description: true }))}
                required
                style={inputStyle(!!(errors.description && touched.description))}
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
        <Button size="4" color="gray" variant="soft" onClick={onExplore}>
          EXPLORE
        </Button>
      </Flex>

      {/* Inline error banner (T22) */}
      {generateError && (
        <Callout.Root color="red" role="alert">
          <Callout.Text>
            <strong>Error: </strong>{generateError}
          </Callout.Text>
          <Button
            size="1"
            variant="ghost"
            color="red"
            style={{ marginLeft: 'auto' }}
            onClick={() => setGenerateError(null)}
            aria-label="Dismiss error"
          >
            ✕
          </Button>
        </Callout.Root>
      )}

      {/* Success banner (T23) */}
      {generateSuccess && (
        <Callout.Root color="green" role="status">
          <Callout.Text>
            <strong>Success! </strong>Your project zip has been downloaded.
          </Callout.Text>
          <Button
            size="1"
            variant="ghost"
            color="green"
            style={{ marginLeft: 'auto' }}
            onClick={() => setGenerateSuccess(false)}
            aria-label="Dismiss success"
          >
            ✕
          </Button>
        </Callout.Root>
      )}
    </div>
  );
};

export default GeneratorForm;
