{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.clang ];
}
