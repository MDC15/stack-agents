// Package project models a project root: it reads AGENTS.md, generates
// the .agents/* rule files and keeps every fixed-path link in sync with
// them, and detects drift between the two.
package project

import (
	"fmt"
	"os"
)

// Project is a project root containing AGENTS.md.
type Project struct {
	Root string
}

// Open returns the project rooted at root. It fails when root is not a
// directory.
func Open(root string) (*Project, error) {
	fi, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("open project %s: %w", root, err)
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("open project %s: not a directory", root)
	}
	return &Project{Root: root}, nil
}

// Symlink creates newname as a symlink to oldname. It is a variable so
// tests can force symlink failure and exercise the copy fallback.
var Symlink = os.Symlink
