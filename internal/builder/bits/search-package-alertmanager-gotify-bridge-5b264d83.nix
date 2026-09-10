{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.alertmanager-gotify-bridge ]; }
