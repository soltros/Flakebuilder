{ pkgs, config, ... }: { boot.loader.grub.device = "/dev/disk/by-id/wwn-0x500001234567890a"; }
