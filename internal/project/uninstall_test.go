package project

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/MDC15/stack-agents"
)

func TestUninstall(t *testing.T) {
	root := t.TempDir()
	p, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.Install(stackagents.Assets, false, map[string]bool{"claude": true, "copilot": true}); err != nil {
		t.Fatalf("Install: %v", err)
	}

	// User-owned files that must survive uninstall: a conventions file,
	// a custom skill, and a foreign .roorules without the banner.
	for _, f := range []string{".ai/conventions.md", ".ai/skills/mine.md", ".roorules"} {
		if err := os.WriteFile(filepath.Join(root, f), []byte("# mine\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	dry, err := p.Uninstall(true)
	if err != nil {
		t.Fatalf("Uninstall(dry): %v", err)
	}
	if len(dry) == 0 {
		t.Error("Uninstall(dry) returned nothing to remove")
	}
	for _, f := range []string{".agents", "CLAUDE.md", "VERSION"} {
		if _, err := os.Lstat(filepath.Join(root, f)); err != nil {
			t.Errorf("Uninstall(dry) touched %s: %v", f, err)
		}
	}

	removed, err := p.Uninstall(false)
	if err != nil {
		t.Fatalf("Uninstall: %v", err)
	}
	if len(removed) == 0 {
		t.Error("Uninstall returned no actions")
	}

	gone := []string{
		".agents", "CLAUDE.md", ".github", ".github/copilot-instructions.md",
		"VERSION", "scripts", "scripts/load_skill.py",
	}
	for _, f := range gone {
		if _, err := os.Lstat(filepath.Join(root, f)); !os.IsNotExist(err) {
			t.Errorf("%s still present after Uninstall (err=%v)", f, err)
		}
	}
	for _, f := range []string{"AGENTS.md", ".ai/conventions.md", ".ai/skills/mine.md", ".roorules"} {
		if _, err := os.Stat(filepath.Join(root, f)); err != nil {
			t.Errorf("%s removed by Uninstall: %v", f, err)
		}
	}

	again, err := p.Uninstall(false)
	if err != nil {
		t.Fatalf("second Uninstall: %v", err)
	}
	if len(again) != 0 {
		t.Errorf("second Uninstall = %v, want nothing left to remove", again)
	}
}

func TestUninstallShippedSetMatchesAssets(t *testing.T) {
	var assetFiles []string
	err := fs.WalkDir(stackagents.Assets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && path != "AGENTS.md" {
			assetFiles = append(assetFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	root := t.TempDir()
	p, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.Install(stackagents.Assets, false, nil); err != nil {
		t.Fatalf("Install: %v", err)
	}
	if _, err := p.Uninstall(false); err != nil {
		t.Fatalf("Uninstall: %v", err)
	}

	// If the embed grows without uninstall's shipped set growing with it,
	// updates would orphan files: every shipped file must be removed.
	for _, rel := range assetFiles {
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(rel))); !os.IsNotExist(err) {
			t.Errorf("shipped file %s survived Uninstall (err=%v)", rel, err)
		}
	}
}
