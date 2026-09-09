package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CheckHardware accepts the ordinary generated hardware expression. Local paths
// and NIX_PATH lookups would make the output depend on files outside flake.nix.
// This conservative scan is not a security boundary for untrusted Nix code.
func CheckHardware(src string) error {
	masked := maskStringsAndComments(src)
	for _, token := range strings.FieldsFunc(masked, func(r rune) bool { return strings.ContainsRune(" \t\r\n(){}[]=;,", r) }) {
		if strings.HasPrefix(token, "./") || strings.HasPrefix(token, "../") || strings.HasPrefix(token, "/") || strings.HasPrefix(token, "<") {
			return fmt.Errorf("hardware contains an external path %q; inline that configuration first", token)
		}
	}
	if strings.Contains(masked, "builtins.path") || strings.Contains(masked, "builtins.toPath") || strings.Contains(masked, "builtins.readFile") || strings.Contains(masked, "builtins.getEnv") {
		return fmt.Errorf("hardware must not read external files or environment variables")
	}
	return nil
}
func maskStringsAndComments(s string) string {
	var out strings.Builder
	for i := 0; i < len(s); {
		switch {
		case s[i] == '#':
			for i < len(s) && s[i] != '\n' {
				i++
			}
			out.WriteByte(' ')
		case strings.HasPrefix(s[i:], "/*"):
			end := strings.Index(s[i+2:], "*/")
			if end < 0 {
				return out.String()
			}
			i += end + 4
			out.WriteByte(' ')
		case s[i] == '"':
			i++
			for i < len(s) {
				if s[i] == '\\' {
					i += 2
					continue
				}
				if s[i] == '"' {
					i++
					break
				}
				i++
			}
			out.WriteByte(' ')
		case strings.HasPrefix(s[i:], "''"):
			i += 2
			for i < len(s) {
				if strings.HasPrefix(s[i:], "'''") {
					i += 3
					continue
				}
				if strings.HasPrefix(s[i:], "''${") {
					i += 4
					continue
				}
				if strings.HasPrefix(s[i:], "''") {
					i += 2
					break
				}
				i++
			}
			out.WriteByte(' ')
		default:
			out.WriteByte(s[i])
			i++
		}
	}
	return out.String()
}

type WriteOptions struct {
	Directory string
	Force     bool
	Lock      bool
	Build     bool
	Log       io.Writer
}

func run(ctx context.Context, log io.Writer, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = log
	cmd.Stderr = log
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", name, err)
	}
	return nil
}
func Parse(ctx context.Context, source string) error {
	file, err := os.CreateTemp("", "flakebuilder-parse-*.nix")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err = file.WriteString(source); err != nil {
		file.Close()
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "nix-instantiate", "--parse", file.Name())
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Nix parse failed: %w\n%s", err, stderr.String())
	}
	return nil
}

// Write validates in a staging directory. It never activates a configuration.
// Only flake.nix and (when locking/building) flake.lock are published.
func Write(ctx context.Context, source string, cfg Config, opt WriteOptions) (string, error) {
	if opt.Log == nil {
		opt.Log = io.Discard
	}
	if opt.Build && strings.TrimSpace(cfg.Hardware) == "" {
		return "", fmt.Errorf("building requires embedded hardware; supply --hardware")
	}
	if err := CheckHardware(cfg.Hardware); err != nil {
		return "", err
	}
	if err := Parse(ctx, source); err != nil {
		return "", err
	}
	dir, err := filepath.Abs(opt.Directory)
	if err != nil {
		return "", err
	}
	if err = os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	target := filepath.Join(dir, "flake.nix")
	if _, err = os.Lstat(target); err == nil && !opt.Force {
		return "", fmt.Errorf("%s exists; use --force after reviewing the preview", target)
	} else if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	stage, err := os.MkdirTemp(dir, ".flakebuilder-stage-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(stage)
	if err = os.WriteFile(filepath.Join(stage, "flake.nix"), []byte(source), 0644); err != nil {
		return "", err
	}
	names := []string{"flake.nix"}
	if opt.Lock || opt.Build {
		old, err := os.ReadFile(filepath.Join(dir, "flake.lock"))
		if err == nil {
			if err = os.WriteFile(filepath.Join(stage, "flake.lock"), old, 0644); err != nil {
				return "", err
			}
		} else if !os.IsNotExist(err) {
			return "", err
		}
		common := []string{"--extra-experimental-features", "nix-command flakes"}
		if err = run(ctx, opt.Log, "nix", append(common, "flake", "lock", "path:"+stage)...); err != nil {
			return "", err
		}
		if err = run(ctx, opt.Log, "nix", append(common, "flake", "check", "--no-build", "--no-write-lock-file", "path:"+stage)...); err != nil {
			return "", err
		}
		if opt.Build {
			if err = run(ctx, opt.Log, "nix", append(common, "build", "--no-link", "--no-write-lock-file", "path:"+stage+"#nixosConfigurations."+cfg.Host+".config.system.build.toplevel")...); err != nil {
				return "", err
			}
		}
		names = []string{"flake.lock", "flake.nix"}
	}
	// Snapshot all destinations before changing either file. Keep a durable backup
	// when replacing existing files; restore those bytes if publication fails.
	old := map[string][]byte{}
	modes := map[string]os.FileMode{}
	for _, name := range names {
		path := filepath.Join(dir, name)
		st, e := os.Lstat(path)
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return "", e
		}
		if !st.Mode().IsRegular() {
			return "", fmt.Errorf("refusing to replace non-regular file %s", path)
		}
		b, e := os.ReadFile(path)
		if e != nil {
			return "", e
		}
		old[name] = b
		modes[name] = st.Mode().Perm()
	}
	backup := ""
	if len(old) > 0 {
		backup, err = os.MkdirTemp(dir, "flakebuilder-backup-")
		if err != nil {
			return "", err
		}
		for name, b := range old {
			if err = os.WriteFile(filepath.Join(backup, name), b, modes[name]); err != nil {
				return "", err
			}
		}
	}
	published := []string{}
	for _, name := range names {
		if err = os.Rename(filepath.Join(stage, name), filepath.Join(dir, name)); err != nil {
			for _, n := range published {
				path := filepath.Join(dir, n)
				var rollback error
				if b, ok := old[n]; ok {
					rollback = os.WriteFile(path, b, modes[n])
				} else {
					rollback = os.Remove(path)
				}
				if rollback != nil {
					return backup, fmt.Errorf("publish failed: %v; restore %s from %s: %w", err, n, backup, rollback)
				}
			}
			return backup, err
		}
		published = append(published, name)
	}
	return backup, nil
}
