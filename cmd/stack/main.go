// Command stack is the stack-agents CLI: a single static binary that ships
// the embedded rules, skills, and scripts.
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// commands maps command names to their handlers.
var commands = map[string]func() error{
	"version":   runVersion,
	"--version": runVersion,
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
