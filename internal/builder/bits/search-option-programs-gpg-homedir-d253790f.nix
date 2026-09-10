{ pkgs, config, ... }: { programs.gpg.homedir = "${config.xdg.dataHome}/gnupg"; }
