{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.rclone ];
}
