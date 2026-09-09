// Package tui implements ConfigBuilder-style selection with dependency feedback.
package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/soltros/Flakebuilder/internal/builder"
)

type Model struct {
	inputForm *inputForm
	Catalog   *builder.Catalog
	Config    builder.Config
	selected  map[string]bool
	cursor    int
	filter    string
	searching bool
	preview   string
	offset    int
	height    int
	message   string
	Confirmed bool
	Err       error
}

func New(c *builder.Catalog, cfg builder.Config) Model {
	m := Model{Catalog: c, Config: cfg, selected: map[string]bool{}, height: 24}
	for _, id := range cfg.Bits {
		m.selected[id] = true
	}
	return m
}
func (m Model) Init() tea.Cmd { return nil }
func (m Model) choices() []builder.Bit {
	var out []builder.Bit
	for _, b := range m.Catalog.Bits {
		if strings.Contains(strings.ToLower(b.Category+" "+b.Label+" "+b.ID), strings.ToLower(m.filter)) {
			out = append(out, b)
		}
	}
	return out
}
func (m Model) Selection() builder.Config {
	cfg := m.Config
	cfg.Bits = []string{}
	for id, on := range m.selected {
		if on {
			cfg.Bits = append(cfg.Bits, id)
		}
	}
	sort.Strings(cfg.Bits)
	return cfg
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
	case tea.KeyMsg:
		key := msg.String()
		if key == "ctrl+c" {
			return m, tea.Quit
		}
		if m.inputForm != nil {
			return m.updateInput(msg)
		}
		if m.preview != "" {
			switch key {
			case "y":
				m.Config = m.Selection()
				m.Confirmed = true
				return m, tea.Quit
			case "esc", "n":
				m.preview = ""
				m.offset = 0
			case "up", "k":
				if m.offset > 0 {
					m.offset--
				}
			case "down", "j":
				if m.offset < len(strings.Split(m.preview, "\n"))-1 {
					m.offset++
				}
			case "pgdown":
				m.offset = min(m.offset+max(1, m.height-7), len(strings.Split(m.preview, "\n"))-1)
			case "pgup":
				m.offset = max(0, m.offset-max(1, m.height-7))
			}
			return m, nil
		}
		if m.searching {
			switch key {
			case "enter", "esc":
				m.searching = false
			case "backspace":
				r := []rune(m.filter)
				if len(r) > 0 {
					m.filter = string(r[:len(r)-1])
				}
			default:
				if msg.Type == tea.KeyRunes {
					m.filter += string(msg.Runes)
				}
			}
			m.cursor = 0
			return m, nil
		}
		choices := m.choices()
		switch key {
		case "q":
			return m, tea.Quit
		case "i":
			m.inputForm = &inputForm{Values: [6]string{"", "", "input", "default", "no", ""}}
		case "/":
			m.searching = true
		case "esc":
			m.filter = ""
			m.cursor = 0
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor+1 < len(choices) {
				m.cursor++
			}
		case " ":
			if len(choices) > 0 {
				b := choices[m.cursor]
				m.selected[b.ID] = !m.selected[b.ID]
			}
			m.message = ""
		case "enter", "g":
			source, _, err := m.Catalog.Render(m.Selection())
			if err != nil {
				m.message = err.Error()
			} else {
				m.preview = source
				m.offset = 0
			}
		}
	}
	return m, nil
}
func (m Model) View() string {
	if m.inputForm != nil {
		return m.inputView()
	}
	if m.Confirmed {
		return "Selection confirmed. Validating generated flake…\n"
	}
	header := fmt.Sprintf("Flakebuilder — %s · %s · nixpkgs %s\n", m.Config.Host, m.Config.System, m.Config.Track)
	if m.preview != "" {
		lines := strings.Split(m.preview, "\n")
		// The metadata line is intentionally long; show its purpose in the preview.
		if len(lines) > 1 {
			lines[1] = "# Saved selection metadata (embedded in the file)"
		}
		end := min(len(lines), m.offset+max(1, m.height-6))
		return header + "\n" + strings.Join(lines[m.offset:end], "\n") + "\n\n↑/↓ PgUp/PgDn: scroll · y: generate this flake · n: edit · Ctrl+C: cancel\n"
	}
	plan, err := m.Catalog.Resolve(m.Selection())
	auto := map[string]bool{}
	for _, id := range plan.Auto {
		auto[id] = true
	}
	var b strings.Builder
	b.WriteString(header)
	names := make([]string, 0, len(m.Config.ExtraInputs))
	for name := range m.Config.ExtraInputs {
		names = append(names, name)
	}
	sort.Strings(names)
	if len(names) > 0 {
		fmt.Fprintf(&b, "Custom inputs: %s\n", strings.Join(names, ", "))
	}
	fmt.Fprintf(&b, "%d explicit choices · %d dependencies · %d external inputs\n", len(m.Selection().Bits), len(plan.Auto), len(plan.Inputs))
	if m.filter != "" || m.searching {
		fmt.Fprintf(&b, "Search: %s\n", m.filter)
	}
	choices := m.choices()
	count := max(1, m.height-10)
	start := max(0, m.cursor-count+1)
	end := min(len(choices), start+count)
	for i := start; i < end; i++ {
		bit := choices[i]
		mark := " "
		if m.selected[bit.ID] {
			mark = "x"
		} else if auto[bit.ID] {
			mark = "+"
		}
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		fmt.Fprintf(&b, "%s [%s] %-22s %s\n", cursor, mark, bit.Category, bit.Label)
	}
	if len(choices) == 0 {
		b.WriteString("No matching bits. Esc clears the filter.\n")
	} else {
		fmt.Fprintf(&b, "\n%s\n", choices[min(m.cursor, len(choices)-1)].Description)
	}
	if err != nil {
		b.WriteString("\n" + err.Error() + "\n")
	} else if m.message != "" {
		b.WriteString("\n" + m.message + "\n")
	}
	b.WriteString("\nSpace: select · /: search · i: add input · Enter: preview · q: cancel\n[x] selected · [+] required automatically (remove its dependents to omit it)\n")
	return b.String()
}
