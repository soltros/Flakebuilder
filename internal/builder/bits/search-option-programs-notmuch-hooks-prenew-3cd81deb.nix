{ pkgs, config, ... }: { programs.notmuch.hooks.preNew = "mbsync --all"; }
