{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.adb-sync ]; }
