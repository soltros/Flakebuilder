{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.telegram-desktop ];
}
