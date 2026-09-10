{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.alertmanager-webhook-logger ]; }
