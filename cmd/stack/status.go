package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MDC15/stack-agents/internal/agents"
	"github.com/MDC15/stack-agents/internal/ui"
)

// runStatus reports per-agent state: ok, drift, missing, or skipped.
func runStatus() error {
	var path string
	fs := flagSet("status", &path)
	if err := fs.Parse(os.Args[2:]); err != nil {
		return err
	}

	p, err := openProject(path)
	if err != nil {
		return err
	}

	results, found := detectAll()
	drifts, err := p.Check(found)
	if err != nil {
		return err
	}
	driftAt := make(map[string]string, len(drifts))
	for _, d := range drifts {
		driftAt[d.Path] = d.Problem
	}
	linkExists := func(rel string) bool {
		_, err := os.Stat(filepath.Join(p.Root, rel))
		return err == nil
	}

	rows := make([]ui.Row, 0, len(results))
	for _, r := range results {
		a := r.Agent
		row := ui.Row{Agent: a.Name, Version: r.Version, Target: a.Link}
		switch {
		case a.Kind == agents.KindDirect:
			row.Status = "✓ ok"
		case driftAt[a.Link] != "":
			row.Status = "✗ " + driftAt[a.Link] + " (" + a.Link + ")"
		case linkExists(a.Link):
			row.Status = "✓ ok"
		case found[a.ID]:
			row.Status = "✗ missing link"
		default:
			row.Status = "✗ skipped (agent not found)"
		}
		rows = append(rows, row)
	}
	fmt.Println(ui.RenderRows(rows))
	return nil
}
