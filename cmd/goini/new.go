package main

import (
	"errors"
	"fmt"
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

	// TODO T34: call generator, extract zip into output directory.
	// TODO T35: print success banner with next-step hints.
	_ = req
	_ = opts.output
	fmt.Println("(prompts complete — generation wired in T34)")
	return nil
}

// promptMissing prompts interactively for any field in opts that was not
// already supplied via a flag. Framework options are filtered by the selected
// project type using SupportedFrameworksMap.
//
// The prompts are split into two groups so that the project type is known
// before the framework select is rendered:
//
//	Group 1 — name, module, description, Go version, project type
//	Group 2 — framework (filtered), cache/database/other addons, docker, output
func promptMissing(cmd *cobra.Command, opts *newOpts) error {
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

	// Description is optional; prompt only when the flag was not explicitly set.
	if !cmd.Flags().Changed("description") {
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
		cacheAddons  []string
		dbAddons     []string
		otherAddons  []string
		promptAddons = !cmd.Flags().Changed("addon")
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

	// Docker: --docker defaults to false so it is never "Changed" unless the
	// user explicitly passed it; always prompt when not set.
	if !cmd.Flags().Changed("docker") {
		group2 = append(group2, huh.NewConfirm().
			Title("Docker support").
			Value(&opts.docker))
	}

	if opts.output == "" {
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
