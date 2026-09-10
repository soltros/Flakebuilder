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
	inputForm      *inputForm
	Catalog        *builder.Catalog
	Config         builder.Config
	selected       map[string]bool
	cursor         int
	categoryCursor int
	activeCategory string
	nurMode        bool
	nurCursor      int
	filter         string
	searching      bool
	preview        string
	offset         int
	height         int
	message        string
	Confirmed      bool
	Err            error
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
		if m.activeCategory != "" && b.Category != m.activeCategory {
			continue
		}
		if strings.Contains(strings.ToLower(b.Category+" "+b.Label+" "+b.ID), strings.ToLower(m.filter)) {
			out = append(out, b)
		}
	}
	return out
}
func (m Model) categories() []string {
	set := map[string]bool{}
	for _, b := range m.Catalog.Bits {
		set[b.Category] = true
	}
	out := make([]string, 0, len(set))
	for c := range set {
		out = append(out, c)
	}
	sort.Strings(out)
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
func containsString(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}
func toggleString(xs []string, value string) []string {
	out := make([]string, 0, len(xs)+1)
	found := false
	for _, x := range xs {
		if x == value {
			found = true
			continue
		}
		out = append(out, x)
	}
	if !found {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
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
		if m.nurMode {
			repos := m.Catalog.NURRepos
			switch key {
			case "esc", "backspace":
				m.nurMode = false
			case "up", "k":
				if m.nurCursor > 0 {
					m.nurCursor--
				}
			case "down", "j":
				if m.nurCursor+1 < len(repos) {
					m.nurCursor++
				}
			case " ":
				if len(repos) > 0 {
					m.Config.NURRepos = toggleString(m.Config.NURRepos, repos[m.nurCursor].Name)
				}
			}
			return m, nil
		}
		if m.preview != "" {
			switch key {
			case "y", "g", "enter":
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
		categoryMode := m.activeCategory == "" && m.filter == ""
		switch key {
		case "q":
			return m, tea.Quit
		case "i":
			m.inputForm = &inputForm{Values: [6]string{"", "", "input", "default", "no", ""}}
		case "n":
			m.nurMode = true
			m.nurCursor = 0
		case "/":
			m.searching = true
		case "esc":
			if m.filter != "" {
				m.filter = ""
				m.cursor = 0
			} else if m.activeCategory != "" {
				m.activeCategory = ""
				m.cursor = 0
			} else {
				m.cursor = 0
			}
		case "backspace":
			if m.activeCategory != "" && m.filter == "" {
				m.activeCategory = ""
				m.cursor = 0
			}
		case "up", "k":
			if categoryMode {
				if m.categoryCursor > 0 {
					m.categoryCursor--
				}
			} else if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if categoryMode {
				if m.categoryCursor+1 < len(m.categories()) {
					m.categoryCursor++
				}
			} else if m.cursor+1 < len(choices) {
				m.cursor++
			}
		case " ":
			if !categoryMode && len(choices) > 0 {
				b := choices[m.cursor]
				m.selected[b.ID] = !m.selected[b.ID]
			}
			m.message = ""
		case "enter":
			if categoryMode {
				cats := m.categories()
				if len(cats) > 0 {
					m.activeCategory = cats[m.categoryCursor]
					m.cursor = 0
				}
				return m, nil
			}
			source, _, err := m.Catalog.Render(m.Selection())
			if err != nil {
				m.message = err.Error()
			} else {
				m.preview = source
				m.offset = 0
			}
		case "g":
			m.Config = m.Selection()
			m.Confirmed = true
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m Model) View() string {
	if m.inputForm != nil {
		return m.inputView()
	}
	if m.Confirmed {
		return "Selection confirmed. Saving generated flake…\n"
	}
	header := fmt.Sprintf("Flakebuilder — %s · %s · nixpkgs %s\n", m.Config.Host, m.Config.System, m.Config.Track)
	if m.nurMode {
		var b strings.Builder
		b.WriteString(header + "\nNUR repositories — Space subscribes/unsubscribes\n\n")
		limit := max(1, m.height-8)
		start := max(0, m.nurCursor-limit+1)
		end := min(len(m.Catalog.NURRepos), start+limit)
		for i := start; i < end; i++ {
			r := m.Catalog.NURRepos[i]
			mark := " "
			if containsString(m.Config.NURRepos, r.Name) {
				mark = "x"
			}
			cur := " "
			if i == m.nurCursor {
				cur = ">"
			}
			fmt.Fprintf(&b, "%s [%s] %-22s %d packages\n", cur, mark, r.Name, len(r.Packages))
		}
		if len(m.Catalog.NURRepos) > 0 {
			r := m.Catalog.NURRepos[m.nurCursor]
			shown := r.Packages
			if len(shown) > 18 {
				shown = shown[:18]
			}
			fmt.Fprintf(&b, "\n%s: %s\n", r.Name, strings.Join(shown, ", "))
		}
		b.WriteString("\nEsc/Backspace: return · q: cancel\n")
		return b.String()
	}
	if m.preview != "" {
		lines := strings.Split(m.preview, "\n")
		// The metadata line is intentionally long; show its purpose in the preview.
		if len(lines) > 1 {
			lines[1] = "# Saved selection metadata (embedded in the file)"
		}
		end := min(len(lines), m.offset+max(1, m.height-6))
		return header + "\n" + strings.Join(lines[m.offset:end], "\n") + "\n\n↑/↓ PgUp/PgDn: scroll · Enter/g/y: generate this flake · n: edit · Ctrl+C: cancel\n"
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
	if m.activeCategory == "" && m.filter == "" && !m.searching {
		cats := m.categories()
		fmt.Fprintf(&b, "\nCategories (%d) · selected bits: %d\n", len(cats), len(m.Selection().Bits))
		for i, category := range cats {
			count, selected := 0, 0
			for _, bit := range m.Catalog.Bits {
				if bit.Category == category {
					count++
					if m.selected[bit.ID] {
						selected++
					}
				}
			}
			cursor := " "
			if i == m.categoryCursor {
				cursor = ">"
			}
			fmt.Fprintf(&b, "%s %-24s %2d bits", cursor, category, count)
			if selected > 0 {
				fmt.Fprintf(&b, " · %d selected", selected)
			}
			b.WriteString("\n")
		}
		b.WriteString("\nEnter: open category · [g] Generate flake · /: search all bits · i: add input · q: cancel\n")
		return b.String()
	}
	if m.activeCategory != "" {
		fmt.Fprintf(&b, "Category: %s · Backspace/Esc: categories\n", m.activeCategory)
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
	b.WriteString("\nSpace: select · /: search · i: add input · n: NUR repositories · Enter: preview · g: Generate flake · Backspace: categories · q: cancel\n[x] selected · [+] required automatically (remove its dependents to omit it)\n")
	return b.String()
}
