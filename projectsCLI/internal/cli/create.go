package cli

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/jackmorganxyz/projectsCLI/internal/agent"
	"github.com/jackmorganxyz/projectsCLI/internal/git"
	"github.com/jackmorganxyz/projectsCLI/internal/project"
	"github.com/jackmorganxyz/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

// NewCreateCmd creates a new project scaffold.
func NewCreateCmd() *cobra.Command {
	var (
		title       string
		description string
		tags        []string
		status      string
	)

	cmd := &cobra.Command{
		Use:   "create [slug]",
		Short: "Create a new project",
		Long: `Create a new project scaffold with directory structure and template files.

If slug is omitted, it is generated from --title.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runtime, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			var slug string
			if len(args) == 1 {
				slug = args[0]
			} else if title != "" {
				slug = Slugify(title)
				if slug == "" {
					return fmt.Errorf("could not generate a valid slug from title %q", title)
				}
			} else {
				return fmt.Errorf("provide a slug argument or --title to auto-generate one")
			}

			if err := ValidateSlug(slug); err != nil {
				return err
			}

			meta := project.NewMeta(slug, title)
			if description != "" {
				meta.Description = description
			}
			if len(tags) > 0 {
				meta.Tags = tags
			}
			if status != "" {
				meta.Status = status
			}

			// Determine the target directory.
			projectsDir := runtime.Config.ProjectsDir
			if runtime.Folder != "" {
				if runtime.Config.FolderByName(runtime.Folder) == nil {
					return fmt.Errorf("folder %q not configured; run 'projects folder add %s --account <gh-user>' first", runtime.Folder, runtime.Folder)
				}
				projectsDir = filepath.Join(runtime.Config.ProjectsDir, runtime.Folder)
			}

			dir, err := project.Scaffold(projectsDir, meta)
			if err != nil {
				return err
			}

			// Auto-init git if configured.
			if runtime.Config.AutoGitInit {
				if err := git.Init(dir); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: git init failed: %v\n", err)
				} else {
					_ = git.AddAll(dir)
					_ = git.Commit(dir, "Initial project scaffold")
				}
			}

			// Regenerate registry.
			_ = project.WriteRegistry(runtime.Config.ProjectsDir)

			if tui.IsJSON() {
				result := map[string]any{
					"status":     "created",
					"slug":       slug,
					"dir":        dir,
					"created_at": meta.CreatedAt,
				}
				if runtime.Folder != "" {
					result["folder"] = runtime.Folder
				}
				return writeJSON(cmd.OutOrStdout(), result)
			}

			w := cmd.OutOrStdout()
			fmt.Fprintln(w, tui.SuccessMessage(fmt.Sprintf("Created project %q — %s", slug, tui.RandomCreateCheer())))
			fmt.Fprintln(w, tui.FormatField("Directory", dir))
			if runtime.Folder != "" {
				fmt.Fprintln(w, tui.FormatField("Folder", runtime.Folder))
			}
			fmt.Fprintln(w, tui.FormatField("Created", time.Now().Format("2006-01-02")))
			if tip := tui.MaybeTip(); tip != "" {
				fmt.Fprintln(w)
				fmt.Fprintln(w, tip)
			}

			// Offer to spawn an AI agent to fill out the scaffold.
			if tui.IsInteractive() && agent.HasAny() {
				fmt.Fprintln(w)
				if err := offerAgentSpawn(cmd, dir, slug, description); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "project title (defaults to slug)")
	cmd.Flags().StringVar(&description, "description", "", "project description")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "project tags (comma-separated)")
	cmd.Flags().StringVar(&status, "status", "active", "project status")

	return cmd
}

// offerAgentSpawn prompts the user to optionally launch an AI agent in the
// newly scaffolded project directory.
func offerAgentSpawn(cmd *cobra.Command, dir, slug, description string) error {
	agents := agent.Detect()
	if len(agents) == 0 {
		return nil
	}

	options := []huh.Option[string]{
		huh.NewOption("Skip — I'll work on it myself", "skip"),
	}
	for _, a := range agents {
		options = append(options, huh.NewOption(fmt.Sprintf("Launch %s", a.Name), a.Command))
	}

	var selected string
	theme := huh.ThemeBase()
	theme.Focused.Title = theme.Focused.Title.Foreground(lipgloss.Color(tui.ColorPrimary))
	theme.Focused.SelectSelector = theme.Focused.SelectSelector.Foreground(lipgloss.Color(tui.ColorPrimary))

	err := huh.NewSelect[string]().
		Title("Spawn an AI agent to fill out the project scaffold?").
		Options(options...).
		Value(&selected).
		WithTheme(theme).
		Run()

	if err != nil || selected == "skip" {
		return nil
	}

	var chosenAgent *agent.Agent
	for _, a := range agents {
		if a.Command == selected {
			chosenAgent = &a
			break
		}
	}
	if chosenAgent == nil {
		return nil
	}

	var userPrompt string
	err = huh.NewText().
		Title("What should the agent work on?").
		Description("Describe what you want the agent to do with the scaffolded files").
		Value(&userPrompt).
		WithTheme(theme).
		Run()

	if err != nil || strings.TrimSpace(userPrompt) == "" {
		return nil
	}

	fullPrompt := fmt.Sprintf(
		"You are working in a new project scaffold. "+
			"The project is called %q. Description: %s. "+
			"Read USAGE.md to understand the project structure, then: %s",
		slug, description, userPrompt,
	)

	w := cmd.OutOrStdout()
	fmt.Fprintln(w, tui.Muted(fmt.Sprintf("Launching %s...", chosenAgent.Name)))

	if err := agent.Spawn(*chosenAgent, dir, fullPrompt); err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), tui.WarningMessage(
			fmt.Sprintf("%s exited with error: %v", chosenAgent.Name, err)))
	}

	return nil
}
