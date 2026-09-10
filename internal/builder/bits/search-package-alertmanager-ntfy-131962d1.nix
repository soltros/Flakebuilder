{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.alertmanager-ntfy ]; }
