{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.alsa-lib ]; }
