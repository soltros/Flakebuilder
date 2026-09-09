{ config, lib, pkgs, ... }: {
environment.systemPackages = with pkgs; [ kopia btrfs-progs ntfs3g ncdu ];
}
