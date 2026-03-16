package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is set at build time via ldflags:
// go build -ldflags "-X main.version=v1.0.0" ./cmd/goini
var version = "dev"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:          "goini",
		Short:        "Go project scaffolding tool",
		Long:         "goini scaffolds production-ready Go projects from the command line.",
		SilenceUsage: false,
	}
	root.AddCommand(newVersionCmd())
	root.AddCommand(newCompletionCmd(root))
	root.AddCommand(newListCmd())
	root.AddCommand(newNewCmd())
	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version and build info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("goini %s\n", version)
		},
	}
}

func newCompletionCmd(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:       "completion [bash|zsh|fish|powershell]",
		Short:     "Generate shell completion script",
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				root.GenBashCompletion(os.Stdout)
			case "zsh":
				root.GenZshCompletion(os.Stdout)
			case "fish":
				root.GenFishCompletion(os.Stdout, true)
			case "powershell":
				root.GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
}
