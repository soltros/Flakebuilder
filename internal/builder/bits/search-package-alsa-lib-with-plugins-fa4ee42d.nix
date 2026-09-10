{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.alsa-lib-with-plugins ]; }
