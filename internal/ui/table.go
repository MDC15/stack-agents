// Package ui renders the tables and prompts the stack CLI shows its users.
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// Row is one line of an agent table. An empty Target stays blank; an
// empty Version renders as an em dash.
type Row struct {
	Agent   string
	Version string
	Status  string
	Target  string
}

var (
	// presentStyle marks a row the command will act on.
	presentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	// absentStyle marks a row the command will skip.
	absentStyle = lipgloss.NewStyle().Faint(true)
)

// RenderRows renders rows as a plain Agent (ACTION) / Version / Status /
// Target table. Statuses starting with ✓ render green, ✗ renders dim.
func RenderRows(rows []Row) string {
	data := make([][]string, 0, len(rows)+1)
	data = append(data, []string{"AGENT", "VERSION", "STATUS", "TARGET"})
	for _, r := range rows {
		status := r.Status
		switch {
		case strings.HasPrefix(status, "✓"):
			status = presentStyle.Render(status)
		case strings.HasPrefix(status, "✗"):
			status = absentStyle.Render(status)
		}
		version := r.Version
		if version == "" {
			version = "—"
		}
		data = append(data, []string{r.Agent, version, status, r.Target})
	}
	return table.New().
		Headers(data[0]...).
		Rows(data[1:]...).
		String()
}
