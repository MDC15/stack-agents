package main

import (
	"fmt"
	"io/fs"
	"strings"

	stackagents "github.com/MDC15/stack-agents"
)

// runVersion prints the embedded version of the stack CLI.
func runVersion() error {
	data, err := fs.ReadFile(stackagents.Assets, "VERSION")
	if err != nil {
		return fmt.Errorf("reading VERSION: %w", err)
	}
	fmt.Println("stack " + strings.TrimSpace(string(data)))
	return nil
}
