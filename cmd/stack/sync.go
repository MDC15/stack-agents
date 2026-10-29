package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/MDC15/stack-agents/internal/project"
)

// runSync regenerates .agents/* and refreshes the planned links.
func runSync() error {
	var path, agentCSV string
	var dry, all bool
	fs := flagSet("sync", &path)
	fs.BoolVar(&dry, "dry-run", false, "print the plan without writing")
	fs.StringVar(&agentCSV, "agents", "", "comma-separated agent IDs to link")
	fs.BoolVar(&all, "all", false, "link every agent")
	if err := fs.Parse(os.Args[2:]); err != nil {
		return err
	}

	p, err := openProject(path)
	if err != nil {
		return err
	}

	results, _ := detectAll()
	plan, err := project.LinkPlan(detectedAgents(results), parseAgentIDs(agentCSV), all)
	if err != nil {
		return err
	}

	linked := make([]string, 0, len(plan))
	for id := range plan {
		linked = append(linked, id)
	}
	sort.Strings(linked)
	if dry {
		fmt.Println("dry-run: would sync .agents/* and links for " + strings.Join(linked, ", "))
		return nil
	}
	if err := p.Sync(plan); err != nil {
		return err
	}
	fmt.Println("synced .agents/* and links for " + strings.Join(linked, ", "))
	return nil
}
