{ pkgs, config, ... }: { boot.initrd.systemd.network.wait-online.timeout = 0; }
