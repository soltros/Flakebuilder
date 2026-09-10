{ config, lib, pkgs, ... }: {
boot.supportedFilesystems = [ "zfs" ]; services.zfs.autoSnapshot.enable = true;
}
