package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/MDC15/stack-agents/internal/project"
	"github.com/MDC15/stack-agents/internal/ui"
)

// runUninstall removes what stack installed, keeping project-owned files.
func runUninstall() error {
	var path string
	var yes, dry bool
	fs := flagSet("uninstall", &path)
	fs.BoolVar(&yes, "yes", false, "skip the confirmation prompt")
	fs.BoolVar(&dry, "dry-run", false, "print what would be removed")
	if err := fs.Parse(os.Args[2:]); err != nil {
		return err
	}

	p, err := project.Open(path)
	if err != nil {
		return err
	}

	if !yes {
		ok, err := ui.Confirm("remove the stack-managed files from "+path+"?", os.Stdout)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Println("aborted; nothing was removed")
			return nil
		}
	}

	removed, err := p.Uninstall(dry)
	if err != nil {
		return err
	}
	if len(removed) == 0 {
		fmt.Println("nothing to remove")
		return nil
	}
	action := "removed"
	if dry {
		action = "would remove"
	}
	fmt.Printf("%s: %s\n", action, strings.Join(removed, ", "))
	return nil
}
