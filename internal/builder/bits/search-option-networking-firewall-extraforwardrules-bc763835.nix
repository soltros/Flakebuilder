{ pkgs, config, ... }: { networking.firewall.extraForwardRules = "iifname wg0 accept"; }
