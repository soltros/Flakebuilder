{ pkgs, config, ... }: { networking.firewall.extraStopCommands = "iptables -P INPUT ACCEPT"; }
