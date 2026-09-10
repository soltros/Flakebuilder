{ pkgs, config, ... }: { nix.gc.options = "--max-freed $((64 * 1024**3))"; }
