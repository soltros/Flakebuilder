{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.OVMF-xen ]; }
