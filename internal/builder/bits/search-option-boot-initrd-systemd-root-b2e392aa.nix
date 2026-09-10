{ pkgs, config, ... }: { boot.initrd.systemd.root = "gpt-auto"; }
