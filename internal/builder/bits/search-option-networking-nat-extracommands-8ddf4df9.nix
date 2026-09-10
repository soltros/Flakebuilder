{ pkgs, config, ... }: { networking.nat.extraCommands = "iptables -A INPUT -p icmp -j ACCEPT"; }
