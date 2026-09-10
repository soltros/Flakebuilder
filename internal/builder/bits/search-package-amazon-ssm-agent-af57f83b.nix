{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.amazon-ssm-agent ]; }
