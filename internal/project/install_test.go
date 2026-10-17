package project

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MDC15/stack-agents"
	"github.com/MDC15/stack-agents/internal/generate"
)

func TestInstallFresh(t *testing.T) {
	root := t.TempDir()
	p, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	plan := map[string]bool{"claude": true, "copilot": true}

	actions, err := p.Install(stackagents.Assets, false, plan)
	if err != nil {
		t.Fatalf("Install: %v", err)
	}
	if len(actions) == 0 {
		t.Error("Install returned no actions")
	}

	for _, f := range []string{
		"AGENTS.md", "VERSION", "scripts/load_skill.py",
		".ai/principles.md", ".ai/skills/tdd.md",
	} {
		if _, err := os.Stat(filepath.Join(root, f)); err != nil {
			t.Errorf("%s missing after Install: %v", f, err)
		}
	}

	drifts, err := p.Check(plan)
	if err != nil {
		t.Fatal(err)
	}
	if len(drifts) != 0 {
		t.Errorf("Check after Install = %+v, want no drift", drifts)
	}
}

func TestInstallRefusesExisting(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "VERSION"), []byte("9.9.9\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.Install(stackagents.Assets, false, nil)
	if err == nil {
		t.Fatal("Install with VERSION present and no force = nil error, want an already-installed error")
	}
	if !strings.Contains(err.Error(), "already installed") {
		t.Errorf("Install error = %q, want it to mention %q", err, "already installed")
	}
}

func TestInstallAdoptsExistingAGENTS(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# mine\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := p.Install(stackagents.Assets, false, nil); err != nil {
		t.Fatalf("Install adopting an existing AGENTS.md: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "# mine\n"; string(got) != want {
		t.Errorf("AGENTS.md = %q, want the adopted %q untouched", got, want)
	}

	if _, err := os.Stat(filepath.Join(root, "VERSION")); err != nil {
		t.Errorf("VERSION not shipped during adoption: %v", err)
	}

	rendered, err := os.ReadFile(filepath.Join(root, ".agents", "claude.md"))
	if err != nil {
		t.Fatalf(".agents/claude.md missing after adoption: %v", err)
	}
	if want := generate.Plain([]byte("# mine\n")); !bytes.Equal(rendered, want) {
		t.Errorf(".agents/claude.md = %q, want it rendered from the adopted file %q", rendered, want)
	}
}

func TestForceKeepsOverrides(t *testing.T) {
	upstream, err := fs.ReadFile(stackagents.Assets, "AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	target := append(bytes.Clone(upstream), []byte("\n## Project overrides\n\nTest: custom-cmd\n")...)
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), target, 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := p.Install(stackagents.Assets, true, nil); err != nil {
		t.Fatalf("Install(force): %v", err)
	}

	got, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(got, upstream) {
		t.Errorf("AGENTS.md = %q, want it to start with the new upstream body", got)
	}
	if !bytes.Contains(got, []byte("Test: custom-cmd")) {
		t.Errorf("AGENTS.md = %q, want the project's overrides kept", got)
	}
}

func TestForceWithoutHeadingOverwrites(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# no heading here\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}

	actions, err := p.Install(stackagents.Assets, true, nil)
	if err != nil {
		t.Fatalf("Install(force): %v", err)
	}

	got, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	upstream, err := fs.ReadFile(stackagents.Assets, "AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, upstream) {
		t.Errorf("AGENTS.md = %q, want it replaced by upstream", got)
	}

	warned := false
	for _, line := range actions {
		if strings.Contains(line, "AGENTS.md") && strings.Contains(line, "## Project overrides") {
			warned = true
		}
	}
	if !warned {
		t.Errorf("Install actions = %v, want a warning about the missing %q heading", actions, "## Project overrides")
	}
}
