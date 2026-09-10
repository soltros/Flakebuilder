// Flakebuilder is a standalone descendant of soltros/configbuilder's menu-based
// NixOS configuration workflow. It emits a single Nix source file.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/soltros/Flakebuilder/internal/builder"
	"github.com/soltros/Flakebuilder/internal/tui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Flakebuilder:", err)
		os.Exit(1)
	}
}
func run() error {
	cfg := builder.DefaultConfig()
	var dir, hardware, preset, selected, from, catalogDir string
	var yes, stdout, list, force, lock, build bool
	var inputs, uses, follows, remove repeatedFlag
	flag.Var(&inputs, "input", "Add/override NAME=URL (repeatable)")
	flag.Var(&uses, "input-use", "Use NAME=overlay:ATTRIBUTE, NAME=module:ATTRIBUTE, or NAME=package:ATTRIBUTE (repeatable)")
	flag.Var(&follows, "input-follows", "Make NAME.inputs.nixpkgs follow root nixpkgs (repeatable)")
	flag.Var(&remove, "remove-input", "Remove a custom input loaded with --from (repeatable)")
	flag.StringVar(&dir, "dir", defaultOutputDirectory, "Output directory (the default keeps backups when generating again)")
	flag.StringVar(&from, "from", "", "Load choices and embedded hardware from a generated flake")
	flag.StringVar(&catalogDir, "catalog", "", "Use a trusted local directory containing catalog.json and bit templates")
	flag.StringVar(&preset, "preset", "", "Start with a preset (see --list)")
	flag.StringVar(&selected, "bits", "", "Explicit comma-separated bit IDs (replaces preset selections)")
	flag.StringVar(&hardware, "hardware", "", "Generated hardware-configuration.nix to embed, not import")
	flag.StringVar(&cfg.Host, "host", cfg.Host, "Hostname and flake output key")
	flag.StringVar(&cfg.System, "system", cfg.System, "Target platform")
	flag.StringVar(&cfg.Track, "nixpkgs", cfg.Track, "nixpkgs track: unstable or 26.05")
	flag.StringVar(&cfg.StateVersion, "state-version", "", "Required: this installation's original stateVersion")
	flag.StringVar(&cfg.User, "user", cfg.User, "Normal user name")
	flag.StringVar(&cfg.Description, "user-description", cfg.Description, "User description")
	flag.StringVar(&cfg.Timezone, "timezone", cfg.Timezone, "Timezone")
	flag.StringVar(&cfg.Locale, "locale", cfg.Locale, "Locale")
	flag.StringVar(&cfg.Keyboard, "keyboard", cfg.Keyboard, "Keyboard layout")
	flag.StringVar(&cfg.Workspace, "workspace", cfg.Workspace, "Workspace directory for optional personal integrations")
	flag.StringVar(&cfg.ConfigDir, "config-dir", cfg.ConfigDir, "Installed configuration directory used by rebuild aliases")
	flag.StringVar(&cfg.HermesModel, "hermes-model", cfg.HermesModel, "Hermes default model")
	flag.StringVar(&cfg.HermesEnv, "hermes-env", cfg.HermesEnv, "Runtime Hermes environment file; contents are never embedded")
	flag.BoolVar(&list, "list", false, "List catalog bits and presets")
	flag.BoolVar(&yes, "yes", false, "Generate noninteractively from the supplied selections")
	flag.BoolVar(&stdout, "stdout", false, "Print generated flake without writing files or running Nix")
	flag.BoolVar(&force, "force", false, "Allow replacing generated files, with backups")
	flag.BoolVar(&lock, "lock", false, "Save the flake, then resolve inputs and evaluate; emits flake.lock on success")
	flag.BoolVar(&build, "build", false, "Save the flake, then lock, evaluate and build; never activate")
	flag.Parse()
	if flag.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(flag.Args(), " "))
	}
	var cat *builder.Catalog
	var err error
	if catalogDir == "" {
		cat, err = builder.Load()
	} else {
		cat, err = builder.LoadFS(os.DirFS(catalogDir))
	}
	if err != nil {
		return err
	}
	if list {
		fmt.Println("Bits ([+] dependencies are selected automatically):")
		for _, b := range cat.Bits {
			fmt.Printf("  %-24s %-22s %s\n", b.ID, b.Category, b.Label)
		}
		fmt.Println("\nPresets:")
		for _, name := range sortedPresets(cat) {
			p := cat.Presets[name]
			fmt.Printf("  %-24s nixpkgs %-8s %s\n", name, p.Track, p.Label)
		}
		return nil
	}
	var automaticBackup bool
	dir, automaticBackup, err = outputDirectory(dir)
	if err != nil {
		return err
	}
	force = force || automaticBackup
	seen := map[string]bool{}
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	overrides := cfg
	if from != "" {
		data, e := os.ReadFile(from)
		if e != nil {
			return e
		}
		cfg, e = builder.ReadConfig(data)
		if e != nil {
			return e
		}
	}
	// Preset changes selections and track, never hardware or user identity.
	if preset != "" {
		p, e := cat.Preset(preset)
		if e != nil {
			return e
		}
		cfg.Bits = p.Bits
		cfg.Track = p.Track
	}
	values := map[string][2]*string{
		"host": {&cfg.Host, &overrides.Host}, "system": {&cfg.System, &overrides.System}, "nixpkgs": {&cfg.Track, &overrides.Track}, "state-version": {&cfg.StateVersion, &overrides.StateVersion},
		"user": {&cfg.User, &overrides.User}, "user-description": {&cfg.Description, &overrides.Description}, "timezone": {&cfg.Timezone, &overrides.Timezone}, "locale": {&cfg.Locale, &overrides.Locale},
		"keyboard": {&cfg.Keyboard, &overrides.Keyboard}, "workspace": {&cfg.Workspace, &overrides.Workspace}, "config-dir": {&cfg.ConfigDir, &overrides.ConfigDir}, "hermes-model": {&cfg.HermesModel, &overrides.HermesModel}, "hermes-env": {&cfg.HermesEnv, &overrides.HermesEnv},
	}
	for name, pair := range values {
		if seen[name] {
			*pair[0] = *pair[1]
		}
	}
	if seen["bits"] {
		cfg.Bits = []string{}
		for _, id := range strings.Split(selected, ",") {
			if id = strings.TrimSpace(id); id != "" {
				cfg.Bits = append(cfg.Bits, id)
			}
		}
	}
	if hardware != "" {
		b, e := os.ReadFile(hardware)
		if e != nil {
			return e
		}
		cfg.Hardware = string(b)
	}
	if err = applyInputs(&cfg, cat, inputs, uses, follows, remove); err != nil {
		return err
	}
	if err = cfg.Validate(); err != nil {
		return err
	}
	if stdout && (build || lock) {
		return fmt.Errorf("--stdout cannot be combined with --lock or --build")
	}
	if !yes && !stdout {
		fmt.Fprintf(os.Stderr, "Output: %s/flake.nix (build: %t, replace existing: %t)\n", dir, build, force)
		result, e := tea.NewProgram(tui.New(cat, cfg), tea.WithAltScreen()).Run()
		if e != nil {
			return e
		}
		m := result.(tui.Model)
		if !m.Confirmed {
			fmt.Fprintln(os.Stderr, "Cancelled; no output files changed.")
			return nil
		}
		cfg = m.Config
	}
	source, plan, err := cat.Render(cfg)
	if err != nil {
		return err
	}
	if stdout {
		fmt.Print(source)
		return nil
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	backup, err := builder.Write(ctx, source, cfg, builder.WriteOptions{Directory: dir, Force: force, Lock: lock, Build: build, DryRun: true, Log: os.Stderr})
	var validation *builder.ValidationError
	if err != nil && !errors.As(err, &validation) {
		return err
	}
	fmt.Printf("Generated %s/flake.nix: %d bits, %d external inputs.\n", dir, len(plan.Bits), len(plan.Inputs))
	if backup != "" {
		fmt.Println("Previous files backed up to", backup)
	}
	if err != nil {
		return err
	}
	if build {
		fmt.Println("System build passed. No activation was performed.")
	} else if lock {
		fmt.Println("Inputs locked and flake evaluation passed.")
	} else {
		fmt.Println("Nix dry-run passed. No activation was performed; use --lock or --build for deeper checks.")
	}
	if cfg.Hardware == "" {
		fmt.Println("Hardware was omitted; embed it with --hardware before building for a real machine.")
	}
	return nil
}

func sortedPresets(c *builder.Catalog) []string {
	out := make([]string, 0, len(c.Presets))
	for name := range c.Presets {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
