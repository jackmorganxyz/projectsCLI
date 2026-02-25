package cli

import (
	"fmt"

	"github.com/jackmorganxyz/projectsCLI/internal/git"
	"github.com/jackmorganxyz/projectsCLI/internal/project"
	"github.com/jackmorganxyz/projectsCLI/internal/tui"
	"github.com/spf13/cobra"
)

// NewPushCmd handles the full git push workflow.
func NewPushCmd() *cobra.Command {
	var (
		message string
		private bool
		noGH    bool
	)

	cmd := &cobra.Command{
		Use:   "push <slug>",
		Short: "Push project to git remote",
		Long:  "Stage, commit, and push changes. Creates GitHub repo if no remote exists.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			runtime, ok := RuntimeFromContext(cmd.Context())
			if !ok {
				return fmt.Errorf("missing runtime context")
			}

			slug := args[0]
			proj, err := project.FindProject(runtime.Config.ProjectsDir, slug)
			if err != nil {
				return err
			}
			dir := proj.Dir

			// Ensure git is initialized.
			if !git.IsRepo(dir) {
				fmt.Fprintln(cmd.ErrOrStderr(), tui.Muted("Setting up git... first commits are special."))
				if err := git.Init(dir); err != nil {
					return fmt.Errorf("git init: %w", err)
				}
			}

			// Stage and commit.
			if err := git.AddAll(dir); err != nil {
				return fmt.Errorf("git add: %w", err)
			}

			hasChanges, err := git.HasUncommitted(dir)
			if err != nil {
				return fmt.Errorf("check changes: %w", err)
			}

			if hasChanges {
				if message == "" {
					message = "Update project"
				}
				if err := git.Commit(dir, message); err != nil {
					return fmt.Errorf("git commit: %w", err)
				}
				fmt.Fprintln(cmd.ErrOrStderr(), tui.SuccessMessage("Changes committed. "+tui.RandomPushCheer()))
			} else {
				fmt.Fprintln(cmd.ErrOrStderr(), tui.Muted(tui.RandomNoChanges()))
			}

			// Create remote if needed.
			if !git.HasRemote(dir) && !noGH {
				if !git.HasGHCLI() {
					return fmt.Errorf("no remote configured and gh CLI not available; add a remote manually or install gh")
				}

				org := runtime.Config.GitHubOrg
				fmt.Fprintln(cmd.ErrOrStderr(), tui.Muted("Creating GitHub repo... your code deserves a home."))
				repoURL, err := git.CreateRepo(dir, slug, org, private)
				if err != nil {
					return fmt.Errorf("create repo: %w", err)
				}

				// Update project metadata with remote URL.
				proj.Meta.GitRemote = repoURL
				if err := project.WriteProjectFile(proj.Dir, proj.Meta, proj.Body); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: failed to save remote URL to PROJECT.md: %v\n", err)
				}

				fmt.Fprintln(cmd.ErrOrStderr(), tui.SuccessMessage(fmt.Sprintf("Repository created: %s", repoURL)))
			} else if git.HasRemote(dir) {
				// Push to existing remote.
				branch, _ := git.CurrentBranch(dir)
				if branch == "" {
					branch = "main"
				}
				if err := git.PushSetUpstream(dir, "origin", branch); err != nil {
					return fmt.Errorf("git push: %w", err)
				}
				fmt.Fprintln(cmd.ErrOrStderr(), tui.SuccessMessage("Pushed to remote. "+tui.RandomPushCheer()))
				if tip := tui.MaybeTip(); tip != "" {
					fmt.Fprintln(cmd.ErrOrStderr(), tip)
				}
			}

			if tui.IsJSON() {
				remote, _ := git.RemoteURL(dir)
				return writeJSON(cmd.OutOrStdout(), map[string]any{
					"status": "pushed",
					"slug":   slug,
					"remote": remote,
				})
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "commit message")
	cmd.Flags().BoolVar(&private, "private", true, "create private GitHub repo")
	cmd.Flags().BoolVar(&noGH, "no-github", false, "skip GitHub repo creation")

	return cmd
}
