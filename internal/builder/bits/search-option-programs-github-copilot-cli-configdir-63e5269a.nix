{ pkgs, config, ... }: { programs.github-copilot-cli.configDir = "${config.xdg.configHome}/copilot"; }
