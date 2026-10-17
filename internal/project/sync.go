package project

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MDC15/stack-agents/internal/agents"
	"github.com/MDC15/stack-agents/internal/generate"
)

// LinkPlan returns the set of agent IDs to link: detected ∪ extraIDs, or
// every non-direct agent when all is set. An unknown extraID is an error.
func LinkPlan(detected []agents.Agent, extraIDs []string, all bool) (map[string]bool, error) {
	extras := make(map[string]agents.Agent, len(extraIDs))
	for _, id := range extraIDs {
		a, ok := agents.ByID(id)
		if !ok {
			return nil, fmt.Errorf("unknown agent %q", id)
		}
		extras[id] = a
	}

	plan := make(map[string]bool)
	if all {
		for _, a := range agents.All() {
			if a.Kind != agents.KindDirect {
				plan[a.ID] = true
			}
		}
	}
	for _, a := range detected {
		if a.Kind != agents.KindDirect {
			plan[a.ID] = true
		}
	}
	for _, a := range extras {
		if a.Kind != agents.KindDirect {
			plan[a.ID] = true
		}
	}
	return plan, nil
}

// Sync writes every .agents/* file from the project's AGENTS.md, then
// creates or refreshes the fixed-path entry for each planned non-direct
// agent: a relative symlink (a banner-carrying copy when Symlink fails),
// replacing whatever sits there when it is stack-owned. Unplanned agents
// keep their existing entries — presence is opt-in — and get fresh
// content: symlinks through the always-overwritten .agents/* files,
// copies by rewrite. Unmanaged entries are left alone and reported as
// the returned error.
func (p *Project) Sync(plan map[string]bool) error {
	body, err := os.ReadFile(filepath.Join(p.Root, "AGENTS.md"))
	if err != nil {
		return fmt.Errorf("read AGENTS.md: %w", err)
	}

	for _, a := range agents.All() {
		if a.File == "" {
			continue
		}
		if err := p.write(filepath.Join(p.Root, a.File), render(a, body)); err != nil {
			return err
		}
	}

	var skipped []string
	for _, a := range agents.All() {
		if a.Kind == agents.KindDirect {
			continue
		}
		if !plan[a.ID] {
			// Presence = opt-in: keep the entry. Symlinks refresh
			// through the .agents/* overwrite above; a plain copy
			// must be rewritten by hand.
			if a.Kind == agents.KindCopy {
				if err := p.refreshCopy(a, body); err != nil {
					return err
				}
			}
			continue
		}
		report, err := p.place(a, body)
		if err != nil {
			return err
		}
		if report != "" {
			skipped = append(skipped, report)
		}
	}
	if len(skipped) > 0 {
		return errors.New(strings.Join(skipped, "; "))
	}
	return nil
}

// entryState classifies what sits at a fixed path.
type entryState int

const (
	entryAbsent entryState = iota
	entryOwned
	entryUnmanaged
)

// place creates or refreshes the fixed-path entry for a, returning a
// report when an unmanaged entry was left untouched.
func (p *Project) place(a agents.Agent, body []byte) (string, error) {
	linkPath := filepath.Join(p.Root, a.Link)
	state, err := p.owned(linkPath)
	if err != nil {
		return "", err
	}
	if state == entryUnmanaged {
		return fmt.Sprintf("%s: exists but is not managed by stack; left alone", a.Link), nil
	}
	if state == entryOwned {
		if err := os.Remove(linkPath); err != nil {
			return "", err
		}
	}

	content := render(a, body)
	if a.Kind == agents.KindCopy {
		return "", p.write(linkPath, content)
	}

	// Symlink target relative to the link's own directory, so the layout
	// survives the repo being moved or cloned anywhere.
	rel, err := filepath.Rel(filepath.Dir(linkPath), filepath.Join(p.Root, a.File))
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(linkPath), 0o755); err != nil {
		return "", err
	}
	if err := Symlink(rel, linkPath); err != nil {
		// Symlinks fail on e.g. Windows checkouts without developer
		// mode and on FAT mounts; a banner-carrying copy keeps the
		// agent working and Check satisfied.
		return "", p.write(linkPath, content)
	}
	return "", nil
}

// refreshCopy rewrites the fixed-path copy of an unplanned KindCopy
// agent when the copy is stack-owned, keeping presence = opt-in links
// current without planning them. Absent or unmanaged paths are untouched.
func (p *Project) refreshCopy(a agents.Agent, body []byte) error {
	path := filepath.Join(p.Root, a.Link)
	state, err := p.owned(path)
	if err != nil || state != entryOwned {
		return err
	}
	return p.write(path, render(a, body))
}

// owned classifies what sits at path: absent, stack-owned (a symlink
// resolving into .agents/ or a banner-carrying file), or unmanaged.
func (p *Project) owned(path string) (entryState, error) {
	fi, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return entryAbsent, nil
	}
	if err != nil {
		return entryAbsent, err
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(path)
		if err != nil {
			return entryAbsent, err
		}
		resolved := filepath.Clean(filepath.Join(filepath.Dir(path), target))
		agentsDir := filepath.Join(p.Root, ".agents")
		if resolved == agentsDir || strings.HasPrefix(resolved, agentsDir+string(filepath.Separator)) {
			return entryOwned, nil
		}
		return entryUnmanaged, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return entryAbsent, err
	}
	if hasBanner(content) {
		return entryOwned, nil
	}
	return entryUnmanaged, nil
}

// write writes content to path, creating parent directories.
func (p *Project) write(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// render returns the generated content for a from the AGENTS.md body.
func render(a agents.Agent, body []byte) []byte {
	if a.ID == "cursor" {
		return generate.Cursor(body)
	}
	return generate.Plain(body)
}

// hasBanner reports whether content carries the generated banner within
// its first 8 lines (the Cursor frontmatter pushes it down a few lines).
func hasBanner(content []byte) bool {
	for i, line := range bytes.Split(content, []byte("\n")) {
		if i >= 8 {
			break
		}
		if bytes.Contains(line, []byte(generate.Banner)) {
			return true
		}
	}
	return false
}
