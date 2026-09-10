{ pkgs, config, ... }: { programs.claude-code.configDir = "${config.xdg.configHome}/claude"; }
