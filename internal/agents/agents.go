// Package agents is the static registry of supported AI coding agents.
// Every fixed path lives here so sync, check and uninstall all agree on
// what stack owns.
package agents

// Kind describes how an agent consumes generated content.
type Kind int

const (
	// KindSymlink agents read the fixed path, which is a relative
	// symlink into .agents/.
	KindSymlink Kind = iota
	// KindCopy agents read the fixed path, which is a banner-carrying
	// copy of the .agents/ file.
	KindCopy
	// KindDirect agents read AGENTS.md itself and need no file or link.
	KindDirect
)

// Agent is one supported AI coding agent.
type Agent struct {
	ID   string // stable identifier, e.g. "claude"
	Name string // display name, e.g. "Claude Code"
	Kind Kind
	File string // generated file under .agents/, e.g. ".agents/claude.md"; "" for KindDirect
	Link string // fixed path the agent reads, relative to project root; "" for KindDirect
}

// all lists every agent in registry order.
var all = []Agent{
	{ID: "claude", Name: "Claude Code", Kind: KindSymlink, File: ".agents/claude.md", Link: "CLAUDE.md"},
	{ID: "opencode", Name: "OpenCode", Kind: KindDirect},
	{ID: "codex", Name: "Codex", Kind: KindDirect},
	{ID: "cursor", Name: "Cursor", Kind: KindSymlink, File: ".agents/cursor.mdc", Link: ".cursor/rules/stack-agents.mdc"},
	{ID: "copilot", Name: "GitHub Copilot", Kind: KindCopy, File: ".agents/copilot.md", Link: ".github/copilot-instructions.md"},
	{ID: "cline", Name: "Cline", Kind: KindSymlink, File: ".agents/cline.md", Link: ".clinerules"},
	{ID: "roo", Name: "Roo Code", Kind: KindSymlink, File: ".agents/roo.md", Link: ".roorules"},
	{ID: "windsurf", Name: "Windsurf", Kind: KindSymlink, File: ".agents/windsurf.md", Link: ".windsurfrules"},
}

// All returns every supported agent in registry order.
func All() []Agent { return all }

// ByID returns the agent with the given ID.
func ByID(id string) (Agent, bool) {
	for _, a := range all {
		if a.ID == id {
			return a, true
		}
	}
	return Agent{}, false
}
