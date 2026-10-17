package project

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MDC15/stack-agents/internal/agents"
	"github.com/MDC15/stack-agents/internal/generate"
)

const rules = "# rules\n"

// fixture writes a temp project whose AGENTS.md holds the rules body and
// returns the project root.
func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte(rules), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// synced opens the project at root and syncs it with plan.
func synced(t *testing.T, root string, plan map[string]bool) *Project {
	t.Helper()
	p, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Sync(plan); err != nil {
		t.Fatalf("Sync(%v): %v", plan, err)
	}
	return p
}

func TestOpenRejectsNonDirectory(t *testing.T) {
	root := fixture(t)
	file := filepath.Join(root, "AGENTS.md")
	if _, err := Open(file); err == nil {
		t.Errorf("Open(%s) on a regular file = nil error, want an error", file)
	}
}

func TestLinkPlan(t *testing.T) {
	claude, _ := agents.ByID("claude")

	t.Run("detected union extraIDs", func(t *testing.T) {
		plan, err := LinkPlan([]agents.Agent{claude}, []string{"cursor"}, false)
		if err != nil {
			t.Fatal(err)
		}
		if !plan["claude"] || !plan["cursor"] || len(plan) != 2 {
			t.Errorf("LinkPlan(detected=[claude], extra=[cursor]) = %v, want {claude, cursor}", plan)
		}
	})

	t.Run("all links every non-direct agent", func(t *testing.T) {
		plan, err := LinkPlan(nil, nil, true)
		if err != nil {
			t.Fatal(err)
		}
		for _, a := range agents.All() {
			if a.Kind == agents.KindDirect {
				if plan[a.ID] {
					t.Errorf("LinkPlan(all) planned direct agent %q", a.ID)
				}
				continue
			}
			if !plan[a.ID] {
				t.Errorf("LinkPlan(all) missed non-direct agent %q", a.ID)
			}
		}
	})

	t.Run("unknown extraID is an error", func(t *testing.T) {
		if _, err := LinkPlan(nil, []string{"nope"}, false); err == nil {
			t.Error("LinkPlan(extra=[nope]) = nil error, want an error")
		}
	})
}

func TestSyncCreatesAllFilesAndPlannedLinks(t *testing.T) {
	root := fixture(t)
	synced(t, root, map[string]bool{"claude": true, "cursor": true, "copilot": true})

	for _, f := range []string{
		".agents/claude.md", ".agents/cursor.mdc", ".agents/copilot.md",
		".agents/cline.md", ".agents/roo.md", ".agents/windsurf.md",
	} {
		if _, err := os.Stat(filepath.Join(root, f)); err != nil {
			t.Errorf("%s missing after Sync: %v", f, err)
		}
	}

	got, err := os.Readlink(filepath.Join(root, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("CLAUDE.md is not a symlink: %v", err)
	}
	if want := ".agents/claude.md"; got != want {
		t.Errorf("CLAUDE.md target = %q, want %q", got, want)
	}

	// The known-risky nested case: the target is relative to the link's
	// own directory, two levels down.
	got, err = os.Readlink(filepath.Join(root, ".cursor", "rules", "stack-agents.mdc"))
	if err != nil {
		t.Fatalf("stack-agents.mdc is not a symlink: %v", err)
	}
	if want := "../../.agents/cursor.mdc"; got != want {
		t.Errorf(".cursor/rules/stack-agents.mdc target = %q, want %q", got, want)
	}

	copilotPath := filepath.Join(root, ".github", "copilot-instructions.md")
	fi, err := os.Lstat(copilotPath)
	if err != nil {
		t.Fatalf("copilot fixed path missing: %v", err)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		t.Fatal(".github/copilot-instructions.md is a symlink, want a regular copy")
	}
	content, err := os.ReadFile(copilotPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(content, []byte(generate.Banner)) {
		t.Errorf(".github/copilot-instructions.md = %q, want the banner", content)
	}

	for _, l := range []string{".clinerules", ".roorules", ".windsurfrules"} {
		if _, err := os.Lstat(filepath.Join(root, l)); !os.IsNotExist(err) {
			t.Errorf("%s exists, want unplanned and absent", l)
		}
	}
}

func TestSyncReplacesBannerCarryingLeftovers(t *testing.T) {
	root := fixture(t)
	leftover := generate.Plain([]byte("# old layout rules\n"))
	if err := os.WriteFile(filepath.Join(root, "CLAUDE.md"), leftover, 0o644); err != nil {
		t.Fatal(err)
	}

	synced(t, root, map[string]bool{"claude": true})

	got, err := os.Readlink(filepath.Join(root, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("CLAUDE.md is not a symlink after Sync: %v", err)
	}
	if want := ".agents/claude.md"; got != want {
		t.Errorf("CLAUDE.md target = %q, want %q", got, want)
	}
}

func TestSyncFallbackCopyWhenSymlinkFails(t *testing.T) {
	root := fixture(t)
	original := Symlink
	Symlink = func(oldname, newname string) error { return errors.New("symlinks disabled") }
	t.Cleanup(func() { Symlink = original })

	p := synced(t, root, map[string]bool{"cline": true})

	path := filepath.Join(root, ".clinerules")
	fi, err := os.Lstat(path)
	if err != nil {
		t.Fatalf(".clinerules missing after Sync: %v", err)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		t.Fatal(".clinerules is a symlink, want the fallback copy")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if want := generate.Plain([]byte(rules)); !bytes.Equal(content, want) {
		t.Errorf(".clinerules = %q, want %q", content, want)
	}

	drifts, err := p.Check(map[string]bool{"cline": true})
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range drifts {
		if d.Path == ".clinerules" {
			t.Errorf("Check reports drift for the accepted fallback copy: %+v", d)
		}
	}
}

func TestSyncPresenceOptIn(t *testing.T) {
	root := fixture(t)
	// A correct link left by an earlier run, kept even though nothing is
	// planned (presence = opt-in).
	if err := os.Symlink(".agents/cline.md", filepath.Join(root, ".clinerules")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# rules v2\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	synced(t, root, map[string]bool{})

	got, err := os.Readlink(filepath.Join(root, ".clinerules"))
	if err != nil {
		t.Fatalf(".clinerules no longer a symlink: %v", err)
	}
	if want := ".agents/cline.md"; got != want {
		t.Errorf(".clinerules target = %q, want %q", got, want)
	}
	content, err := os.ReadFile(filepath.Join(root, ".clinerules"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(content, []byte("# rules v2\n")) {
		t.Errorf(".clinerules content = %q, want refreshed rules", content)
	}
}

func TestSyncRefreshesUnplannedCopy(t *testing.T) {
	root := fixture(t)
	synced(t, root, map[string]bool{"copilot": true})

	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# rules v2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	synced(t, root, map[string]bool{})

	content, err := os.ReadFile(filepath.Join(root, ".github", "copilot-instructions.md"))
	if err != nil {
		t.Fatalf("copilot copy removed: %v", err)
	}
	if !bytes.Contains(content, []byte("# rules v2\n")) {
		t.Errorf("copilot copy = %q, want refreshed rules", content)
	}
}

func TestSyncLeavesForeignFileAlone(t *testing.T) {
	root := fixture(t)
	foreign := []byte("# my own rules\n")
	if err := os.WriteFile(filepath.Join(root, ".clinerules"), foreign, 0o644); err != nil {
		t.Fatal(err)
	}

	p, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	err = p.Sync(map[string]bool{"cline": true})
	if err == nil {
		t.Fatal("Sync with a foreign .clinerules = nil error, want a report")
	}
	if !strings.Contains(err.Error(), ".clinerules") {
		t.Errorf("Sync error %q does not name .clinerules", err)
	}

	got, readErr := os.ReadFile(filepath.Join(root, ".clinerules"))
	if readErr != nil {
		t.Fatalf(".clinerules was removed: %v", readErr)
	}
	if !bytes.Equal(got, foreign) {
		t.Errorf(".clinerules = %q, want the untouched %q", got, foreign)
	}
}
