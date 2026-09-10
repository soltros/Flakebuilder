{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.fuzzel ];
}
