package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	t.Run("clean after sync", func(t *testing.T) {
		root := fixture(t)
		p := synced(t, root, map[string]bool{"claude": true, "copilot": true})

		drifts, err := p.Check(map[string]bool{"claude": true})
		if err != nil {
			t.Fatal(err)
		}
		if len(drifts) != 0 {
			t.Errorf("Check = %+v, want no drift", drifts)
		}
	})

	t.Run("stale agents file", func(t *testing.T) {
		root := fixture(t)
		p := synced(t, root, map[string]bool{"claude": true})
		appendAGENTS(t, root, "# more rules\n")

		drifts, err := p.Check(map[string]bool{"claude": true})
		if err != nil {
			t.Fatal(err)
		}
		wantDrift(t, drifts, ".agents/claude.md", "out of date")
	})

	t.Run("missing link for a detected agent", func(t *testing.T) {
		root := fixture(t)
		p := synced(t, root, map[string]bool{"claude": true})
		if err := os.Remove(filepath.Join(root, "CLAUDE.md")); err != nil {
			t.Fatal(err)
		}

		drifts, err := p.Check(map[string]bool{"claude": true})
		if err != nil {
			t.Fatal(err)
		}
		wantDrift(t, drifts, "CLAUDE.md", "missing")
	})

	t.Run("missing link for a non-detected agent passes", func(t *testing.T) {
		root := fixture(t)
		p := synced(t, root, map[string]bool{"claude": true})
		if err := os.Remove(filepath.Join(root, "CLAUDE.md")); err != nil {
			t.Fatal(err)
		}

		drifts, err := p.Check(nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range drifts {
			if d.Path == "CLAUDE.md" {
				t.Errorf("missing link for non-detected agent reported: %+v", d)
			}
		}
	})

	t.Run("foreign symlink", func(t *testing.T) {
		root := fixture(t)
		p := synced(t, root, map[string]bool{"claude": true})
		if err := os.Remove(filepath.Join(root, "CLAUDE.md")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("elsewhere.md", filepath.Join(root, "CLAUDE.md")); err != nil {
			t.Fatal(err)
		}

		drifts, err := p.Check(map[string]bool{"claude": true})
		if err != nil {
			t.Fatal(err)
		}
		wantDrift(t, drifts, "CLAUDE.md", "wrong target")
	})

	t.Run("edited copy", func(t *testing.T) {
		root := fixture(t)
		p := synced(t, root, map[string]bool{"copilot": true})
		copyPath := filepath.Join(root, ".github", "copilot-instructions.md")
		edited, err := os.ReadFile(copyPath)
		if err != nil {
			t.Fatal(err)
		}
		edited = append(edited, "# tweaked\n"...)
		if err := os.WriteFile(copyPath, edited, 0o644); err != nil {
			t.Fatal(err)
		}

		drifts, err := p.Check(map[string]bool{"copilot": true})
		if err != nil {
			t.Fatal(err)
		}
		wantDrift(t, drifts, ".github/copilot-instructions.md", "out of date")
	})
}

// wantDrift asserts that drifts holds one entry for path whose Problem
// contains problem.
func wantDrift(t *testing.T, drifts []Drift, path, problem string) {
	t.Helper()
	for _, d := range drifts {
		if d.Path == path && strings.Contains(d.Problem, problem) {
			return
		}
	}
	t.Errorf("Check = %+v, want a drift for %s containing %q", drifts, path, problem)
}

// appendAGENTS appends extra to the project's AGENTS.md.
func appendAGENTS(t *testing.T, root, extra string) {
	t.Helper()
	path := filepath.Join(root, "AGENTS.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(content, extra...), 0o644); err != nil {
		t.Fatal(err)
	}
}
