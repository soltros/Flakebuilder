{ pkgs, config, ... }: { powerManagement.bootCommands = "${pkgs.networkmanager}/bin/nmcli radio wifi on"; }
