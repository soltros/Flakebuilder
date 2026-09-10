{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.add-determinism ]; }
