package project

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/MDC15/stack-agents/internal/agents"
)

// shipped lists every relative file path an install puts in the project
// root except AGENTS.md, which uninstall always keeps. It mirrors
// stackagents.Assets — the embed stays test-only so the CLI remains the
// owner of the ship set — and TestUninstallShippedSetMatchesAssets fails
// when the two drift, because drifted updates would orphan files.
var shipped = []string{
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

// Uninstall removes what an install ships and returns the removed paths,
// one line each. It takes only stack-owned things: the whole .agents/
// directory, the fixed paths that are symlinks into .agents/ or
// banner-carrying copies, VERSION, scripts/load_skill.py, and the files
// under .ai/ enumerated in shipped. AGENTS.md, .ai/conventions.md and
// custom skills are user-owned and survive, as do directories that keep
// any user file. With dry, nothing is touched.
func (p *Project) Uninstall(dry bool) ([]string, error) {
	shippedSet := make(map[string]bool, len(shipped))
	for _, name := range shipped {
		shippedSet[name] = true
	}

	var files []string // stack-owned paths to remove, root-relative

	if fi, err := os.Stat(filepath.Join(p.Root, ".agents")); err == nil && fi.IsDir() {
		files = append(files, ".agents/")
	} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	for _, a := range agents.All() {
		if a.Link == "" {
			continue
		}
		state, err := p.owned(filepath.Join(p.Root, a.Link))
		if err != nil {
			return nil, err
		}
		if state == entryOwned {
			files = append(files, a.Link)
		}
	}

	for _, name := range []string{"VERSION", "scripts/load_skill.py"} {
		if _, err := os.Lstat(filepath.Join(p.Root, name)); err == nil {
			files = append(files, name)
		} else if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}

	err := fs.WalkDir(os.DirFS(p.Root), ".ai", func(rel string, d fs.DirEntry, err error) error {
		if errors.Is(err, fs.ErrNotExist) {
			return nil // no .ai/ to clean
		}
		if err != nil {
			return err
		}
		if !d.IsDir() && shippedSet[filepath.ToSlash(rel)] {
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	dirs, err := p.emptyDirs(files)
	if err != nil {
		return nil, err
	}

	if !dry {
		for _, rel := range files {
			path := filepath.Join(p.Root, rel)
			if rel == ".agents/" {
				if err := os.RemoveAll(path); err != nil {
					return nil, err
				}
				continue
			}
			if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Errorf("remove %s: %w", rel, err)
			}
		}
		for _, rel := range dirs {
			if err := os.Remove(filepath.Join(p.Root, rel)); err != nil && !errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Errorf("remove dir %s: %w", rel, err)
			}
		}
	}

	removed := append(append([]string{}, files...), dirs...)
	sort.Strings(removed)
	return removed, nil
}

// emptyDirs returns the directories that hold only removed files and
// other removed directories — parents of the removed files, deepest
// first — so uninstall can clean up after itself without touching user
// content. A directory keeping any other file or subdirectory survives.
func (p *Project) emptyDirs(files []string) ([]string, error) {
	remove := make(map[string]bool, len(files))
	for _, f := range files {
		remove[f] = true
	}

	candidates := map[string]bool{}
	for _, f := range files {
		if f == ".agents/" {
			continue // the whole directory goes as one
		}
		for dir := filepath.Dir(f); dir != "." && dir != ""; dir = filepath.Dir(dir) {
			candidates[dir] = true
		}
	}

	sorted := make([]string, 0, len(candidates))
	for dir := range candidates {
		sorted = append(sorted, dir)
	}
	// A child path is always longer than its parents, so this visits
	// every directory after the directories it contains.
	sort.Slice(sorted, func(i, j int) bool { return len(sorted[i]) > len(sorted[j]) })

	removable := make(map[string]bool, len(sorted))
	var dirs []string
	for _, dir := range sorted {
		entries, err := os.ReadDir(filepath.Join(p.Root, dir))
		if err != nil {
			return nil, err
		}
		empty := true
		for _, e := range entries {
			rel := filepath.ToSlash(filepath.Join(dir, e.Name()))
			if e.IsDir() {
				if !removable[rel] {
					empty = false
					break
				}
				continue
			}
			if !remove[rel] {
				empty = false
				break
			}
		}
		if empty {
			removable[dir] = true
			dirs = append(dirs, dir)
		}
	}
	return dirs, nil
}
