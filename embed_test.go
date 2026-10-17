package stackagents

import (
	"io/fs"
	"testing"
)

// TestAssetsShipSet asserts that Assets embeds exactly the files an install
// ships and nothing else.
func TestAssetsShipSet(t *testing.T) {
	want := []string{
		"AGENTS.md",
		"VERSION",
		"scripts/load_skill.py",
		".ai/principles.md",
		".ai/skills/architect.md",
		".ai/skills/arena.md",
		".ai/skills/checkpoint.md",
		".ai/skills/conventions.md",
		".ai/skills/red-team.md",
		".ai/skills/tdd.md",
		".ai/skills/unslop.md",
		".ai/skills/verify.md",
		".ai/playbooks/debug.md",
		".ai/playbooks/refactor.md",
	}
	wantSet := make(map[string]bool, len(want))
	for _, p := range want {
		wantSet[p] = true
	}

	got := map[string]bool{}
	err := fs.WalkDir(Assets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			got[path] = true
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking Assets: %v", err)
	}

	for _, p := range want {
		if !got[p] {
			t.Errorf("missing embedded file %q", p)
		}
	}
	for p := range got {
		if !wantSet[p] {
			t.Errorf("unexpected embedded file %q", p)
		}
	}
}
