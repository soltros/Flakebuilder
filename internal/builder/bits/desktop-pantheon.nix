{ config, lib, pkgs, ... }: {
services.xserver.enable = true;
services.xserver.displayManager.lightdm.enable = true;
services.desktopManager.pantheon = { enable = true; extraSwitchboardPlugs = [ pkgs.pantheon-tweaks ]; };
environment.systemPackages = [ pkgs.pantheon-tweaks ];
}
