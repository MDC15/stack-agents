// Package detect finds which supported AI coding agents are installed on
// this machine. Detection is best-effort: every probe degrades to
// not-found rather than failing the run.
package detect

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/MDC15/stack-agents/internal/agents"
)

// Finder abstracts the machine probes detection needs so tests inject fakes.
type Finder interface {
	// InPath reports whether bin resolves on PATH.
	InPath(bin string) bool
	// HomeExists reports whether rel exists under the user's home directory.
	HomeExists(rel string) bool
	// VSCodeExts lists installed VS Code extensions; (nil, false) when the
	// code binary is absent or the listing fails.
	VSCodeExts() ([]string, bool)
	// Version returns the trimmed first line of <bin> --version, or "" on
	// failure or absence.
	Version(bin string) string
}

// Result is the detection outcome for one agent.
type Result struct {
	Agent   agents.Agent
	Found   bool
	Version string // filled only when Found
}

// All probes each registered agent through f. Version is filled only for
// found agents; nothing here panics or fails hard.
func All(f Finder) []Result {
	exts, haveExts := f.VSCodeExts()

	all := agents.All()
	results := make([]Result, 0, len(all))
	for _, a := range all {
		r := Result{Agent: a, Found: found(f, a.ID, exts, haveExts)}
		if r.Found {
			r.Version = f.Version(a.ID)
		}
		results = append(results, r)
	}
	return results
}

// found applies the per-agent detection signals from the design spec.
// An unknown ID is not found.
func found(f Finder, id string, exts []string, haveExts bool) bool {
	switch id {
	case "claude":
		return f.InPath("claude") || f.HomeExists(".claude")
	case "opencode":
		return f.InPath("opencode") || f.HomeExists(".config/opencode")
	case "codex":
		return f.InPath("codex") || f.HomeExists(".codex")
	case "cursor":
		return f.InPath("cursor") || f.HomeExists(".cursor")
	case "copilot":
		return haveExts && hasExt(exts, "github.copilot")
	case "cline":
		return haveExts && hasExt(exts, "saoudrizwan.claude-dev")
	case "roo":
		return haveExts && hasExt(exts, "rooveterinaryinc.roo-cline")
	case "windsurf":
		return f.InPath("windsurf") || f.HomeExists(".codeium")
	}
	return false
}

// hasExt reports whether exts contains want, compared case-insensitively.
func hasExt(exts []string, want string) bool {
	want = strings.ToLower(want)
	for _, e := range exts {
		if strings.ToLower(e) == want {
			return true
		}
	}
	return false
}

// machineFinder is the real Finder: PATH lookups, home-directory stats and
// short-lived shelled-out commands.
type machineFinder struct {
	once sync.Once
	exts []string
	ok   bool
}

// NewFinder returns a Finder that probes this machine.
func NewFinder() Finder { return &machineFinder{} }

func (m *machineFinder) InPath(bin string) bool {
	_, err := exec.LookPath(bin)
	return err == nil
}

func (m *machineFinder) HomeExists(rel string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(home, rel))
	return err == nil
}

// VSCodeExts runs code --list-extensions once per Finder and caches the
// outcome, including the negative one.
func (m *machineFinder) VSCodeExts() ([]string, bool) {
	m.once.Do(m.loadExts)
	return m.exts, m.ok
}

func (m *machineFinder) loadExts() {
	if _, err := exec.LookPath("code"); err != nil {
		return // code absent: (nil, false)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "code", "--list-extensions").Output()
	if err != nil {
		return // listing failed: (nil, false)
	}
	for _, line := range strings.Split(string(out), "\n") {
		if s := strings.TrimSpace(line); s != "" {
			m.exts = append(m.exts, s)
		}
	}
	m.ok = true
}

// Version shells out <bin> --version with a 2s timeout and returns the
// trimmed first line, or "" on any failure including a missing binary.
func (m *machineFinder) Version(bin string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, bin, "--version").Output()
	if err != nil {
		return ""
	}
	first, _, _ := strings.Cut(string(out), "\n")
	return strings.TrimSpace(first)
}
