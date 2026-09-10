{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.rsync ];
}
