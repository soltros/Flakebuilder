{ config, lib, pkgs, ... }: {
services.xserver = { enable = true; windowManager.i3.enable = true; };
}
