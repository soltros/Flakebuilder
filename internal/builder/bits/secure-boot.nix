{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.sbctl ];
}
