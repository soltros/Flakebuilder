{ pkgs, config, ... }: { environment.systemPackages = [ pkgs._389-ds-base ]; }
