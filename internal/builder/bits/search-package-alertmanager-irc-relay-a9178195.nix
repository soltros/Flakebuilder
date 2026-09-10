{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.alertmanager-irc-relay ]; }
