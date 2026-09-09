package main

import (
	"github.com/soltros/Flakebuilder/internal/builder"
	"testing"
)

func TestInputFlags(t *testing.T) {
	c, e := builder.Load()
	if e != nil {
		t.Fatal(e)
	}
	cfg := builder.DefaultConfig()
	e = applyInputs(&cfg, c, []string{"other=github:example/repo"}, []string{"other=overlay:default", "other=package:cli"}, []string{"other"}, nil)
	if e != nil {
		t.Fatal(e)
	}
	in := cfg.ExtraInputs["other"]
	if len(in.Uses) != 2 || in.Input.Follows["nixpkgs"] != "nixpkgs" {
		t.Fatal(in)
	}
	if e = applyInputs(&cfg, c, nil, []string{"absent=overlay:default"}, nil, nil); e == nil {
		t.Fatal("accepted undeclared input")
	}
}
