package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MDC15/stack-agents/internal/agents"
	"github.com/MDC15/stack-agents/internal/generate"
	"github.com/MDC15/stack-agents/internal/project"
)

// runCheck reports every stale .agents file and link. Drift exits 1.
func runCheck() error {
	var path string
	var diff bool
	fs := flagSet("check", &path)
	fs.BoolVar(&diff, "diff", false, "print unified diffs of stale .agents files")
	if err := fs.Parse(os.Args[2:]); err != nil {
		return err
	}

	p, err := openProject(path)
	if err != nil {
		return err
	}

	_, found := detectAll()
	drifts, err := p.Check(found)
	if err != nil {
		return err
	}
	if len(drifts) == 0 {
		fmt.Println("no drift; everything matches AGENTS.md")
		return nil
	}

	var lines []string
	for _, d := range drifts {
		fmt.Printf("%s: %s\n", d.Problem, d.Path)
		lines = append(lines, d.Problem+": "+d.Path)
	}
	if diff {
		printDiffs(p, drifts)
	}
	return fmt.Errorf("%s", strings.Join(lines, "; "))
}

// printDiffs shows the rendered body versus the file on disk for each
// stale .agents file.
func printDiffs(p *project.Project, drifts []project.Drift) {
	body, err := os.ReadFile(filepath.Join(p.Root, "AGENTS.md"))
	if err != nil {
		return
	}
	for _, d := range drifts {
		if d.Problem != "out of date" || !strings.HasPrefix(d.Path, ".agents/") {
			continue
		}
		var a agents.Agent
		for _, cand := range agents.All() {
			if cand.File == d.Path {
				a = cand
			}
		}
		current, _ := os.ReadFile(filepath.Join(p.Root, d.Path))
		want := renderExpected(a, body)
		if len(current) == 0 {
			continue
		}
		fmt.Println("diff " + d.Path + ":")
		for _, line := range unifiedDiff(current, want) {
			fmt.Println(line)
		}
	}
}

// renderExpected re-renders the content stack would generate for the file.
func renderExpected(a agents.Agent, body []byte) []byte {
	if a.ID == "cursor" {
		return generate.Cursor(body)
	}
	return generate.Plain(body)
}

// unifiedDiff returns a compact unified-style diff: it trims the common
// prefix and suffix, then lists the removed lines (-) and the added
// lines (+) between them.
func unifiedDiff(before, after []byte) []string {
	lines := func(b []byte) []string {
		return strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	}
	a, b := lines(before), lines(after)

	pre := 0
	for pre < len(a) && pre < len(b) && a[pre] == b[pre] {
		pre++
	}
	suf := 0
	for suf < len(a)-pre && suf < len(b)-pre && a[len(a)-1-suf] == b[len(b)-1-suf] {
		suf++
	}

	out := []string{"--- current", "+++ expected"}
	for _, line := range a[pre : len(a)-suf] {
		out = append(out, "-"+line)
	}
	for _, line := range b[pre : len(b)-suf] {
		out = append(out, "+"+line)
	}
	return out
}
