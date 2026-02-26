package editor

import "testing"

func TestHasCLI(t *testing.T) {
	// "ls" should exist on every Unix-like system.
	if !HasCLI("ls") {
		t.Error("expected HasCLI(\"ls\") to be true")
	}

	// A nonsense name should not be found.
	if HasCLI("__nonexistent_editor_xyz__") {
		t.Error("expected HasCLI for a random name to be false")
	}
}

func TestFindByCommand(t *testing.T) {
	editors := []Editor{
		{Name: "Vim", Command: "vim", Type: Terminal},
		{Name: "VS Code", Command: "code", Type: GUI},
	}

	found := FindByCommand(editors, "code")
	if found == nil {
		t.Fatal("expected to find editor with command \"code\"")
	}
	if found.Name != "VS Code" {
		t.Errorf("expected Name \"VS Code\", got %q", found.Name)
	}

	if FindByCommand(editors, "missing") != nil {
		t.Error("expected nil for a missing command")
	}
}

func TestDetectReturnsEditors(t *testing.T) {
	editors := Detect()
	// On a dev machine we should find at least one editor.
	// Skip on CI / bare containers where nothing may be installed.
	if len(editors) == 0 {
		t.Skip("no editors detected (likely CI environment)")
	}
	for _, ed := range editors {
		if ed.Name == "" || ed.Command == "" {
			t.Errorf("editor has empty name or command: %+v", ed)
		}
	}
}
