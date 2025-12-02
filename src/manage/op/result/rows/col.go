package result

import "github.com/charmbracelet/lipgloss"

type Col interface {
	Style() lipgloss.Style
	Value() string
}
type defaultCol struct {
	value string
}

var _ Col = &defaultCol{}

func (c *defaultCol) Style() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (c *defaultCol) Value() string {
	return c.value + "VAL"
}
