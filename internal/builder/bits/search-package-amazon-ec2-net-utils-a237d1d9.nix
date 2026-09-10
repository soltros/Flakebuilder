{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.amazon-ec2-net-utils ]; }
