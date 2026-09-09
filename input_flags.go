package main

import (
	"fmt"
	"github.com/soltros/Flakebuilder/internal/builder"
	"strings"
)

type repeatedFlag []string

func (f *repeatedFlag) String() string     { return strings.Join(*f, ", ") }
func (f *repeatedFlag) Set(s string) error { *f = append(*f, s); return nil }
func applyInputs(cfg *builder.Config, c *builder.Catalog, inputs, uses, follows, remove []string) error {
	if cfg.ExtraInputs == nil {
		cfg.ExtraInputs = map[string]builder.ExtraInput{}
	}
	for _, name := range remove {
		delete(cfg.ExtraInputs, name)
	}
	for _, spec := range inputs {
		pair := strings.SplitN(spec, "=", 2)
		if len(pair) != 2 {
			return fmt.Errorf("--input expects NAME=URL")
		}
		name, url := pair[0], pair[1]
		in := cfg.ExtraInputs[name]
		in.Input.URL = url
		cfg.ExtraInputs[name] = in
	}
	ensure := func(name string) (builder.ExtraInput, error) {
		if in, ok := cfg.ExtraInputs[name]; ok {
			return in, nil
		}
		if in, ok := c.Inputs[name]; ok {
			return builder.ExtraInput{Input: in}, nil
		}
		return builder.ExtraInput{}, fmt.Errorf("declare input %s with --input first", name)
	}
	for _, name := range follows {
		in, err := ensure(name)
		if err != nil {
			return err
		}
		in.Input.Follows = map[string]string{"nixpkgs": "nixpkgs"}
		cfg.ExtraInputs[name] = in
	}
	for _, spec := range uses {
		pair := strings.SplitN(spec, "=", 2)
		if len(pair) != 2 {
			return fmt.Errorf("--input-use expects NAME=overlay:ATTRIBUTE, module:ATTRIBUTE or package:ATTRIBUTE")
		}
		path := strings.SplitN(pair[1], ":", 2)
		if len(path) != 2 {
			return fmt.Errorf("--input-use expects KIND:ATTRIBUTE")
		}
		in, err := ensure(pair[0])
		if err != nil {
			return err
		}
		use := builder.InputUse{Kind: path[0], Attribute: path[1]}
		found := false
		for _, u := range in.Uses {
			if u == use {
				found = true
			}
		}
		if !found {
			in.Uses = append(in.Uses, use)
		}
		cfg.ExtraInputs[pair[0]] = in
	}
	for name, in := range cfg.ExtraInputs {
		if err := builder.ValidateExtraInput(name, in); err != nil {
			return err
		}
	}
	return nil
}
