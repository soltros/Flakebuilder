{ config, lib, pkgs, ... }: {
services.xserver = { enable = true; desktopManager.cinnamon.enable = true; };
}
