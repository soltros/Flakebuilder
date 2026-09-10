{ config, lib, pkgs, ... }: {
hardware.graphics.extraPackages = [ pkgs.mesa.drivers ];
}
