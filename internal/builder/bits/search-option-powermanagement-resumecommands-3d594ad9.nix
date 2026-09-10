{ pkgs, config, ... }: { powerManagement.resumeCommands = "${pkgs.util-linux}/bin/rfkill unblock all"; }
