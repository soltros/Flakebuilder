{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.abseil-cpp ]; }
