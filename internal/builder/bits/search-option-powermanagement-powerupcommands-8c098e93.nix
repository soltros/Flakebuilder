{ pkgs, config, ... }: { powerManagement.powerUpCommands = "${pkgs.powertop}/bin/powertop --auto-tune"; }
