{ config, lib, pkgs, ... }: {
powerManagement = { enable = true; cpuFreqGovernor = "performance"; powertop.enable = false; };
services.power-profiles-daemon.enable = false;
}
