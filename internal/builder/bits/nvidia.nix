{ config, lib, pkgs, ... }: {
services.xserver.videoDrivers = [ "nvidia" ]; hardware.graphics.enable = true;
}
