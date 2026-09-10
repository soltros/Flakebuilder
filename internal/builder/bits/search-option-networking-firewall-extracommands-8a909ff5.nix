{ pkgs, config, ... }: { networking.firewall.extraCommands = "iptables -A INPUT -p icmp -j ACCEPT"; }
