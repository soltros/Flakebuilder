{ config, lib, pkgs, ... }: {
powerManagement = { enable = true; cpuFreqGovernor = "schedutil"; powertop.enable = true; };
}
