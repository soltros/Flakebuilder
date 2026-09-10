{ pkgs, config, ... }: { boot.loader.grub.extraPerEntryConfig = "root (hd0)"; }
