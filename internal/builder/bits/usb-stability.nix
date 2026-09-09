{ config, lib, pkgs, ... }: {
boot.kernelParams = [ "pcie_aspm=off" "transparent_hugepage=never" "usbcore.autosuspend=-1" ];
services.udev.extraRules = ''
 ACTION=="add|change", SUBSYSTEM=="usb", TEST=="power/control", ATTR{power/control}="on"
 ACTION=="add|change", SUBSYSTEM=="pci", TEST=="power/control", ATTR{power/control}="on"
'';
}
