{ config, lib, pkgs, ... }: {
services.xserver = { enable = true; desktopManager.lxqt.enable = true; };
}
