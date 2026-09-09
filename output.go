package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultOutputDirectory = "~/generated_flakes"

// outputDirectory expands the documented default (including a quoted ~/ path).
// The designated generation directory can be reused with automatic backups;
// other directories retain explicit overwrite protection.
func outputDirectory(value string) (string, bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", false, fmt.Errorf("locate generated_flakes directory: %w", err)
	}
	return resolveOutputDirectory(value, home)
}

func resolveOutputDirectory(value, home string) (string, bool, error) {
	defaultDir := filepath.Join(home, "generated_flakes")
	if value == "~" {
		value = home
	} else if strings.HasPrefix(value, "~/") {
		value = filepath.Join(home, value[2:])
	}
	dir, err := filepath.Abs(value)
	if err != nil {
		return "", false, err
	}
	return dir, dir == defaultDir, nil
}
