{ config, lib, pkgs, ... }: {
boot.kernelModules = [ "virtiofs" ];
}
