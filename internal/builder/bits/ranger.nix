{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.ranger ];
}
