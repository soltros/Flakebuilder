{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.borgbackup ];
}
