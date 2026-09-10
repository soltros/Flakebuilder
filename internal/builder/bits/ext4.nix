{ config, lib, pkgs, ... }: {
boot.supportedFilesystems = [ "ext4" ]; environment.systemPackages = [ pkgs.e2fsprogs ];
}
