{ config, lib, pkgs, ... }: {
xdg.portal = { enable = true; extraPortals = [ pkgs.xdg-desktop-portal-gtk ]; config.common.default = lib.mkDefault [ "gtk" ]; };
}
