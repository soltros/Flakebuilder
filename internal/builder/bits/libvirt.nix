{ config, lib, pkgs, ... }: {
virtualisation.libvirtd.enable = true;
environment.systemPackages = [ pkgs.virt-manager ];
}
