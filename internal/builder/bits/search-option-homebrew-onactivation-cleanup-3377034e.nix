{ pkgs, config, ... }: { homebrew.onActivation.cleanup = "uninstall"; }
