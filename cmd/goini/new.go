package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/neo7337/go-initializer/internal/generator"
)

// newOpts holds the values collected from `goini new` flags.
// Any field left as its zero value was not supplied via flag and will be
// filled in by the interactive prompts in promptMissing.
type newOpts struct {
	name        string
	module      string
	description string
	goVersion   string
	projectType string
	framework   string
	addons      []string // raw "category=value" strings from --addon flags
	docker      bool
	output      string
}

func newNewCmd() *cobra.Command {
	opts := &newOpts{}

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Scaffold a new Go project",
		Long: `Scaffold a new Go project interactively or fully via flags.

Any flag provided skips its corresponding interactive prompt, making the
command scriptable in CI pipelines.

Examples:
  goini new
  goini new --name myapp --type microservice --framework gin
  goini new --name myapp --module github.com/acme/myapp --type microservice \
            --framework gin --go-version 1.25.0 --addon cache=redis --docker`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNew(cmd, opts)
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "Project name")
	cmd.Flags().StringVar(&opts.module, "module", "", "Go module path (e.g. github.com/acme/myapp)")
	cmd.Flags().StringVar(&opts.description, "description", "", "Short project description")
	cmd.Flags().StringVar(&opts.goVersion, "go-version", "", "Go version (e.g. 1.25.0)")
	cmd.Flags().StringVar(&opts.projectType, "type", "", "Project type: microservice, simple-project, cli-app, api-server")
	cmd.Flags().StringVar(&opts.framework, "framework", "", "Framework name (must be valid for the selected project type)")
	cmd.Flags().StringArrayVar(&opts.addons, "addon", nil, "Addon in category=value format, repeatable (e.g. --addon cache=redis --addon database=gorm)")
	cmd.Flags().BoolVar(&opts.docker, "docker", false, "Generate a multi-stage Dockerfile")
	cmd.Flags().StringVar(&opts.output, "output", "", "Output directory (default: ./<name>)")

	return cmd
}

func runNew(cmd *cobra.Command, opts *newOpts) error {
	// Fill in any fields that were not provided via flags using interactive prompts.
	if err := promptMissing(cmd, opts); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Fprintln(cmd.ErrOrStderr(), "Aborted.")
			return nil
		}
		return err
	}

	// Parse --addon flags (and interactively-collected addon strings) into the
	// map[string][]string format expected by CreateProjectRequest.
	addonsMap, err := parseAddons(opts.addons)
	if err != nil {
		return err
	}

	req := generator.CreateProjectRequest{
		Name:          opts.name,
		ModuleName:    opts.module,
		Description:   opts.description,
		GoVersion:     opts.goVersion,
		ProjectType:   opts.projectType,
		Framework:     opts.framework,
		Addons:        addonsMap,
		DockerSupport: opts.docker,
	}

	// Look up and run the generator for the chosen project type.
	gen, ok := generator.GeneratorRegistry[req.ProjectType]
	if !ok {
		return fmt.Errorf("unknown project type %q; run 'goini list types' to see valid values", req.ProjectType)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Generating %s project…\n", req.ProjectType)

	zipBuf, err := gen.Generate(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	// Resolve the absolute output path.
	outDir, err := filepath.Abs(opts.output)
	if err != nil {
		return fmt.Errorf("resolving output path: %w", err)
	}

	// Extract the zip (whose entries are rooted at <name>/) into outDir,
	// stripping the leading <name>/ prefix so files land directly in outDir.
	if err := extractZip(zipBuf, outDir); err != nil {
		return fmt.Errorf("extracting project: %w", err)
	}

	printSuccess(cmd, req.Name, req.ProjectType, outDir)
	return nil
}

// printSuccess prints the success banner and contextual next-step hints.
// The make target differs by project type: cli-app projects build a binary,
// all other types run the service directly.
func printSuccess(cmd *cobra.Command, name, projectType, outDir string) {
	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "\nProject created at %s\n", outDir)
	fmt.Fprintln(w, "\nNext steps:")
	fmt.Fprintf(w, "  cd %s\n", name)
	fmt.Fprintln(w, "  go mod tidy")
	if projectType == "cli-app" {
		fmt.Fprintln(w, "  make build")
	} else {
		fmt.Fprintln(w, "  make run")
	}
}

// maxExtractSize is the maximum number of bytes copied from a single zip entry
// to guard against decompression bombs (G110).
const maxExtractSize = 100 * 1024 * 1024 // 100 MiB

// extractZip extracts all entries from buf into outDir, stripping the first
// path component of each zip entry (the generator always nests files under a
// top-level <name>/ folder).
//
// os.Root scopes every file operation to outDir, preventing directory traversal.
func extractZip(buf *bytes.Buffer, outDir string) error {
	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return fmt.Errorf("read zip: %w", err)
	}

	// Ensure the output directory exists.
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	// OpenRoot scopes all subsequent writes to outDir; no path can escape it.
	root, err := os.OpenRoot(outDir)
	if err != nil {
		return fmt.Errorf("open output root: %w", err)
	}
	defer root.Close()

	for _, f := range r.File {
		// Strip the leading path component (e.g. "myapp/cmd/main.go" → "cmd/main.go").
		stripped := stripFirstComponent(f.Name)
		if stripped == "" {
			// The entry IS the top-level directory itself; skip it.
			continue
		}

		entryPath := filepath.FromSlash(stripped)

		if f.FileInfo().IsDir() {
			if err := root.MkdirAll(entryPath, 0o750); err != nil {
				return err
			}
			continue
		}

		if dir := filepath.Dir(entryPath); dir != "." {
			if err := root.MkdirAll(dir, 0o750); err != nil {
				return err
			}
		}

		if err := writeZipEntry(f, root, entryPath); err != nil {
			return err
		}
	}
	return nil
}

// writeZipEntry writes a single zip file entry into root at entryPath.
// The copy is capped at maxExtractSize to prevent decompression bombs.
func writeZipEntry(f *zip.File, root *os.Root, entryPath string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := root.OpenFile(entryPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, io.LimitReader(rc, maxExtractSize))
	return err
}

// stripFirstComponent removes the first path segment from a forward-slash
// separated path (e.g. "myapp/cmd/main.go" → "cmd/main.go").
func stripFirstComponent(path string) string {
	idx := strings.Index(path, "/")
	if idx < 0 {
		return ""
	}
	return path[idx+1:]
}

// stdinIsTerminal reports whether os.Stdin is an interactive terminal.
// When false (e.g. in CI or when stdin is redirected) optional prompts are
// skipped and their zero-value defaults are used instead.
func stdinIsTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// promptMissing prompts interactively for any field in opts that was not
// already supplied via a flag. Framework options are filtered by the selected
// project type using SupportedFrameworksMap.
//
// Required fields (name, module, goVersion, projectType, framework) always
// trigger a prompt when missing; the command fails in non-TTY environments if
// these are not supplied via flags.
//
// Optional fields (description, addons, docker, output) only prompt when stdin
// is a TTY; in non-TTY environments they default silently to their zero values.
//
// The prompts are split into two groups so that the project type is known
// before the framework select is rendered:
//
//	Group 1 — name, module, description, Go version, project type
//	Group 2 — framework (filtered), cache/database/other addons, docker, output
func promptMissing(cmd *cobra.Command, opts *newOpts) error {
	isTTY := stdinIsTerminal()
	// ── Group 1: basic metadata + project type ────────────────────────────────
	var group1 []huh.Field

	if opts.name == "" {
		group1 = append(group1, huh.NewInput().
			Title("Project name").
			Value(&opts.name).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return errors.New("project name is required")
				}
				return nil
			}))
	}

	if opts.module == "" {
		group1 = append(group1, huh.NewInput().
			Title("Module path").
			Description("e.g. github.com/acme/myapp").
			Value(&opts.module).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return errors.New("module path is required")
				}
				return nil
			}))
	}

	// Description is optional; prompt only when the flag was not explicitly set
	// AND stdin is a TTY (non-TTY: default to empty string silently).
	if !cmd.Flags().Changed("description") && isTTY {
		group1 = append(group1, huh.NewInput().
			Title("Description").
			Description("Optional — press Enter to skip").
			Value(&opts.description))
	}

	if opts.goVersion == "" {
		group1 = append(group1, huh.NewSelect[string]().
			Title("Go version").
			Options(buildSelectOptsDesc(generator.SupportedGoVersionsMap)...).
			Value(&opts.goVersion))
	}

	if opts.projectType == "" {
		group1 = append(group1, huh.NewSelect[string]().
			Title("Project type").
			Options(buildLabelSelectOpts(generator.SupportedProjectTypesLabelsMap)...).
			Value(&opts.projectType))
	}

	// ── Group 2: framework + addons + docker + output ─────────────────────────
	// Framework options are computed dynamically from opts.projectType so that
	// the correct framework list is shown whether projectType came from a flag
	// or was chosen in Group 1.
	var (
		cacheAddons []string
		dbAddons    []string
		otherAddons []string
		// Prompt for addons only when neither --addon flag was set AND stdin is
		// a TTY; in non-TTY / fully-flagged CI runs default to no addons.
		promptAddons = !cmd.Flags().Changed("addon") && isTTY
	)
	var group2 []huh.Field

	if opts.framework == "" {
		group2 = append(group2, huh.NewSelect[string]().
			Title("Framework").
			OptionsFunc(func() []huh.Option[string] {
				return buildSelectOpts(generator.SupportedFrameworksMap[opts.projectType])
			}, &opts.projectType).
			Value(&opts.framework))
	}

	if promptAddons {
		if cacheOpts := buildSelectOpts(generator.SupportedAddonsMap["cache"]); len(cacheOpts) > 0 {
			group2 = append(group2, huh.NewMultiSelect[string]().
				Title("Cache addons").
				Options(cacheOpts...).
				Value(&cacheAddons))
		}
		if dbOpts := buildSelectOpts(generator.SupportedAddonsMap["database"]); len(dbOpts) > 0 {
			group2 = append(group2, huh.NewMultiSelect[string]().
				Title("Database addons").
				Options(dbOpts...).
				Value(&dbAddons))
		}
		if otherOpts := buildSelectOpts(generator.SupportedAddonsMap["other"]); len(otherOpts) > 0 {
			group2 = append(group2, huh.NewMultiSelect[string]().
				Title("Other addons").
				Options(otherOpts...).
				Value(&otherAddons))
		}
	}

	// Docker: prompt only when not explicitly flagged AND stdin is a TTY.
	if !cmd.Flags().Changed("docker") && isTTY {
		group2 = append(group2, huh.NewConfirm().
			Title("Docker support").
			Value(&opts.docker))
	}

	if opts.output == "" && isTTY {
		group2 = append(group2, huh.NewInput().
			Title("Output directory").
			Description("Leave blank to use ./<name>").
			Value(&opts.output))
	}

	// ── Run the form ──────────────────────────────────────────────────────────
	var groups []*huh.Group
	if len(group1) > 0 {
		groups = append(groups, huh.NewGroup(group1...))
	}
	if len(group2) > 0 {
		groups = append(groups, huh.NewGroup(group2...))
	}

	if len(groups) > 0 {
		if err := huh.NewForm(groups...).Run(); err != nil {
			return err
		}
	}

	// Merge interactively collected addons back into opts.addons so that the
	// downstream parseAddons call handles everything uniformly.
	if promptAddons {
		for _, a := range cacheAddons {
			opts.addons = append(opts.addons, "cache="+a)
		}
		for _, a := range dbAddons {
			opts.addons = append(opts.addons, "database="+a)
		}
		for _, a := range otherAddons {
			opts.addons = append(opts.addons, "other="+a)
		}
	}

	// Default output directory to ./<name> when the user left it blank.
	if opts.output == "" {
		opts.output = "./" + opts.name
	}

	return nil
}

// buildSelectOpts returns a sorted slice of huh.Option[string] from a
// map[string]bool, using each key as both the display label and stored value.
func buildSelectOpts(m map[string]bool) []huh.Option[string] {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	opts := make([]huh.Option[string], len(keys))
	for i, k := range keys {
		opts[i] = huh.NewOption(k, k)
	}
	return opts
}

// buildSelectOptsDesc is like buildSelectOpts but sorted in descending order,
// which is used to show the newest Go version first.
func buildSelectOptsDesc(m map[string]bool) []huh.Option[string] {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	opts := make([]huh.Option[string], len(keys))
	for i, k := range keys {
		opts[i] = huh.NewOption(k, k)
	}
	return opts
}

// buildLabelSelectOpts creates sorted huh.Option[string] from a map[key]label
// where the label is the user-visible display string and the key is the stored value.
func buildLabelSelectOpts(m map[string]string) []huh.Option[string] {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	opts := make([]huh.Option[string], len(keys))
	for i, k := range keys {
		opts[i] = huh.NewOption(m[k], k)
	}
	return opts
}

// parseAddons converts a slice of "category=value" strings (from --addon flags
// and interactive prompt selections) into the map[string][]string format used
// by CreateProjectRequest.Addons.
func parseAddons(raw []string) (map[string][]string, error) {
	result := make(map[string][]string)
	for _, entry := range raw {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid --addon value %q: expected category=value (e.g. --addon cache=redis)", entry)
		}
		result[parts[0]] = append(result[parts[0]], parts[1])
	}
	return result, nil
}
