// Package generate renders .agents/* file content from a project's AGENTS.md.
package generate

// Banner marks every generated file so sync and uninstall can recognize
// stack-owned output and refuse to touch user-owned files.
const Banner = "<!-- Generated from AGENTS.md by stack sync. Edit AGENTS.md, not this file. -->"

// cursorFrontmatter is the YAML frontmatter Cursor's rule files require.
const cursorFrontmatter = "---\ndescription: stack-agents operating rules\nalwaysApply: true\n---\n\n"

// Plain returns body preceded by the banner and a blank line.
func Plain(body []byte) []byte {
	out := make([]byte, 0, len(Banner)+2+len(body))
	out = append(out, Banner...)
	out = append(out, '\n', '\n')
	out = append(out, body...)
	return out
}

// Cursor returns body preceded by the Cursor frontmatter and the Plain header.
func Cursor(body []byte) []byte {
	plain := Plain(body)
	out := make([]byte, 0, len(cursorFrontmatter)+len(plain))
	out = append(out, cursorFrontmatter...)
	out = append(out, plain...)
	return out
}
