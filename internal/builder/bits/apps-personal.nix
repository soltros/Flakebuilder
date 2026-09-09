{ config, lib, pkgs, ... }: {
environment.systemPackages = with pkgs; [ bitwarden-desktop spotify winetricks wine64 appimage-run distrobox wget screen unzip pavucontrol pamixer caffeine-ng adapta-gtk-theme mlocate lxrandr python311Packages.pip ];
}
