package ui

import (
	"strings"
	"testing"
)

// TestRenderRows pins the three visible facts the brief asks for: every
// agent name appears, present/absent marks appear, and an empty version
// renders as an em dash.
func TestRenderRows(t *testing.T) {
	rows := []Row{
		{Agent: "Claude Code", Version: "1.2.3", Status: "✓ link", Target: "CLAUDE.md"},
		{Agent: "OpenCode", Version: "", Status: "✗ skip", Target: "AGENTS.md"},
	}
	out := RenderRows(rows)
	for _, want := range []string{"Claude Code", "OpenCode", "1.2.3", "✓", "✗", "—"} {
		if !strings.Contains(out, want) {
			t.Errorf("RenderRows output %q missing %q", out, want)
		}
	}
}
