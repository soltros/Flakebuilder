{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.adbfs-rootless ]; }
