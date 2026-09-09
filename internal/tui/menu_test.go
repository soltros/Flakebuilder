package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/soltros/Flakebuilder/internal/builder"
	"testing"
)

func model(t *testing.T) Model {
	c, e := builder.Load()
	if e != nil {
		t.Fatal(e)
	}
	cfg := builder.DefaultConfig()
	cfg.StateVersion = "26.05"
	return New(c, cfg)
}
func key(m Model, k string) Model {
	v, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)})
	return v.(Model)
}
func TestPreviewFreezesSelection(t *testing.T) {
	m := model(t)
	m = key(m, "g")
	if m.preview == "" {
		t.Fatal("no preview")
	}
	count := len(m.Selection().Bits)
	m = key(m, " ")
	if len(m.Selection().Bits) != count {
		t.Fatal("selection changed in preview")
	}
	m = key(m, "y")
	if !m.Confirmed {
		t.Fatal("confirmation lost")
	}
}
func TestCancelAndEmptySearch(t *testing.T) {
	m := model(t)
	m = key(m, "/")
	m = key(m, "not-a-real-bit")
	m.searching = false
	m = key(m, " ")
	_ = m.View()
	m = key(m, "q")
	if m.Confirmed {
		t.Fatal("quit confirmed output")
	}
}
