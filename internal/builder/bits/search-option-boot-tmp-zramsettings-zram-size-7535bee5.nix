{ pkgs, config, ... }: { boot.tmp.zramSettings.zram-size = "min(ram / 2, 4096)"; }
