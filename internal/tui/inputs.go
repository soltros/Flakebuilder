package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/soltros/Flakebuilder/internal/builder"
	"strings"
)

type inputForm struct {
	Field  int
	Values [6]string
	Error  string
}

var inputLabels = []string{"Name (example: my-packages)", "URL (example: github:owner/repository)", "Use: input, overlay, module, package, or remove", "Output attribute (default, vulkan, etc.; unused for input/remove)", "Follow root nixpkgs? yes/no", "Confirm: press Enter to save, Esc to cancel"}

func (m Model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	f := m.inputForm
	key := msg.String()
	if key == "esc" {
		m.inputForm = nil
		return m, nil
	}
	if key == "enter" {
		if f.Field < 5 {
			f.Field++
			return m, nil
		}
		name, url, kind, attr, follow := f.Values[0], f.Values[1], f.Values[2], f.Values[3], f.Values[4]
		if kind == "remove" {
			delete(m.Config.ExtraInputs, name)
			m.inputForm = nil
			m.message = "Removed custom input " + name
			return m, nil
		}
		extra := builder.ExtraInput{Input: builder.Input{URL: url}}
		if follow != "yes" && follow != "no" {
			f.Error = "Follow must be yes or no"
			return m, nil
		}
		if follow == "yes" {
			extra.Input.Follows = map[string]string{"nixpkgs": "nixpkgs"}
		}
		if kind != "input" {
			extra.Uses = []builder.InputUse{{Kind: kind, Attribute: attr}}
		}
		if err := builder.ValidateExtraInput(name, extra); err != nil {
			f.Error = err.Error()
			return m, nil
		}
		if m.Config.ExtraInputs == nil {
			m.Config.ExtraInputs = map[string]builder.ExtraInput{}
		}
		// Adding another use of the same declaration should preserve previous uses.
		if old, ok := m.Config.ExtraInputs[name]; ok && old.Input.URL == url {
			for _, use := range old.Uses {
				found := false
				for _, newUse := range extra.Uses {
					if use == newUse {
						found = true
					}
				}
				if !found {
					extra.Uses = append(extra.Uses, use)
				}
			}
		}
		m.Config.ExtraInputs[name] = extra
		m.inputForm = nil
		m.message = "Added input " + name
		return m, nil
	}
	switch key {
	case "up", "shift+tab":
		if f.Field > 0 {
			f.Field--
		}
	case "down", "tab":
		if f.Field < 5 {
			f.Field++
		}
	case "ctrl+u":
		if f.Field < 5 {
			f.Values[f.Field] = ""
		}
	case "backspace":
		if f.Field < 5 {
			r := []rune(f.Values[f.Field])
			if len(r) > 0 {
				f.Values[f.Field] = string(r[:len(r)-1])
			}
		}
	default:
		if msg.Type == tea.KeyRunes && f.Field < 5 {
			f.Values[f.Field] += string(msg.Runes)
		}
	}
	f.Error = ""
	return m, nil
}
func (m Model) inputView() string {
	var out strings.Builder
	out.WriteString("Flakebuilder — Add an input\n\n")
	out.WriteString("Overlay → inputs.NAME.overlays.ATTRIBUTE\nModule  → inputs.NAME.nixosModules.ATTRIBUTE\nPackage → inputs.NAME.packages.SYSTEM.ATTRIBUTE\n\n")
	for i, label := range inputLabels {
		cursor := " "
		if i == m.inputForm.Field {
			cursor = ">"
		}
		fmt.Fprintf(&out, "%s %s\n  %s\n", cursor, label, m.inputForm.Values[i])
	}
	if m.inputForm.Error != "" {
		out.WriteString("\n" + m.inputForm.Error + "\n")
	}
	out.WriteString("\nEnter/Tab: next · ↑: previous · Ctrl+U: clear field · Esc: cancel\n")
	out.WriteString("Repeat with the same name/URL to add another output. Use ‘remove’ to delete a custom input.\n")
	return out.String()
}
