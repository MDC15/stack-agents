// Command stack is the stack-agents CLI: a single static binary that ships
// the embedded rules, skills, and scripts.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MDC15/stack-agents/internal/project"
)

// commands maps command names to their handlers.
var commands = map[string]func() error{
	"version":   runVersion,
	"--version": runVersion,
	"init":      runInit,
	"status":    runStatus,
	"sync":      runSync,
	"check":     runCheck,
	"uninstall": runUninstall,
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	name := os.Args[1]
	cmd, ok := commands[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "stack: unknown command %q\n", name)
		usage()
		os.Exit(2)
	}
	if err := cmd(); err != nil {
		fmt.Fprintln(os.Stderr, "stack:", err)
		os.Exit(1)
	}
}

func usage() {
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	fmt.Fprintln(os.Stderr, "usage: stack <command> [args]")
	fmt.Fprintln(os.Stderr, "commands:", strings.Join(names, ", "))
}

// flagSet builds a command's flag set with the shared --path flag.
func flagSet(name string, path *string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.StringVar(path, "path", ".", "project directory")
	return fs
}

// openProject opens path as a stack project, naming AGENTS.md in the
// error when the project has not been initialized.
func openProject(path string) (*project.Project, error) {
	p, err := project.Open(path)
	if err != nil {
		return nil, err
	}
	if _, statErr := os.Stat(filepath.Join(path, "AGENTS.md")); os.IsNotExist(statErr) {
		return nil, fmt.Errorf("AGENTS.md: project is not initialized; run `stack init --path %s`", path)
	}
	return p, nil
}
