package agent

import (
	"fmt"
	"os"
	"os/exec"
)

// AgentType distinguishes between supported AI agents.
type AgentType int

const (
	ClaudeCode AgentType = iota
	CodexCLI
)

// Agent represents an AI coding agent.
type Agent struct {
	Name    string
	Command string
	Type    AgentType
}

var knownAgents = []Agent{
	{Name: "Claude Code", Command: "claude", Type: ClaudeCode},
	{Name: "Codex CLI", Command: "codex", Type: CodexCLI},
}

// Detect returns a list of AI agents installed on the current system.
func Detect() []Agent {
	var found []Agent
	for _, a := range knownAgents {
		if _, err := exec.LookPath(a.Command); err == nil {
			found = append(found, a)
		}
	}
	return found
}

// HasAny reports whether at least one AI agent is installed.
func HasAny() bool {
	return len(Detect()) > 0
}

// Spawn runs the given agent in workDir with the given prompt.
// stdin/stdout/stderr are attached so the agent gets full terminal control.
func Spawn(a Agent, workDir string, prompt string) error {
	cmd := exec.Command(a.Command, prompt)
	cmd.Dir = workDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SpawnWithFile runs the given agent in workDir with a prompt that references
// a specific file path.
func SpawnWithFile(a Agent, workDir string, filePath string, prompt string) error {
	fullPrompt := fmt.Sprintf("Edit the file %s: %s", filePath, prompt)
	return Spawn(a, workDir, fullPrompt)
}
