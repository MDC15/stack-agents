package project

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

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
	embedded, err := fs.ReadFile(stackagents.Assets, "AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}
	rules := rulesBefore(embedded)

	// The project's file: rules as an older install shipped them, plus
	// the project's own overrides section.
	root := t.TempDir()
	target := bytes.Clone(rules)
	target = append(target, []byte("## Project overrides\n\nTest: custom-cmd\n")...)
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), target, 0o644); err != nil {
		t.Fatal(err)
	}

	// Upstream has since changed: a new marker line in the rules body,
	// its own overrides section still following it.
	changed := bytes.Clone(rules)
	changed = append(changed, []byte("# upstream revision 2\n\n")...)
	changed = append(changed, embedded[len(rules):]...)

	p, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.Install(changedUpstream(t, changed), true, nil); err != nil {
		t.Fatalf("Install(force): %v", err)
	}

	got, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(got)
	if !strings.Contains(text, "# upstream revision 2") {
		t.Errorf("AGENTS.md = %q, want the new upstream rules body", got)
	}
	if strings.Contains(text, "Repo-specific commands") {
		t.Errorf("AGENTS.md = %q, want upstream's own overrides section truncated", got)
	}
	if !strings.Contains(text, "Test: custom-cmd") {
		t.Errorf("AGENTS.md = %q, want the project's overrides kept", got)
	}

	const heading = "## Project overrides"
	count := 0
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == heading {
			count++
		}
	}
	if count != 1 {
		t.Errorf("AGENTS.md has %d %q lines, want exactly 1:\n%s", count, heading, got)
	}
	if i := strings.Index(text, heading); i < 0 || !strings.Contains(text[i:], "Test: custom-cmd") {
		t.Errorf("AGENTS.md = %q, want %q under the overrides heading", got, "Test: custom-cmd")
	}
}

// rulesBefore returns content up to (not including) its first overrides
// heading line, or all of it when the heading is missing.
func rulesBefore(content []byte) []byte {
	offset := 0
	for _, line := range bytes.Split(content, []byte("\n")) {
		if strings.TrimSpace(string(line)) == "## Project overrides" {
			return content[:offset]
		}
		offset += len(line) + 1
	}
	return content
}

// changedUpstream returns the embedded assets with AGENTS.md replaced
// by the given body, simulating an upstream that changed since install.
func changedUpstream(t *testing.T, agentsMD []byte) fs.FS {
	t.Helper()
	return overlayFS{FS: stackagents.Assets, files: map[string][]byte{"AGENTS.md": agentsMD}}
}

// overlayFS serves the embedded assets, overriding named files.
type overlayFS struct {
	fs.FS
	files map[string][]byte
}

func (o overlayFS) Open(name string) (fs.File, error) {
	if data, ok := o.files[name]; ok {
		return fstest.MapFS{name: &fstest.MapFile{Data: data}}.Open(name)
	}
	return o.FS.Open(name)
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
