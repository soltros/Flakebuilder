{ config, lib, pkgs, ... }: {
environment.systemPackages = [ pkgs.gstreamer ];
}
