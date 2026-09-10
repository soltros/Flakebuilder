{ pkgs, config, ... }: { networking.nat.extraStopCommands = "iptables -D INPUT -p icmp -j ACCEPT || true"; }
