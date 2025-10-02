// Converts supportedFrameworks object to Record<string, string[]> for UI
export function toSupportedFrameworkOptionsMap(supportedFrameworks: Record<string, Record<string, boolean>>): Record<string, string[]> {
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
  const result: Record<string, string[]> = {};
  Object.entries(supportedFrameworks).forEach(([ptype, frameworks]) => {
    const arr = Object.keys(frameworks)
      .filter((fw) => frameworks[fw])
      .map((fw) => (labelMap[fw] || fw));
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
    const labelMap: Record<string, string> = {

    }
    const result: Record<string, string[]> = {};
    Object.entries(supportedAddons).forEach(([category, addons]) => {
        const arr = Object.keys(addons)
            .filter((fw) => addons[fw])
            .map((fw) => (labelMap[fw] || fw));
        result[category] = arr;
    });
    return result;
}