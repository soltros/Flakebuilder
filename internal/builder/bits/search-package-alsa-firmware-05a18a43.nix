{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.alsa-firmware ]; }
