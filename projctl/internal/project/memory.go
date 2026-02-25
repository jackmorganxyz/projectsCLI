package project

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MemoryEntry represents a parsed memory entry from a memory file.
type MemoryEntry struct {
	Heading string
	Content string
}

// LoadMemory reads and parses the memory file for a project.
func LoadMemory(projectDir string) ([]MemoryEntry, error) {
	memPath := filepath.Join(projectDir, "memory", "MEMORY.md")
	return ParseMemoryFile(memPath)
}

// ParseMemoryFile reads a markdown memory file and extracts entries by heading.
func ParseMemoryFile(path string) ([]MemoryEntry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read memory file: %w", err)
	}

	return parseMemoryEntries(string(data)), nil
}

// LatestMemoryEntries returns the last n entries from the memory file.
func LatestMemoryEntries(projectDir string, n int) ([]MemoryEntry, error) {
	entries, err := LoadMemory(projectDir)
	if err != nil {
		return nil, err
	}

	if n <= 0 || n >= len(entries) {
		return entries, nil
	}

	return entries[len(entries)-n:], nil
}

// parseMemoryEntries splits markdown content into entries by heading.
func parseMemoryEntries(content string) []MemoryEntry {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var entries []MemoryEntry
	var current *MemoryEntry
	var contentLines []string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "# ") || strings.HasPrefix(line, "## ") {
			// Save previous entry.
			if current != nil {
				current.Content = strings.TrimSpace(strings.Join(contentLines, "\n"))
				entries = append(entries, *current)
			}
			current = &MemoryEntry{
				Heading: strings.TrimSpace(strings.TrimLeft(line, "#")),
			}
			contentLines = nil
		} else if current != nil {
			contentLines = append(contentLines, line)
		}
	}

	// Save last entry.
	if current != nil {
		current.Content = strings.TrimSpace(strings.Join(contentLines, "\n"))
		entries = append(entries, *current)
	}

	return entries
}
