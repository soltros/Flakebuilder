package builder

import (
	"fmt"
	"regexp"
	"strings"
)

// ExtraInput is a user's explicit input declaration and optional consumers.
// Uses are structured output paths, never executable Nix snippets.
type ExtraInput struct {
	Input Input      `json:"input"`
	Uses  []InputUse `json:"uses,omitempty"`
}
type InputUse struct {
	Kind      string `json:"kind"` // overlay, module, package
	Attribute string `json:"attribute"`
}

var inputNamePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]*$`)
var attrPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+(\.[A-Za-z0-9_-]+)*$`)

func ValidateExtraInput(name string, in ExtraInput) error {
	if !inputNamePattern.MatchString(name) || name == "self" || name == "nixpkgs" {
		return fmt.Errorf("invalid/reserved input name %q", name)
	}
	url := in.Input.URL
	if !(strings.HasPrefix(url, "github:") || strings.HasPrefix(url, "gitlab:") || strings.HasPrefix(url, "git+https://") || strings.HasPrefix(url, "git+ssh://") || strings.HasPrefix(url, "https://")) {
		return fmt.Errorf("input %s needs a github:, gitlab:, git+https://, git+ssh:// or https:// URL", name)
	}
	for _, r := range url {
		if r <= 32 || r == 127 {
			return fmt.Errorf("input URL must not contain whitespace or control characters")
		}
	}
	for child, target := range in.Input.Follows {
		if !inputNamePattern.MatchString(child) || target != "nixpkgs" {
			return fmt.Errorf("%s: follows must map a child input to root nixpkgs", name)
		}
	}
	seen := map[string]bool{}
	for _, use := range in.Uses {
		if use.Kind != "overlay" && use.Kind != "module" && use.Kind != "package" {
			return fmt.Errorf("input %s: use must be overlay, module or package", name)
		}
		if !attrPattern.MatchString(use.Attribute) {
			return fmt.Errorf("invalid %s attribute %q", use.Kind, use.Attribute)
		}
		key := use.Kind + ":" + use.Attribute
		if seen[key] {
			return fmt.Errorf("duplicate input use %s for %s", key, name)
		}
		seen[key] = true
	}
	return nil
}
func InputExpression(name string, use InputUse) string {
	prefix := "inputs." + NixString(name) + "."
	switch use.Kind {
	case "overlay":
		prefix += "overlays"
	case "module":
		prefix += "nixosModules"
	case "package":
		prefix += "packages.${pkgs.stdenv.hostPlatform.system}"
	}
	for _, part := range strings.Split(use.Attribute, ".") {
		prefix += "." + NixString(part)
	}
	return prefix
}
