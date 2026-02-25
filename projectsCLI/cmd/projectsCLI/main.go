package main

import (
	"fmt"
	"os"

	"github.com/jackmorganxyz/projectsCLI/internal/cli"
	"github.com/jackmorganxyz/projectsCLI/internal/config"
	"github.com/jackmorganxyz/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	defaultConfigPath, err := config.ConfigPath()
	if err != nil {
		defaultConfigPath = "~/.projects/config.toml"
	}

	var jsonOutput bool
	configPath := defaultConfigPath

	rootCmd := &cobra.Command{
		Use:          "projectsCLI",
		Short:        "Manage project scaffolds ✨",
		Long:         "projectsCLI manages project scaffolds under ~/.projects/projects/.\nAgents use it via projectsCLI <command> --json; humans get a polished TUI.",
		Version:      version,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, _ []string) {
			if !tui.IsJSON() && tui.IsInteractive() {
				fmt.Fprintln(cmd.OutOrStdout(), tui.Banner())
			}
			_ = cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// First-run setup: detect openclaw, let user choose project location.
			if config.NeedsSetup() && tui.IsInteractive() && !tui.IsJSON() {
				if _, err := config.RunSetup(); err != nil {
					return fmt.Errorf("setup: %w", err)
				}
			}

			if err := config.EnsureDirs(); err != nil {
				return fmt.Errorf("create directories: %w", err)
			}

			cfg, err := config.LoadFromPath(configPath)
			if err != nil {
				return fmt.Errorf("load config %q: %w", configPath, err)
			}

			runtime := cli.RuntimeContext{
				Config:     cfg,
				ConfigPath: configPath,
				JSON:       jsonOutput,
			}
			cmd.SetContext(cli.WithRuntimeContext(cmd.Context(), runtime))
			tui.SetJSON(jsonOutput)

			return nil
		},
	}

	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output JSON (auto-enabled when piped)")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", defaultConfigPath, "path to config file")

	rootCmd.AddCommand(
		cli.NewCreateCmd(),
		cli.NewListCmd(),
		cli.NewLoadCmd(),
		cli.NewDeleteCmd(),
		cli.NewViewCmd(),
		cli.NewEditCmd(),
		cli.NewStatusCmd(),
		cli.NewPushCmd(),
	)

	return rootCmd
}
