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
	// Open a category, then Enter previews the generated flake.
	m = key(m, "enter")
	m = key(m, "enter")
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

func TestCategoryNavigationKeepsSelectionFocused(t *testing.T) {
	m := model(t)
	if len(m.categories()) < 2 {
		t.Fatal("catalog should expose multiple categories")
	}
	// Enter opens the first category instead of placing the cursor in the full catalog.
	m = key(m, "enter")
	if m.activeCategory == "" {
		t.Fatal("category menu did not open a category")
	}
	if len(m.choices()) == 0 {
		t.Fatal("opened category has no bits")
	}
	first := m.choices()[0].ID
	m = key(m, " ")
	if !m.selected[first] {
		t.Fatal("space did not select a bit inside the category")
	}
	m = key(m, "backspace")
	if m.activeCategory != "" {
		t.Fatal("backspace did not return to category menu")
	}
}

func TestGenerateFromPreview(t *testing.T) {
	for _, confirm := range []string{"enter", "y"} {
		t.Run(confirm, func(t *testing.T) {
			m := key(model(t), "enter")
			m = key(m, "enter")
			if m.preview == "" {
				t.Fatal("preview unavailable")
			}
			if confirm == "enter" {
				v, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
				m = v.(Model)
			} else {
				m = key(m, confirm)
			}
			if !m.Confirmed {
				t.Fatal("generate key did not confirm output")
			}
		})
	}
}

func TestGenerateKeyConfirmsDirectly(t *testing.T) {
	m := key(model(t), "g")
	if !m.Confirmed || m.preview != "" {
		t.Fatal("g should generate directly from the selection menu")
	}
}
