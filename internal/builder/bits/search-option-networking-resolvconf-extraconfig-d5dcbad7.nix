{ pkgs, config, ... }: { networking.resolvconf.extraConfig = "libc=NO"; }
