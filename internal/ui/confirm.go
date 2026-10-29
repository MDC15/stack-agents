package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbletea"
)

// confirmModel is a yes/no prompt where every key but n/esc answers yes:
// Enter, y, Y, EOF and a closed input all keep the default.
type confirmModel struct {
	question string
	result   bool
}

func (m confirmModel) Init() tea.Cmd { return nil }

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "y", "Y", "enter":
		m.result = true
		return m, tea.Quit
	case "n", "N", "esc", "ctrl+c", "q":
		m.result = false
		return m, tea.Quit
	}
	return m, nil
}

func (m confirmModel) View() string {
	return m.question + " [Y/n] "
}

// Confirm asks q and reports the answer, defaulting to yes when the user
// presses Enter, EOF being treated as the same input. When the prompt
// cannot run (no terminal attached), it declines to assume an answer and
// reports the problem; callers that need a scripted yes pass explicit
// confirmation instead.
func Confirm(q string, out io.Writer) (bool, error) {
	m := &confirmModel{question: q, result: true}
	p := tea.NewProgram(m, tea.WithOutput(out), tea.WithInput(os.Stdin))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(out, "%s [prompt unavailable: %v; refusing to assume]\n", q, err)
		return false, fmt.Errorf("no confirmation: %w", err)
	}
	return m.result, nil
}
