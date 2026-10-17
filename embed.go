// Package stackagents embeds the files that an install of the stack CLI ships.
package stackagents

import "embed"

// Assets holds the ship set: agent rules, the version file, skill and
// playbook definitions, and the skill loader script.
//
//go:embed AGENTS.md VERSION .ai scripts/load_skill.py
var Assets embed.FS
