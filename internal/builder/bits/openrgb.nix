{ config, lib, pkgs, ... }: {
services.udev.packages = [ pkgs.openrgb ]; boot.kernelModules = [ "i2c-dev" ];
}
