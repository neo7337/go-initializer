import { describe, it, expect } from 'vitest';
import {
  toSupportedFrameworkOptionsMap,
  toGoVersionOptions,
  toSupportedProjectTypes,
  toSupportedAddons,
  toAddonOptions,
} from './utils';

// ── toSupportedFrameworkOptionsMap ────────────────────────────────────────────

describe('toSupportedFrameworkOptionsMap', () => {
  it('maps known keys through the label map', () => {
    const result = toSupportedFrameworkOptionsMap({
      microservice: { gin: true, echo: true },
    });
    expect(result.microservice).toEqual([
      { label: 'Gin', value: 'gin' },
      { label: 'Echo', value: 'echo' },
    ]);
  });

  it('applies golly → "golly (recommended)"', () => {
    const result = toSupportedFrameworkOptionsMap({
      'simple-project': { golly: true },
    });
    expect(result['simple-project'][0].label).toBe('golly (recommended)');
  });

  it('falls back to the raw key for unknown framework names', () => {
    const result = toSupportedFrameworkOptionsMap({
      custom: { myfx: true },
    });
    expect(result.custom[0]).toEqual({ label: 'myfx', value: 'myfx' });
  });

  it('filters out false-valued frameworks', () => {
    const result = toSupportedFrameworkOptionsMap({
      microservice: { gin: true, echo: false, fiber: true },
    });
    expect(result.microservice.map(o => o.value)).toEqual(['gin', 'fiber']);
  });

  it('returns an empty array when all frameworks are false', () => {
    const result = toSupportedFrameworkOptionsMap({
      microservice: { gin: false },
    });
    expect(result.microservice).toEqual([]);
  });

  it('handles multiple project types independently', () => {
    const result = toSupportedFrameworkOptionsMap({
      microservice: { gin: true },
      'cli-app': { cobra: true },
    });
    expect(result.microservice[0].value).toBe('gin');
    expect(result['cli-app'][0].value).toBe('cobra');
  });
});

// ── toGoVersionOptions ────────────────────────────────────────────────────────

describe('toGoVersionOptions', () => {
  it('sorts versions in descending order', () => {
    const result = toGoVersionOptions({ '1.21': true, '1.23': true, '1.22': true });
    expect(result.map(v => v.version)).toEqual(['1.23', '1.22', '1.21']);
  });

  it('appends "(latest stable)" to the first (highest) version only', () => {
    const result = toGoVersionOptions({ '1.21': true, '1.23': true });
    expect(result[0].label).toMatch(/1\.23.*latest stable/);
    expect(result[1].label).toBe('1.21');
  });

  it('handles a single version input', () => {
    const result = toGoVersionOptions({ '1.22.5': true });
    expect(result).toHaveLength(1);
    expect(result[0].label).toMatch(/latest stable/);
  });

  it('handles patch-level version sorting', () => {
    const result = toGoVersionOptions({ '1.22.1': true, '1.22.10': true, '1.22.2': true });
    expect(result[0].version).toBe('1.22.10');
  });
});

// ── toSupportedProjectTypes ───────────────────────────────────────────────────

describe('toSupportedProjectTypes', () => {
  it('converts a record to {type, label}[] array', () => {
    const result = toSupportedProjectTypes({ microservice: 'Microservice', 'cli-app': 'CLI App' });
    expect(result).toEqual(
      expect.arrayContaining([
        { type: 'microservice', label: 'Microservice' },
        { type: 'cli-app', label: 'CLI App' },
      ]),
    );
  });

  it('returns an empty array for an empty input', () => {
    expect(toSupportedProjectTypes({})).toEqual([]);
  });
});

// ── toSupportedAddons ─────────────────────────────────────────────────────────

describe('toSupportedAddons', () => {
  it('returns only addons where the value is true', () => {
    const result = toSupportedAddons({
      cache: { redis: true, memcached: false },
      database: { gorm: true },
    });
    expect(result.cache).toEqual(['redis']);
    expect(result.database).toEqual(['gorm']);
  });

  it('returns empty arrays when all entries are false', () => {
    const result = toSupportedAddons({ cache: { redis: false } });
    expect(result.cache).toEqual([]);
  });

  it('handles an empty category map', () => {
    const result = toSupportedAddons({ cache: {} });
    expect(result.cache).toEqual([]);
  });
});

// ── toAddonOptions ────────────────────────────────────────────────────────────

describe('toAddonOptions', () => {
  it('applies the label map for known keys', () => {
    const result = toAddonOptions({ cache: { redis: true } });
    expect(result.cache[0]).toEqual({ value: 'redis', label: 'Redis' });
  });

  it('applies label map: memcached, gorm, ent, zap, logrus, cobra, urfave, kingpin', () => {
    const input = { other: { zap: true, logrus: true, gorm: true } };
    const result = toAddonOptions(input);
    const labels = result.other.map(o => o.label);
    expect(labels).toContain('Zap');
    expect(labels).toContain('Logrus');
    expect(labels).toContain('Gorm');
  });

  it('falls back to the raw key for unknown addon names', () => {
    const result = toAddonOptions({ other: { myaddon: true } });
    expect(result.other[0]).toEqual({ value: 'myaddon', label: 'myaddon' });
  });

  it('filters out false entries', () => {
    const result = toAddonOptions({ cache: { redis: true, memcached: false } });
    expect(result.cache.map(o => o.value)).toEqual(['redis']);
  });

  it('returns empty array for a category with no true entries', () => {
    const result = toAddonOptions({ cache: { redis: false } });
    expect(result.cache).toEqual([]);
  });
});
