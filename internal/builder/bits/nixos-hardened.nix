{ config, lib, pkgs, ... }: {
boot.kernelPackages = pkgs.linuxKernel.packages.linux_hardened; security.lockKernelModules = true; security.protectKernelImage = true;
}
