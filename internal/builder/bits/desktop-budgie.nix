{ config, lib, pkgs, ... }: {
services.xserver = { enable = true; desktopManager.budgie.enable = true; };
}
