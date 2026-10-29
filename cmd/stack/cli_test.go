package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runScratch returns a fresh empty project directory.
func runScratch(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

// TestStatusFreshDir: status on a non-project must say so by naming
// AGENTS.md and fail with exit 1.
func TestStatusFreshDir(t *testing.T) {
	path := runScratch(t)
	code, out := run(t, "status", "--path", path)
	if code != 1 {
		t.Fatalf("status on fresh dir: exit %d, want 1 (%s)", code, out)
	}
	if !strings.Contains(out, "AGENTS.md") {
		t.Errorf("status output %q does not name AGENTS.md", out)
	}
}

// TestInitAgents: init with explicit agents installs, links exactly the
// two planned agents, and leaves unplanned paths absent.
func TestInitAgents(t *testing.T) {
	path := runScratch(t)
	code, out := run(t, "init", "--path", path, "--yes", "--agents", "claude,cursor")
	if code != 0 {
		t.Fatalf("init --yes --agents claude,cursor: exit %d, want 0 (%s)", code, out)
	}
	if _, err := os.Stat(filepath.Join(path, "AGENTS.md")); err != nil {
		t.Errorf("AGENTS.md not installed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(path, ".agents")); err != nil {
		t.Errorf(".agents/ not generated: %v", err)
	}
	l, err := os.Lstat(filepath.Join(path, "CLAUDE.md"))
	if err != nil || l.Mode()&os.ModeSymlink == 0 {
		t.Errorf("CLAUDE.md should be a symlink: %v (mode %v)", err, l)
	}
	l, err = os.Lstat(filepath.Join(path, ".cursor/rules/stack-agents.mdc"))
	if err != nil || l.Mode()&os.ModeSymlink == 0 {
		t.Errorf(".cursor/rules/stack-agents.mdc should be a symlink: %v (mode %v)", err, l)
	}
	if _, err := os.Lstat(filepath.Join(path, ".windsurfrules")); !os.IsNotExist(err) {
		t.Errorf(".windsurfrules should be absent: %v", err)
	}
}

// TestInitTwiceRefuses: a second non-force init fails naming the install.
func TestInitTwiceRefuses(t *testing.T) {
	path := runScratch(t)
	if code, out := run(t, "init", "--path", path, "--yes", "--agents", "claude"); code != 0 {
		t.Fatalf("first init: exit %d (%s)", code, out)
	}
	code, out := run(t, "init", "--path", path, "--yes")
	if code != 1 {
		t.Fatalf("second init: exit %d, want 1 (%s)", code, out)
	}
	if !strings.Contains(out, "already installed") {
		t.Errorf("second init output %q does not mention already installed", out)
	}
}

// TestSyncAll: sync --all makes every fixed path exist.
func TestSyncAll(t *testing.T) {
	path := runScratch(t)
	if code, out := run(t, "init", "--path", path, "--yes"); code != 0 {
		t.Fatalf("init: exit %d (%s)", code, out)
	}
	code, out := run(t, "sync", "--all", "--path", path)
	if code != 0 {
		t.Fatalf("sync --all: exit %d (%s)", code, out)
	}
	for _, link := range []string{"CLAUDE.md", ".cursor/rules/stack-agents.mdc", ".github/copilot-instructions.md", ".clinerules", ".roorules", ".windsurfrules"} {
		if _, err := os.Lstat(filepath.Join(path, link)); err != nil {
			t.Errorf("%s absent after sync --all: %v", link, err)
		}
	}
}

// TestCheckDrift: check passes on a synced project, fails after AGENTS.md
// is edited, and --diff shows the appended line.
func TestCheckDrift(t *testing.T) {
	path := runScratch(t)
	if code, out := run(t, "init", "--path", path, "--yes", "--agents", "claude"); code != 0 {
		t.Fatalf("init: exit %d (%s)", code, out)
	}
	if code, out := run(t, "check", "--path", path); code != 0 {
		t.Fatalf("check on fresh install: exit %d, want 0 (%s)", code, out)
	}

	agent := "stack: proof of drift"
	f, err := os.OpenFile(filepath.Join(path, "AGENTS.md"), os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("open AGENTS.md: %v", err)
	}
	f.WriteString("\n" + agent + "\n")
	f.Close()

	code, out := run(t, "check", "--path", path)
	if code != 1 {
		t.Fatalf("check after edit: exit %d, want 1 (%s)", code, out)
	}
	if !strings.Contains(out, ".agents/claude.md") {
		t.Errorf("check output %q does not name .agents/claude.md", out)
	}

	code, out = run(t, "check", "--path", path, "--diff")
	if code != 1 {
		t.Fatalf("check --diff: exit %d, want 1 (%s)", code, out)
	}
	if !strings.Contains(out, agent) {
		t.Errorf("check --diff output %q does not contain the appended line", out)
	}
}

// TestStatusListsAllAgents: status reports something for all eight agents,
// including not-found rows when detection has nothing to report.
func TestStatusListsAllAgents(t *testing.T) {
	path := runScratch(t)
	if code, out := run(t, "init", "--path", path, "--yes", "--agents", "claude"); code != 0 {
		t.Fatalf("init: exit %d (%s)", code, out)
	}
	code, out := run(t, "status", "--path", path)
	if code != 0 {
		t.Fatalf("status: exit %d, want 0 (%s)", code, out)
	}
	for _, name := range []string{"Claude Code", "OpenCode", "Codex", "Cursor", "GitHub Copilot", "Cline", "Roo Code", "Windsurf"} {
		if !strings.Contains(out, name) {
			t.Errorf("status output %q missing agent %q", out, name)
		}
	}
	if !strings.Contains(out, "skipped") {
		t.Errorf("status output %q should include not-found (skipped) rows", out)
	}
}

// TestStatusShowsFileDrift: after AGENTS.md is edited, status must stop
// saying ok for the stale agent — status is the health gauge users look
// at, check is what CI runs.
func TestStatusShowsFileDrift(t *testing.T) {
	path := runScratch(t)
	if code, out := run(t, "init", "--path", path, "--yes", "--agents", "claude"); code != 0 {
		t.Fatalf("init: exit %d (%s)", code, out)
	}
	f, err := os.OpenFile(filepath.Join(path, "AGENTS.md"), os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatalf("open AGENTS.md: %v", err)
	}
	if _, err := f.WriteString("\nstack: drift proof\n"); err != nil {
		t.Fatalf("append AGENTS.md: %v", err)
	}
	f.Close()
	code, out := run(t, "status", "--path", path)
	if code != 0 {
		t.Fatalf("status: exit %d (%s)", code, out)
	}
	if !strings.Contains(out, ".agents/claude.md") {
		t.Errorf("status output %q does not name the stale .agents/claude.md", out)
	}
	if strings.Contains(out, "✓ ok\n") && !strings.Contains(out, "drift") {
		t.Errorf("status output %q still reports ok rows while drift exists", out)
	}
}

// TestUninstallRefusesWithoutConfirmation: when the prompt cannot run
// (no terminal attached), uninstall must abort, not assume yes.
func TestUninstallRefusesWithoutConfirmation(t *testing.T) {
	path := runScratch(t)
	if code, out := run(t, "init", "--path", path, "--yes", "--agents", "claude"); code != 0 {
		t.Fatalf("init: exit %d (%s)", code, out)
	}
	code, out := run(t, "uninstall", "--path", path)
	if code != 1 {
		t.Fatalf("uninstall without --yes in a non-tty: exit %d, want 1 (%s)", code, out)
	}
	if _, err := os.Stat(filepath.Join(path, ".agents")); err != nil {
		t.Errorf(".agents should survive a refused uninstall: %v", err)
	}
}

// TestUninstallWithoutAGENTS: a broken install (AGENTS.md deleted by
// hand) can still be uninstalled; the remedy must not be "run init".
func TestUninstallWithoutAGENTS(t *testing.T) {
	path := runScratch(t)
	if code, out := run(t, "init", "--path", path, "--yes", "--agents", "claude"); code != 0 {
		t.Fatalf("init: exit %d (%s)", code, out)
	}
	if err := os.Remove(filepath.Join(path, "AGENTS.md")); err != nil {
		t.Fatalf("removing AGENTS.md: %v", err)
	}
	code, out := run(t, "uninstall", "--path", path, "--yes")
	if code != 0 {
		t.Fatalf("uninstall without AGENTS.md: exit %d, want 0 (%s)", code, out)
	}
	if _, err := os.Stat(filepath.Join(path, ".agents")); !os.IsNotExist(err) {
		t.Errorf(".agents should be gone: %v", err)
	}
}

// TestInitNoDetection: with no --agents and no detected agents, init still
// succeeds but creates no fixed-path files.
func TestInitNoDetection(t *testing.T) {
	path := runScratch(t)
	if code, out := run(t, "init", "--path", path, "--yes"); code != 0 {
		t.Fatalf("init without agents: exit %d (%s)", code, out)
	}
	for _, link := range []string{"CLAUDE.md", ".clinerules", ".roorules", ".windsurfrules", ".cursor", ".github"} {
		if _, err := os.Lstat(filepath.Join(path, link)); !os.IsNotExist(err) {
			t.Errorf("%s should not exist: %v", link, err)
		}
	}
}

// TestUninstallCmd: uninstall removes .agents but keeps AGENTS.md.
func TestUninstallCmd(t *testing.T) {
	path := runScratch(t)
	if code, out := run(t, "init", "--path", path, "--yes", "--agents", "claude"); code != 0 {
		t.Fatalf("init: exit %d (%s)", code, out)
	}
	if code, out := run(t, "uninstall", "--path", path, "--yes"); code != 0 {
		t.Fatalf("uninstall: exit %d (%s)", code, out)
	}
	if _, err := os.Stat(filepath.Join(path, "AGENTS.md")); err != nil {
		t.Errorf("AGENTS.md should survive uninstall: %v", err)
	}
	if _, err := os.Stat(filepath.Join(path, ".agents")); !os.IsNotExist(err) {
		t.Errorf(".agents should be gone: %v", err)
	}
}
