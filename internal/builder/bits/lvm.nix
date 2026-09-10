{ config, lib, pkgs, ... }: {
boot.initrd.kernelModules = [ "dm_snapshot" "dm_crypt" ]; environment.systemPackages = [ pkgs.lvm2 ];
}
