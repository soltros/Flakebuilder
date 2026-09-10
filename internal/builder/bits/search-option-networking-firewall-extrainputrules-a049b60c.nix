{ pkgs, config, ... }: { networking.firewall.extraInputRules = "ip6 saddr { fc00::/7, fe80::/10 } tcp dport 24800 accept"; }
