{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.ubootTools ];
}
