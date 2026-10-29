package main

import (
	"fmt"
	"os"
	"strings"

	stackagents "github.com/MDC15/stack-agents"
	"github.com/MDC15/stack-agents/internal/agents"
	"github.com/MDC15/stack-agents/internal/detect"
	"github.com/MDC15/stack-agents/internal/project"
	"github.com/MDC15/stack-agents/internal/ui"
)

// parseAgentIDs splits a comma-separated agent list, dropping blanks.
func parseAgentIDs(csv string) []string {
	ids := []string{}
	for _, id := range strings.Split(csv, ",") {
		if id = strings.TrimSpace(id); id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

// detectAll returns the results and their found IDs, computed once.
func detectAll() ([]detect.Result, map[string]bool) {
	results := detect.All(detect.NewFinder())
	ids := make(map[string]bool, len(results))
	for _, r := range results {
		ids[r.Agent.ID] = r.Found
	}
	return results, ids
}

// detectedAgents returns the found Results' agents for LinkPlan.
func detectedAgents(results []detect.Result) []agents.Agent {
	out := []agents.Agent{}
	for _, r := range results {
		if r.Found {
			out = append(out, r.Agent)
		}
	}
	return out
}

// runInit installs stack-agents into a project directory: detect, show
// the plan as a table, confirm, then install and link.
func runInit() error {
	var (
		path, agentCSV       string
		yes, dry, force, all bool
	)
	fs := flagSet("init", &path)
	fs.BoolVar(&yes, "yes", false, "skip the confirmation prompt")
	fs.BoolVar(&dry, "dry-run", false, "print the plan without installing")
	fs.BoolVar(&force, "force", false, "update or overwrite an existing install")
	fs.StringVar(&agentCSV, "agents", "", "comma-separated agent IDs to link")
	fs.BoolVar(&all, "all", false, "link every agent")
	if err := fs.Parse(os.Args[2:]); err != nil {
		return err
	}

	results, _ := detectAll()
	plan, err := project.LinkPlan(detectedAgents(results), parseAgentIDs(agentCSV), all)
	if err != nil {
		return err
	}

	rows := make([]ui.Row, 0, len(results))
	for _, r := range results {
		a := r.Agent
		row := ui.Row{Agent: a.Name, Version: r.Version, Target: a.Link}
		switch {
		case a.Kind == agents.KindDirect:
			row.Status = "✓ reads AGENTS.md"
		case plan[a.ID]:
			row.Status = "✓ link"
		default:
			row.Status = "✗ skip"
		}
		rows = append(rows, row)
	}
	fmt.Println(ui.RenderRows(rows))

	switch {
	case dry:
		fmt.Println("dry-run: no files written")
		return nil
	case !yes:
		ok, err := ui.Confirm("install into "+path+"?", os.Stdout)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("aborted; nothing was written")
			return nil
		}
	}

	p, err := project.Open(path)
	if err != nil {
		return err
	}
	actions, err := p.Install(stackagents.Assets, force, plan)
	if err != nil {
		return err
	}
	for _, action := range actions {
		fmt.Println(action)
	}
	fmt.Println("\nNext: edit AGENTS.md, then run `stack sync --path " + path + "`")
	return nil
}
