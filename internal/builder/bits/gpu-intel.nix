{ config, lib, pkgs, ... }: {
hardware.graphics = { enable = true; enable32Bit = true; extraPackages = with pkgs; [ intel-media-driver libva libva-utils vpl-gpu-rt vulkan-loader intel-gpu-tools ]; };
services.xserver.videoDrivers = [ "modesetting" ];
environment.variables.LIBVA_DRIVER_NAME = "iHD";
}
