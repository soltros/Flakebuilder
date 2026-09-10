{ pkgs, config, ... }: { networking.firewall.pingLimit = "--limit 1/minute --limit-burst 5"; }
