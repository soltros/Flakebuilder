package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/soltros/Flakebuilder/internal/builder"
)

func TestOutputDirectory(t *testing.T) {
	home := t.TempDir()
	for _, value := range []string{defaultOutputDirectory, filepath.Join(home, "generated_flakes"), filepath.Join(home, "generated_flakes") + "/"} {
		dir, backup, err := resolveOutputDirectory(value, home)
		if err != nil || dir != filepath.Join(home, "generated_flakes") || !backup {
			t.Fatalf("%s: %s %t %v", value, dir, backup, err)
		}
	}
	custom := filepath.Join(home, "custom")
	dir, backup, err := resolveOutputDirectory(custom, home)
	if err != nil || dir != custom || backup {
		t.Fatalf("custom directory: %s %t %v", dir, backup, err)
	}
}

func TestGenerateAgainInDefaultDirectory(t *testing.T) {
	dir, force, err := resolveOutputDirectory(defaultOutputDirectory, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	// Generation must work even without a Nix executable. Validation may fail,
	// but both first and subsequent exports must remain available on disk.
	t.Setenv("PATH", t.TempDir())
	cfg := builder.DefaultConfig()
	options := builder.WriteOptions{Directory: dir, Force: force}
	builder.Write(context.Background(), "{ description = \"first\"; }", cfg, options)
	backup, _ := builder.Write(context.Background(), "{ description = \"second\"; }", cfg, options)
	current, err := os.ReadFile(filepath.Join(dir, "flake.nix"))
	if err != nil || string(current) != "{ description = \"second\"; }" {
		t.Fatalf("regeneration failed: %s %v", current, err)
	}
	original, err := os.ReadFile(filepath.Join(backup, "flake.nix"))
	if err != nil || string(original) != "{ description = \"first\"; }" {
		t.Fatalf("backup missing: %s %v", original, err)
	}
}
