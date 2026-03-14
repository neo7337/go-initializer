// Converts supportedFrameworks object to Record<string, { label: string; value: string }[]> for UI
export function toSupportedFrameworkOptionsMap(supportedFrameworks: Record<string, Record<string, boolean>>): Record<string, { label: string; value: string }[]> {
  const labelMap: Record<string, string> = {
    golly: 'golly (recommended)',
    gin: 'Gin',
    echo: 'Echo',
    fiber: 'Fiber',
    gokit: 'Go kit',
    chi: 'Chi',
    cobra: 'Cobra',
    urfave: 'urfave/cli',
    kingpin: 'Kingpin',
  };
  const result: Record<string, { label: string; value: string }[]> = {};
  Object.entries(supportedFrameworks).forEach(([ptype, frameworks]) => {
    const arr = Object.keys(frameworks)
      .filter((fw) => frameworks[fw])
      .map((fw) => ({ 'label': (labelMap[fw] || fw), 'value': fw }));
    result[ptype] = arr;
  });
  return result;
}
// Converts an array of Go version strings to an array of objects with {version, label}
// Accepts an object like { "1.23.12": true, ... } and returns array of {version, label}
export function toGoVersionOptions(versions: Record<string, boolean>): { version: string; label: string }[] {
  return Object.keys(versions)
    .sort((a, b) => {
      // Sort descending by version (semver-like)
      const pa = a.split('.').map(Number);
      const pb = b.split('.').map(Number);
      for (let i = 0; i < Math.max(pa.length, pb.length); i++) {
        if ((pb[i] || 0) !== (pa[i] || 0)) return (pb[i] || 0) - (pa[i] || 0);
      }
      return 0;
    })
    .map((version, idx, arr) => ({
      version,
      label: idx === 0 ? `${version} (latest stable)` : version,
    }));
}

export function toSupportedProjectTypes(projectTypes: Record<string, string>): { type: string; label: string }[] {
  return Object.keys(projectTypes)
    .map((type) => ({
      type,
      label: projectTypes[type],
    }));
}

export function toSupportedAddons(supportedAddons: Record<string, Record<string, boolean>>): Record<string, string[]> {
    const labelMap: Record<string, string> = {}
    const result: Record<string, string[]> = {};
    Object.entries(supportedAddons).forEach(([category, addons]) => {
        const arr = Object.keys(addons)
            .filter((fw) => addons[fw])
            .map((fw) => (labelMap[fw] || fw));
        result[category] = arr;
    });
    return result;
}

/** Converts supportedAddons map from the API into labelled option pairs for the UI. */
export function toAddonOptions(
  supportedAddons: Record<string, Record<string, boolean>>,
): Record<string, { value: string; label: string }[]> {
  const labelMap: Record<string, string> = {
    redis: 'Redis',
    memcached: 'Memcached',
    gorm: 'Gorm',
    ent: 'Ent',
    zap: 'Zap',
    logrus: 'Logrus',
    cobra: 'Cobra',
  };
  const result: Record<string, { value: string; label: string }[]> = {};
  Object.entries(supportedAddons).forEach(([category, addons]) => {
    result[category] = Object.keys(addons)
      .filter(a => addons[a])
      .map(a => ({ value: a, label: labelMap[a] || a }));
  });
  return result;
}