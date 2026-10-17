package project

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// overridesHeading marks the section of AGENTS.md the project owns; a
// force install keeps everything from it onwards.
const overridesHeading = "## Project overrides"

// installRoot lists the root entries an install ships, in check order.
var installRoot = []string{"AGENTS.md", "VERSION", ".ai", "scripts"}

// Install ships the asset set into the project root and syncs the
// planned agents. It returns the actions taken, one line each.
//
// A root that already holds any shipped entry is refused as already
// installed unless force is set — except when it holds an AGENTS.md but
// no VERSION: that is a project adopting stack, whose AGENTS.md is kept
// as-is and becomes the source of truth for Sync. With force, an
// existing AGENTS.md is merged instead: everything before its
// "## Project overrides" heading is replaced by the upstream body and
// the project's tail from that heading on is kept; a missing heading
// means a full overwrite and a returned warning line.
func (p *Project) Install(assets fs.FS, force bool, plan map[string]bool) ([]string, error) {
	var present []string
	presence := make(map[string]bool, len(installRoot))
	for _, name := range installRoot {
		if _, err := os.Lstat(filepath.Join(p.Root, name)); err == nil {
			presence[name] = true
			present = append(present, name)
		} else if !os.IsNotExist(err) {
			return nil, err
		}
	}
	adopt := presence["AGENTS.md"] && !presence["VERSION"]
	if !force && !adopt && len(present) > 0 {
		return nil, fmt.Errorf("already installed: %s (use force to overwrite)", strings.Join(present, ", "))
	}

	upstream, err := fs.ReadFile(assets, "AGENTS.md")
	if err != nil {
		return nil, fmt.Errorf("read upstream AGENTS.md: %w", err)
	}

	var actions []string
	switch {
	case !presence["AGENTS.md"]:
		if err := p.write(filepath.Join(p.Root, "AGENTS.md"), upstream); err != nil {
			return nil, err
		}
		actions = append(actions, "AGENTS.md")
	case adopt && !force:
		// The project's own AGENTS.md becomes the source of truth.
		actions = append(actions, "AGENTS.md (kept the existing file)")
	default: // force with an existing AGENTS.md: merge, keep overrides
		warning, err := p.mergeAGENTS(upstream)
		if err != nil {
			return nil, err
		}
		actions = append(actions, "AGENTS.md")
		if warning != "" {
			actions = append(actions, warning)
		}
	}

	// Ship the rest of the set: overwriting shipped names, keeping any
	// extra files the project already has in place.
	err = fs.WalkDir(assets, ".", func(rel string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || rel == "AGENTS.md" {
			return nil
		}
		content, err := fs.ReadFile(assets, rel)
		if err != nil {
			return err
		}
		if err := p.write(filepath.Join(p.Root, filepath.FromSlash(rel)), content); err != nil {
			return err
		}
		actions = append(actions, rel)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("ship assets: %w", err)
	}

	if err := p.Sync(plan); err != nil {
		return actions, err
	}
	return actions, nil
}

// mergeAGENTS replaces the project's AGENTS.md up to (not including)
// its overrides heading with upstream, keeping the project's overrides.
// It returns a warning when the heading is missing and the file was
// replaced entirely.
func (p *Project) mergeAGENTS(upstream []byte) (string, error) {
	path := filepath.Join(p.Root, "AGENTS.md")
	target, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	offset, found := overridesOffset(target)
	if !found {
		if err := p.write(path, upstream); err != nil {
			return "", err
		}
		return fmt.Sprintf("AGENTS.md: no %q heading, replaced with upstream", overridesHeading), nil
	}
	merged := bytes.Clone(upstream)
	if len(merged) > 0 && merged[len(merged)-1] != '\n' {
		merged = append(merged, '\n')
	}
	merged = append(merged, target[offset:]...)
	return "", p.write(path, merged)
}

// overridesOffset returns the byte offset of the first line in content
// equal to the overrides heading after trimming surrounding space.
func overridesOffset(content []byte) (int, bool) {
	offset := 0
	for _, line := range bytes.Split(content, []byte("\n")) {
		if strings.TrimSpace(string(line)) == overridesHeading {
			return offset, true
		}
		offset += len(line) + 1
	}
	return 0, false
}
