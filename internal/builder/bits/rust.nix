{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.rustc pkgs.cargo ];
}
