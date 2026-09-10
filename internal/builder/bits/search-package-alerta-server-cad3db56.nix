{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.alerta-server ]; }
