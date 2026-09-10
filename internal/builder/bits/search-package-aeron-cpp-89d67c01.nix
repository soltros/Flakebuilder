{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.aeron-cpp ]; }
