{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.alpine-make-rootfs ]; }
