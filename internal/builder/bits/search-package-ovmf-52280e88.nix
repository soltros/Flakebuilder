{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.OVMF ]; }
