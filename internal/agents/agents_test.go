package agents

import (
	"strings"
	"testing"
)

func TestAll(t *testing.T) {
	all := All()
	if len(all) != 8 {
		t.Fatalf("len(All()) = %d, want 8", len(all))
	}

	want := []struct {
		id   string
		kind Kind
		file string
		link string
	}{
		{id: "claude", kind: KindSymlink, file: ".agents/claude.md", link: "CLAUDE.md"},
		{id: "opencode", kind: KindDirect},
		{id: "codex", kind: KindDirect},
		{id: "cursor", kind: KindSymlink, file: ".agents/cursor.mdc", link: ".cursor/rules/stack-agents.mdc"},
		{id: "copilot", kind: KindCopy, file: ".agents/copilot.md", link: ".github/copilot-instructions.md"},
		{id: "cline", kind: KindSymlink, file: ".agents/cline.md", link: ".clinerules"},
		{id: "roo", kind: KindSymlink, file: ".agents/roo.md", link: ".roorules"},
		{id: "windsurf", kind: KindSymlink, file: ".agents/windsurf.md", link: ".windsurfrules"},
	}

	for i, w := range want {
		got := all[i]
		if got.ID != w.id || got.Kind != w.kind || got.File != w.file || got.Link != w.link {
			t.Errorf("All()[%d] = {%q %v %q %q}, want {%q %v %q %q}",
				i, got.ID, got.Kind, got.File, got.Link, w.id, w.kind, w.file, w.link)
		}
	}

	for _, a := range all {
		if a.Kind == KindDirect {
			continue
		}
		if !strings.HasPrefix(a.File, ".agents/") {
			t.Errorf("%s: File = %q, want prefix .agents/", a.ID, a.File)
		}
		if a.Link == "" {
			t.Errorf("%s: Link is empty for non-direct agent", a.ID)
		}
	}

	seen := make(map[string]string)
	for _, a := range all {
		if a.Link == "" {
			continue
		}
		if prev, dup := seen[a.Link]; dup {
			t.Errorf("Link %q used by both %s and %s", a.Link, prev, a.ID)
		}
		seen[a.Link] = a.ID
	}

	var copies, directs []string
	for _, a := range all {
		switch a.Kind {
		case KindCopy:
			copies = append(copies, a.ID)
		case KindDirect:
			directs = append(directs, a.ID)
		}
	}
	if len(copies) != 1 || copies[0] != "copilot" {
		t.Errorf("KindCopy agents = %v, want exactly [copilot]", copies)
	}
	if len(directs) != 2 || directs[0] != "opencode" || directs[1] != "codex" {
		t.Errorf("KindDirect agents = %v, want exactly [opencode codex]", directs)
	}
}

func TestByID(t *testing.T) {
	a, ok := ByID("claude")
	if !ok {
		t.Fatal(`ByID("claude") ok = false, want true`)
	}
	if a.ID != "claude" {
		t.Errorf(`ByID("claude").ID = %q, want "claude"`, a.ID)
	}

	if _, ok := ByID("nope"); ok {
		t.Error(`ByID("nope") ok = true, want false`)
	}
}
