{ pkgs, config, ... }: { boot.initrd.systemd.tpm2.pcrphases.enable = true; }
