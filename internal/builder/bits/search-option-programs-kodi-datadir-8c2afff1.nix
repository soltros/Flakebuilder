{ pkgs, config, ... }: { programs.kodi.datadir = "${config.xdg.dataHome}/kodi"; }
