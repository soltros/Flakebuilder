{ pkgs, config, ... }: { powerManagement.powerDownCommands = "${pkgs.hdparm}/sbin/hdparm -B 255 /dev/sda"; }
