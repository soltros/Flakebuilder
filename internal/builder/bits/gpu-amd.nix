{ config, lib, pkgs, ... }: {
boot.initrd.kernelModules = [ "amdgpu" ];
services.xserver.videoDrivers = [ "amdgpu" ];
hardware.graphics = { enable = true; enable32Bit = pkgs.stdenv.hostPlatform.isx86_64; };
}
