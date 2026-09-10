{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.amazon-ecr-credential-helper ]; }
