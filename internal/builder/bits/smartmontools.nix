{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.smartmontools ];
}
