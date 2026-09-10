{ pkgs, config, ... }: { boot.initrd.systemd.network.wait-online.enable = false; }
