{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.amazon-cloudwatch-agent ]; }
