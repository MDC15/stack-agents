package project

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MDC15/stack-agents/internal/agents"
)

// Drift is a generated file or fixed-path link that no longer matches
// the project's AGENTS.md.
type Drift struct {
	Path    string
	Problem string
}

// Check compares every .agents/* file and fixed-path link against the
// current AGENTS.md. detected holds the agent IDs considered installed:
// their links must exist; missing links for other agents pass.
func (p *Project) Check(detected map[string]bool) ([]Drift, error) {
	body, err := os.ReadFile(filepath.Join(p.Root, "AGENTS.md"))
	if err != nil {
		return nil, fmt.Errorf("read AGENTS.md: %w", err)
	}

	var drifts []Drift
	for _, a := range agents.All() {
		if a.File != "" {
			d, err := p.checkFile(a, body)
			if err != nil {
				return nil, err
			}
			drifts = append(drifts, d...)
		}
		if a.Kind == agents.KindDirect {
			continue
		}
		d, err := p.checkLink(a, body, detected[a.ID])
		if err != nil {
			return nil, err
		}
		drifts = append(drifts, d...)
	}
	return drifts, nil
}

// checkFile reports drift for a's generated .agents/* file.
func (p *Project) checkFile(a agents.Agent, body []byte) ([]Drift, error) {
	got, err := os.ReadFile(filepath.Join(p.Root, a.File))
	if os.IsNotExist(err) {
		return []Drift{{Path: a.File, Problem: "missing"}}, nil
	}
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(got, render(a, body)) {
		return []Drift{{Path: a.File, Problem: "out of date"}}, nil
	}
	return nil, nil
}

// checkLink reports drift for a's fixed-path entry: a correct relative
// symlink or an up-to-date banner copy passes; a missing link passes
// only when the agent is not detected.
func (p *Project) checkLink(a agents.Agent, body []byte, detected bool) ([]Drift, error) {
	path := filepath.Join(p.Root, a.Link)
	fi, err := os.Lstat(path)
	if os.IsNotExist(err) {
		if detected {
			return []Drift{{Path: a.Link, Problem: "missing"}}, nil
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(path)
		if err != nil {
			return nil, err
		}
		resolved := filepath.Clean(filepath.Join(filepath.Dir(path), target))
		if resolved != filepath.Join(p.Root, a.File) {
			return []Drift{{Path: a.Link, Problem: "wrong target"}}, nil
		}
		return nil, nil
	}

	got, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if !hasBanner(got) {
		return []Drift{{Path: a.Link, Problem: "wrong target"}}, nil
	}
	if !bytes.Equal(got, render(a, body)) {
		return []Drift{{Path: a.Link, Problem: "out of date"}}, nil
	}
	return nil, nil
}
