package builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

func catalog(t *testing.T) *Catalog {
	t.Helper()
	c, e := Load()
	if e != nil {
		t.Fatal(e)
	}
	return c
}
func config() Config { c := DefaultConfig(); c.StateVersion = "26.05"; return c }
func TestDependencyClosure(t *testing.T) {
	c := catalog(t)
	cfg := config()
	cfg.Bits = []string{"plasma-voxtype"}
	p, e := c.Resolve(cfg)
	if e != nil {
		t.Fatal(e)
	}
	for _, id := range []string{"desktop-plasma", "voxtype", "typing-tools", "user"} {
		if !contains(p.Auto, id) {
			t.Errorf("missing dependency %s", id)
		}
	}
	if len(p.Inputs) != 1 || p.Inputs["voxtype"].URL == "" {
		t.Fatal("wrong input closure", p.Inputs)
	}
}
func TestConflictsAndTracks(t *testing.T) {
	c := catalog(t)
	for _, ids := range [][]string{{"desktop-plasma", "desktop-pantheon"}, {"fish", "zsh"}, {"power-laptop", "power-desktop"}, {"does-not-exist"}, {"desktop-unity"}} {
		cfg := config()
		cfg.Bits = ids
		if _, e := c.Resolve(cfg); e == nil {
			t.Errorf("accepted %v", ids)
		}
	}
}
func TestCycle(t *testing.T) {
	c := catalog(t)
	a := c.byID["user"]
	a.Requires = []string{"zsh"}
	c.byID["user"] = a
	cfg := config()
	cfg.Bits = []string{"zsh"}
	if _, e := c.Resolve(cfg); e == nil || !strings.Contains(e.Error(), "cycle") {
		t.Fatal(e)
	}
}
func TestDeterminismAndRoundTrip(t *testing.T) {
	c := catalog(t)
	cfg := config()
	cfg.Bits = []string{"voxtype", "user", "voxtype"}
	cfg.Hardware = `{ ... }: { fileSystems."/" = { device = "/dev/example"; fsType = "ext4"; }; }`
	a, _, e := c.Render(cfg)
	if e != nil {
		t.Fatal(e)
	}
	cfg.Bits = []string{"user", "voxtype"}
	b, _, e := c.Render(cfg)
	if e != nil {
		t.Fatal(e)
	}
	if a != b {
		t.Fatal("output depends on selection order/duplicates")
	}
	restored, e := ReadConfig([]byte(a))
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(restored, cfg) {
		t.Fatalf("round trip failed: %#v != %#v", restored, cfg)
	}
	if strings.Contains(a, "./hardware-configuration.nix") {
		t.Fatal("external hardware import")
	}
}
func TestStringEscaping(t *testing.T) {
	got := NixString("x\"\\\n${builtins.abort \"injected\"}")
	want := `"x\"\\\n\${builtins.abort \"injected\"}"`
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}
func TestHardware(t *testing.T) {
	for _, src := range []string{`{ imports = [ ./other.nix ]; }`, `{ imports = [ /etc/nixos/hardware.nix ]; }`, `{ imports = [ <nixos-config> ]; }`} {
		if CheckHardware(src) == nil {
			t.Fatal("accepted external hardware", src)
		}
	}
	if e := CheckHardware(`# generated file
{ modulesPath, ... }: { imports = [ (modulesPath + "/installer/scan/not-detected.nix") ]; fileSystems."/".device = "/dev/disk/by-uuid/abc"; }`); e != nil {
		t.Fatal(e)
	}
}
func TestAllBitsAndPresetsParse(t *testing.T) {
	if _, e := exec.LookPath("nix-instantiate"); e != nil {
		t.Skip("Nix parser unavailable")
	}
	c := catalog(t)
	for _, bit := range c.Bits {
		t.Run(bit.ID, func(t *testing.T) {
			cfg := config()
			cfg.Bits = []string{bit.ID}
			if len(bit.Tracks) > 0 {
				cfg.Track = bit.Tracks[0]
			}
			src, _, e := c.Render(cfg)
			if e != nil {
				t.Fatal(e)
			}
			if e = Parse(context.Background(), src); e != nil {
				t.Fatal(e)
			}
		})
	}
	for name, p := range c.Presets {
		t.Run("preset-"+name, func(t *testing.T) {
			cfg := config()
			cfg.Bits = p.Bits
			cfg.Track = p.Track
			src, _, e := c.Render(cfg)
			if e != nil {
				t.Fatal(e)
			}
			if e = Parse(context.Background(), src); e != nil {
				t.Fatal(e)
			}
		})
	}
}
func TestWriteProtectsExisting(t *testing.T) {
	if _, e := exec.LookPath("nix-instantiate"); e != nil {
		t.Skip("Nix parser unavailable")
	}
	c := catalog(t)
	cfg := config()
	src, _, e := c.Render(cfg)
	if e != nil {
		t.Fatal(e)
	}
	dir := t.TempDir()
	target := filepath.Join(dir, "flake.nix")
	if e = os.WriteFile(target, []byte("original"), 0644); e != nil {
		t.Fatal(e)
	}
	if _, e = Write(context.Background(), src, cfg, WriteOptions{Directory: dir}); e == nil {
		t.Fatal("overwrote without force")
	}
	got, _ := os.ReadFile(target)
	if string(got) != "original" {
		t.Fatal("original changed")
	}
	backup, e := Write(context.Background(), src, cfg, WriteOptions{Directory: dir, Force: true})
	if e != nil {
		t.Fatal(e)
	}
	old, e := os.ReadFile(filepath.Join(backup, "flake.nix"))
	if e != nil || string(old) != "original" {
		t.Fatal("missing backup", e)
	}
	got, _ = os.ReadFile(target)
	if string(got) != src {
		t.Fatal("wrong source written")
	}
}
func TestFailedBuildDoesNotPublish(t *testing.T) {
	dir := t.TempDir()
	tools := t.TempDir()
	shell, e := exec.LookPath("sh")
	if e != nil {
		t.Fatal(e)
	}
	commandLog := filepath.Join(t.TempDir(), "nix-command")
	t.Setenv("FLAKEBUILDER_TEST_LOG", commandLog)
	for n, s := range map[string]string{"nix-instantiate": "#!" + shell + "\nexit 0\n", "nix": "#!" + shell + "\nprintf '%s\\n' \"$@\" > \"$FLAKEBUILDER_TEST_LOG\"\nexit 1\n"} {
		if e := os.WriteFile(filepath.Join(tools, n), []byte(s), 0755); e != nil {
			t.Fatal(e)
		}
	}
	t.Setenv("PATH", tools+":"+os.Getenv("PATH"))
	cfg := config()
	cfg.Hardware = "{}"
	target := filepath.Join(dir, "flake.nix")
	os.WriteFile(target, []byte("old"), 0644)
	os.WriteFile(filepath.Join(dir, "flake.lock"), []byte("old-lock"), 0644)
	if _, e := Write(context.Background(), "{}", cfg, WriteOptions{Directory: dir, Force: true, Build: true}); e == nil {
		t.Fatal("failure not propagated")
	}
	commands, e := os.ReadFile(commandLog)
	if e != nil || !strings.Contains(string(commands), "lock") {
		t.Fatalf("Nix locking was not reached: %s %v", commands, e)
	}
	for n, want := range map[string]string{"flake.nix": "old", "flake.lock": "old-lock"} {
		b, _ := os.ReadFile(filepath.Join(dir, n))
		if string(b) != want {
			t.Fatal("changed", n)
		}
	}
}
func TestMalformedCatalog(t *testing.T) {
	for _, src := range []string{`{"version":2}`, `{"version":1,"bits":[{"id":"x","file":"../oops"}]}`, `{"version":1,"inputs":{"self":{"url":"github:x/y"}}}`, `{"version":1,"unknown":true}`} {
		_, e := LoadFS(fstest.MapFS{"catalog.json": &fstest.MapFile{Data: []byte(src)}})
		if e == nil {
			t.Fatal("accepted", src)
		}
	}
}

func TestCustomInputs(t *testing.T) {
	c := catalog(t)
	cfg := config()
	cfg.ExtraInputs = map[string]ExtraInput{"my_packages": {Input: Input{URL: "github:example/packages", Follows: map[string]string{"nixpkgs": "nixpkgs"}}, Uses: []InputUse{{Kind: "overlay", Attribute: "default"}, {Kind: "module", Attribute: "service"}, {Kind: "package", Attribute: "tools.cli"}}}}
	src, p, e := c.Render(cfg)
	if e != nil {
		t.Fatal(e)
	}
	if len(p.Inputs) != 1 {
		t.Fatal("input missing")
	}
	for _, want := range []string{`inputs."nixpkgs".follows = "nixpkgs"`, `inputs."my_packages".overlays."default"`, `inputs."my_packages".nixosModules."service"`, `inputs."my_packages".packages.${pkgs.stdenv.hostPlatform.system}."tools"."cli"`} {
		if !strings.Contains(src, want) {
			t.Error("missing", want)
		}
	}
	if _, e = exec.LookPath("nix-instantiate"); e == nil {
		if e = Parse(context.Background(), src); e != nil {
			t.Fatal(e)
		}
	}
	got, e := ReadConfig([]byte(src))
	if e != nil || !reflect.DeepEqual(got.ExtraInputs, cfg.ExtraInputs) {
		t.Fatal("custom input roundtrip", e)
	}
}
func TestInputValidation(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input ExtraInput
	}{{"self", ExtraInput{Input: Input{URL: "github:a/b"}}}, {"local", ExtraInput{Input: Input{URL: "path:../outside"}}}, {"other", ExtraInput{Input: Input{URL: "github:a/b"}, Uses: []InputUse{{Kind: "overlay", Attribute: `x; abort "bad"`}}}}} {
		if ValidateExtraInput(tc.name, tc.input) == nil {
			t.Fatal("accepted invalid input", tc)
		}
	}
}
func TestInputOverride(t *testing.T) {
	c := catalog(t)
	cfg := config()
	cfg.Bits = []string{"voxtype"}
	cfg.ExtraInputs = map[string]ExtraInput{"voxtype": {Input: Input{URL: "github:peteonrails/voxtype/new-ref", Follows: map[string]string{"nixpkgs": "nixpkgs"}}}}
	p, e := c.Resolve(cfg)
	if e != nil {
		t.Fatal(e)
	}
	if p.Inputs["voxtype"].URL != cfg.ExtraInputs["voxtype"].Input.URL {
		t.Fatal("explicit override ignored")
	}
}

func TestLargeGeneratedFileParses(t *testing.T) {
	if _, e := exec.LookPath("nix-instantiate"); e != nil {
		t.Skip("Nix parser unavailable")
	}
	if e := Parse(context.Background(), "{ description = "+NixString(strings.Repeat("x", 200000))+"; }"); e != nil {
		t.Fatal(e)
	}
}
func TestBuildSequenceAndFailure(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(fmt.Sprint(fail), func(t *testing.T) {
			dir := t.TempDir()
			tools := t.TempDir()
			shell, e := exec.LookPath("sh")
			if e != nil {
				t.Fatal(e)
			}
			commandLog := filepath.Join(t.TempDir(), "commands")
			t.Setenv("FLAKEBUILDER_TEST_LOG", commandLog)
			if fail {
				t.Setenv("FLAKEBUILDER_TEST_FAIL", "yes")
			} else {
				t.Setenv("FLAKEBUILDER_TEST_FAIL", "no")
			}
			script := `printf '%s\n' "$*" >> "$FLAKEBUILDER_TEST_LOG"
if [ "$3" = flake ] && [ "$4" = lock ]; then
 printf '{}\n' > "${5#path:}/flake.lock"
fi
if [ "$3" = build ] && [ "$FLAKEBUILDER_TEST_FAIL" = yes ]; then exit 1; fi
exit 0
`
			for n, s := range map[string]string{"nix-instantiate": "exit 0\n", "nix": script} {
				if e = os.WriteFile(filepath.Join(tools, n), []byte("#!"+shell+"\n"+s), 0755); e != nil {
					t.Fatal(e)
				}
			}
			t.Setenv("PATH", tools+":"+os.Getenv("PATH"))
			cfg := config()
			cfg.Hardware = "{}"
			target := filepath.Join(dir, "flake.nix")
			if e = os.WriteFile(target, []byte("old"), 0644); e != nil {
				t.Fatal(e)
			}
			_, e = Write(context.Background(), "{ }", cfg, WriteOptions{Directory: dir, Force: true, Build: true})
			if (e != nil) != fail {
				t.Fatal("unexpected outcome", e)
			}
			log, e := os.ReadFile(commandLog)
			if e != nil {
				t.Fatal(e)
			}
			lines := strings.Split(strings.TrimSpace(string(log)), "\n")
			if len(lines) != 3 || !strings.Contains(lines[0], "flake lock") || !strings.Contains(lines[1], "flake check --no-build") || !strings.Contains(lines[2], "build --no-link --no-write-lock-file") || !strings.Contains(lines[2], "#nixosConfigurations.nixos.config.system.build.toplevel") {
				t.Fatal(string(log))
			}
			data, e := os.ReadFile(target)
			if e != nil {
				t.Fatal(e)
			}
			want := "{ }"
			if fail {
				want = "old"
			}
			if string(data) != want {
				t.Fatalf("got %s want %s", data, want)
			}
		})
	}
}
