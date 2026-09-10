{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.gcc ];
}
