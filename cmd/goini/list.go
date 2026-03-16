package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/neo7337/go-initializer/internal/generator"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List supported project types, frameworks, or addons",
	}
	cmd.AddCommand(newListTypesCmd())
	cmd.AddCommand(newListFrameworksCmd())
	cmd.AddCommand(newListAddonsCmd())
	return cmd
}

func newListTypesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "types",
		Short: "Print all supported project types",
		Run: func(cmd *cobra.Command, args []string) {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "TYPE\tLABEL")
			fmt.Fprintln(w, "----\t-----")
			keys := make([]string, 0, len(generator.SupportedProjectTypesLabelsMap))
			for k := range generator.SupportedProjectTypesLabelsMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Fprintf(w, "%s\t%s\n", k, generator.SupportedProjectTypesLabelsMap[k])
			}
			w.Flush()
		},
	}
}

func newListFrameworksCmd() *cobra.Command {
	var projectType string
	cmd := &cobra.Command{
		Use:   "frameworks",
		Short: "Print supported frameworks for a project type",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectType == "" {
				return fmt.Errorf("--type is required (e.g. --type microservice)")
			}
			frameworks, ok := generator.SupportedFrameworksMap[projectType]
			if !ok {
				return fmt.Errorf("unknown project type %q; run 'goini list types' to see valid values", projectType)
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "FRAMEWORK")
			fmt.Fprintln(w, "---------")
			keys := make([]string, 0, len(frameworks))
			for k := range frameworks {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Fprintln(w, k)
			}
			w.Flush()
			return nil
		},
	}
	cmd.Flags().StringVarP(&projectType, "type", "t", "", "Project type (e.g. microservice, api-server, cli-app, simple-project)")
	return cmd
}

func newListAddonsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "addons",
		Short: "Print all supported addons",
		Run: func(cmd *cobra.Command, args []string) {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "CATEGORY\tADDON")
			fmt.Fprintln(w, "--------\t-----")
			categories := make([]string, 0, len(generator.SupportedAddonsMap))
			for k := range generator.SupportedAddonsMap {
				categories = append(categories, k)
			}
			sort.Strings(categories)
			for _, cat := range categories {
				addons := make([]string, 0, len(generator.SupportedAddonsMap[cat]))
				for a := range generator.SupportedAddonsMap[cat] {
					addons = append(addons, a)
				}
				sort.Strings(addons)
				for _, a := range addons {
					fmt.Fprintf(w, "%s\t%s\n", cat, a)
				}
			}
			w.Flush()
		},
	}
}
