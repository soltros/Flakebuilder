{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.go ];
}
