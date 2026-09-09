{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.waterfox ];
}
