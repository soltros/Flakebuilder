// Package builder resolves a catalog of NixOS bits into one self-contained flake.
package builder

import (
	"bytes"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

//go:embed catalog.json nur_repos.json bits/*.nix
var bundled embed.FS

type Input struct {
	URL     string            `json:"url"`
	Follows map[string]string `json:"follows,omitempty"`
}
type Bit struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	File        string   `json:"file"`
	Requires    []string `json:"requires,omitempty"`
	Conflicts   []string `json:"conflicts,omitempty"`
	Group       string   `json:"group,omitempty"`
	Inputs      []string `json:"inputs,omitempty"`
	Tracks      []string `json:"tracks,omitempty"`
	Systems     []string `json:"systems,omitempty"`
}
type Preset struct {
	Label string   `json:"label"`
	Track string   `json:"track"`
	Bits  []string `json:"bits"`
}
type Catalog struct {
	Version  int               `json:"version"`
	Inputs   map[string]Input  `json:"inputs"`
	Bits     []Bit             `json:"bits"`
	Presets  map[string]Preset `json:"presets"`
	NURRepos []NURRepo         `json:"-"`
	files    fs.FS
	byID     map[string]Bit
}
type NURRepo struct {
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Packages []string `json:"packages"`
}
type Config struct {
	ExtraInputs  map[string]ExtraInput `json:"extraInputs,omitempty"`
	Version      int                   `json:"version"`
	Host         string                `json:"host"`
	System       string                `json:"system"`
	Track        string                `json:"track"`
	StateVersion string                `json:"stateVersion"`
	User         string                `json:"user"`
	Description  string                `json:"description"`
	Timezone     string                `json:"timezone"`
	Locale       string                `json:"locale"`
	Keyboard     string                `json:"keyboard"`
	Workspace    string                `json:"workspace"`
	ConfigDir    string                `json:"configDir"`
	HermesModel  string                `json:"hermesModel"`
	HermesEnv    string                `json:"hermesEnv"`
	Bits         []string              `json:"bits"`
	Hardware     string                `json:"hardware,omitempty"`
	NURRepos     []string              `json:"nurRepos,omitempty"`
}
type Plan struct {
	Bits   []Bit
	Inputs map[string]Input
	Auto   []string
}

var idPattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
var hostPattern = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)
var userPattern = regexp.MustCompile(`^[a-z_][a-z0-9_-]*$`)
var versionPattern = regexp.MustCompile(`^[0-9]{2}\.(05|11)$`)

func Load() (*Catalog, error) { return LoadFS(bundled) }

// LoadFS allows a reviewed local catalog to be used without recompiling the app.
func LoadFS(files fs.FS) (*Catalog, error) {
	b, err := fs.ReadFile(files, "catalog.json")
	if err != nil {
		return nil, err
	}
	var c Catalog
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if err = d.Decode(&c); err != nil {
		return nil, err
	}
	if err = d.Decode(new(any)); err != io.EOF {
		return nil, fmt.Errorf("catalog contains trailing data")
	}
	if c.Version != 1 {
		return nil, fmt.Errorf("unsupported catalog version %d", c.Version)
	}
	if nur, readErr := fs.ReadFile(files, "nur_repos.json"); readErr == nil {
		if err = json.Unmarshal(nur, &c.NURRepos); err != nil {
			return nil, fmt.Errorf("invalid nur_repos.json: %w", err)
		}
	}
	c.files = files
	c.byID = map[string]Bit{}
	for _, b := range c.Bits {
		if !idPattern.MatchString(b.ID) {
			return nil, fmt.Errorf("invalid bit ID %q", b.ID)
		}
		if _, ok := c.byID[b.ID]; ok {
			return nil, fmt.Errorf("duplicate bit %s", b.ID)
		}
		if !fs.ValidPath(b.File) {
			return nil, fmt.Errorf("invalid file path for %s", b.ID)
		}
		if _, err := fs.ReadFile(files, b.File); err != nil {
			return nil, err
		}
		c.byID[b.ID] = b
	}
	for name, in := range c.Inputs {
		if !idPattern.MatchString(name) || name == "nixpkgs" || name == "self" {
			return nil, fmt.Errorf("invalid/reserved input %s", name)
		}
		if in.URL == "" {
			return nil, fmt.Errorf("empty input URL for %s", name)
		}
		for _, target := range in.Follows {
			if target != "nixpkgs" {
				return nil, fmt.Errorf("input %s: only root nixpkgs follows is supported", name)
			}
		}
	}
	for _, b := range c.Bits {
		for _, r := range append(append([]string{}, b.Requires...), b.Conflicts...) {
			if _, ok := c.byID[r]; !ok {
				return nil, fmt.Errorf("%s references unknown bit %s", b.ID, r)
			}
		}
		for _, in := range b.Inputs {
			if _, ok := c.Inputs[in]; !ok {
				return nil, fmt.Errorf("%s references unknown input %s", b.ID, in)
			}
		}
	}
	for name, p := range c.Presets {
		if !idPattern.MatchString(strings.ReplaceAll(name, "_", "-")) {
			return nil, fmt.Errorf("invalid preset %s", name)
		}
		for _, id := range p.Bits {
			if _, ok := c.byID[id]; !ok {
				return nil, fmt.Errorf("preset %s references %s", name, id)
			}
		}
	}
	sort.Slice(c.Bits, func(i, j int) bool {
		a, b := c.Bits[i], c.Bits[j]
		if a.Category != b.Category {
			return a.Category < b.Category
		}
		return a.Label < b.Label
	})
	return &c, nil
}
func DefaultConfig() Config {
	return Config{Version: 1, Host: "nixos", System: "x86_64-linux", Track: "unstable", User: "user", Description: "NixOS user", Timezone: "UTC", Locale: "en_US.UTF-8", Keyboard: "us", Workspace: "/data/workspace", ConfigDir: "/etc/nixos", HermesModel: "gpt-5.5", HermesEnv: "/run/secrets/hermes-env", Bits: []string{"nix-settings", "networkmanager", "localization", "user", "boot-systemd"}}
}
func (c *Catalog) Preset(name string) (Preset, error) {
	p, ok := c.Presets[name]
	if !ok {
		return p, fmt.Errorf("unknown preset %q", name)
	}
	return p, nil
}
func (c *Catalog) Resolve(cfg Config) (Plan, error) {
	p := Plan{Inputs: map[string]Input{}}
	if len(cfg.NURRepos) > 0 {
		known := map[string]bool{}
		for _, repo := range c.NURRepos {
			known[repo.Name] = true
		}
		for _, name := range cfg.NURRepos {
			if !known[name] {
				return p, fmt.Errorf("unknown NUR repository %q", name)
			}
		}
		if nur, ok := c.Inputs["nur"]; ok {
			p.Inputs["nur"] = nur
		}
	}
	state := map[string]int{}
	explicit := map[string]bool{}
	groups := map[string]string{}
	for _, id := range cfg.Bits {
		explicit[id] = true
	}
	var visit func(string) error
	visit = func(id string) error {
		if state[id] == 2 {
			return nil
		}
		if state[id] == 1 {
			return fmt.Errorf("dependency cycle involving %s", id)
		}
		b, ok := c.byID[id]
		if !ok {
			return fmt.Errorf("unknown bit %q", id)
		}
		if len(b.Tracks) > 0 && !contains(b.Tracks, cfg.Track) {
			return fmt.Errorf("%s supports nixpkgs tracks %s, selected %s", id, strings.Join(b.Tracks, ", "), cfg.Track)
		}
		if len(b.Systems) > 0 && !contains(b.Systems, cfg.System) {
			return fmt.Errorf("%s does not support %s", id, cfg.System)
		}
		state[id] = 1
		deps := append([]string{}, b.Requires...)
		sort.Strings(deps)
		for _, dep := range deps {
			if err := visit(dep); err != nil {
				return err
			}
		}
		if b.Group != "" {
			if prev, ok := groups[b.Group]; ok && prev != id {
				return fmt.Errorf("choose one %s: %s conflicts with %s", b.Group, prev, id)
			}
			groups[b.Group] = id
		}
		for _, in := range b.Inputs {
			p.Inputs[in] = c.Inputs[in]
		}
		state[id] = 2
		p.Bits = append(p.Bits, b)
		if !explicit[id] {
			p.Auto = append(p.Auto, id)
		}
		return nil
	}
	ids := append([]string{}, cfg.Bits...)
	sort.Strings(ids)
	for _, id := range ids {
		if err := visit(id); err != nil {
			return p, err
		}
	}
	for _, b := range p.Bits {
		for _, other := range b.Conflicts {
			if state[other] == 2 {
				return p, fmt.Errorf("%s conflicts with %s", b.ID, other)
			}
		}
	}
	for name, input := range cfg.ExtraInputs {
		if err := ValidateExtraInput(name, input); err != nil {
			return p, err
		}
		if original, ok := c.Inputs[name]; ok {
			// An explicit user input may change its URL while retaining required follows.
			for child, target := range original.Follows {
				if value, exists := input.Input.Follows[child]; !exists || value != target {
					return p, fmt.Errorf("input %s requires inputs.%s.follows = %s; preserve that mapping when overriding", name, child, target)
				}
			}
		}
		p.Inputs[name] = input.Input
	}
	sort.Strings(p.Auto)
	return p, nil
}
func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}
func NixString(s string) string {
	// Nix strings differ from JSON/Go escapes; notably ${ opens interpolation.
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "\r", `\r`, "\t", `\t`, "${", `\${`)
	return `"` + r.Replace(s) + `"`
}
func (cfg Config) Validate() error {
	if cfg.Version != 1 {
		return fmt.Errorf("unsupported selection version %d", cfg.Version)
	}
	if !hostPattern.MatchString(cfg.Host) {
		return fmt.Errorf("invalid hostname %q", cfg.Host)
	}
	if cfg.System != "x86_64-linux" && cfg.System != "aarch64-linux" {
		return fmt.Errorf("unsupported system %q", cfg.System)
	}
	if cfg.Track != "unstable" && cfg.Track != "26.05" {
		return fmt.Errorf("nixpkgs track must be unstable or 26.05")
	}
	if !versionPattern.MatchString(cfg.StateVersion) {
		return fmt.Errorf("set --state-version to this installation's original NixOS release (e.g. 26.05)")
	}
	if !userPattern.MatchString(cfg.User) {
		return fmt.Errorf("invalid username %q", cfg.User)
	}
	for _, v := range []string{cfg.Description, cfg.Timezone, cfg.Locale, cfg.Keyboard, cfg.Workspace, cfg.ConfigDir, cfg.HermesModel, cfg.HermesEnv} {
		for _, r := range v {
			if r < 32 || r == 127 {
				return fmt.Errorf("settings must not contain control characters")
			}
		}
	}
	if !strings.HasPrefix(cfg.Workspace, "/") || !strings.HasPrefix(cfg.ConfigDir, "/") || !strings.HasPrefix(cfg.HermesEnv, "/") {
		return fmt.Errorf("workspace, config directory and Hermes environment path must be absolute")
	}
	return nil
}

const metadataPrefix = "# flakebuilder-selection-v1: "

// ReadConfig recovers the exact choices, including embedded hardware, from a generated flake.
func ReadConfig(source []byte) (Config, error) {
	var cfg Config
	for _, line := range strings.Split(string(source), "\n") {
		if strings.HasPrefix(line, metadataPrefix) {
			b, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(line, metadataPrefix))
			if err != nil {
				return cfg, err
			}
			err = json.Unmarshal(b, &cfg)
			if err != nil {
				return cfg, err
			}
			return cfg, cfg.Validate()
		}
	}
	return cfg, fmt.Errorf("no Flakebuilder selection metadata found")
}
func (c *Catalog) Render(cfg Config) (string, Plan, error) {
	if err := cfg.Validate(); err != nil {
		return "", Plan{}, err
	}
	plan, err := c.Resolve(cfg)
	if err != nil {
		return "", plan, err
	}
	cfg.Bits = append([]string{}, cfg.Bits...)
	sort.Strings(cfg.Bits)
	cfg.Bits = unique(cfg.Bits)
	meta, err := json.Marshal(cfg)
	if err != nil {
		return "", plan, err
	}
	var out strings.Builder
	out.WriteString("# Generated by Flakebuilder. All selected configuration bits are embedded below.\n")
	out.WriteString(metadataPrefix + base64.StdEncoding.EncodeToString(meta) + "\n")
	out.WriteString("{\n  description = " + NixString("Flakebuilder configuration for "+cfg.Host) + ";\n  inputs = {\n")
	track := cfg.Track
	if track == "unstable" {
		track = "nixos-unstable"
	} else {
		track = "nixos-" + track
	}
	out.WriteString("    nixpkgs.url = " + NixString("github:NixOS/nixpkgs/"+track) + ";\n")
	names := make([]string, 0, len(plan.Inputs))
	for n := range plan.Inputs {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		in := plan.Inputs[n]
		fmt.Fprintf(&out, "    %s = {\n      url = %s;\n", NixString(n), NixString(in.URL))
		follows := make([]string, 0, len(in.Follows))
		for k := range in.Follows {
			follows = append(follows, k)
		}
		sort.Strings(follows)
		for _, k := range follows {
			fmt.Fprintf(&out, "      inputs.%s.follows = %s;\n", NixString(k), NixString(in.Follows[k]))
		}
		out.WriteString("    };\n")
	}
	out.WriteString("  };\n\n  outputs = inputs@{ nixpkgs, ... }: {\n    nixosConfigurations.")
	out.WriteString(NixString(cfg.Host) + " = nixpkgs.lib.nixosSystem {\n      system = " + NixString(cfg.System) + ";\n      specialArgs = { inherit inputs; };\n      modules = [\n")
	fmt.Fprintf(&out, "        ({ ... }: { networking.hostName = %s; system.stateVersion = %s; })\n", NixString(cfg.Host), NixString(cfg.StateVersion))
	if cfg.Hardware != "" {
		out.WriteString("\n        # Embedded hardware configuration\n        (\n")
		out.WriteString(indent(cfg.Hardware, 10))
		out.WriteString("\n        )\n")
	}
	for _, bit := range plan.Bits {
		src, err := fs.ReadFile(c.files, bit.File)
		if err != nil {
			return "", plan, err
		}
		t, err := template.New(bit.ID).Delims("[[", "]]").Option("missingkey=error").Funcs(template.FuncMap{"nix": NixString, "home": func() string { return "/home/" + cfg.User }}).Parse(string(src))
		if err != nil {
			return "", plan, err
		}
		var body bytes.Buffer
		if err = t.Execute(&body, cfg); err != nil {
			return "", plan, err
		}
		fmt.Fprintf(&out, "\n        # Bit: %s — %s\n        (\n%s\n        )\n", bit.ID, bit.Label, indent(strings.TrimSpace(body.String()), 10))
	}
	for _, name := range names {
		for _, use := range cfg.ExtraInputs[name].Uses {
			expr := InputExpression(name, use)
			fmt.Fprintf(&out, "\n        # User input: %s (%s)\n", name, use.Kind)
			switch use.Kind {
			case "overlay":
				fmt.Fprintf(&out, "        ({ ... }: { nixpkgs.overlays = [ %s ]; })\n", expr)
			case "module":
				fmt.Fprintf(&out, "        %s\n", expr)
			case "package":
				fmt.Fprintf(&out, "        ({ pkgs, ... }: { environment.systemPackages = [ %s ]; })\n", expr)
			}
		}
	}
	if len(cfg.NURRepos) > 0 {
		out.WriteString("\n        # Selected NUR repositories are exposed through the NUR overlay.\n        ({ ... }: { nixpkgs.overlays = [ inputs.nur.overlay ]; })\n")
	}
	out.WriteString("      ];\n    };\n  };\n}\n")
	return out.String(), plan, nil
}
func indent(s string, n int) string {
	pad := strings.Repeat(" ", n)
	return pad + strings.ReplaceAll(strings.TrimSpace(s), "\n", "\n"+pad)
}
func unique(xs []string) []string {
	out := []string{}
	for _, x := range xs {
		if len(out) == 0 || out[len(out)-1] != x {
			out = append(out, x)
		}
	}
	return out
}
