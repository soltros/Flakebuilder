{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.openvpn ];
}
