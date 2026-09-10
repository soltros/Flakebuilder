{ pkgs, config, ... }: { networking.firewall.extraReversePathFilterRules = "fib daddr . mark . iif type local accept"; }
