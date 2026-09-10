{ pkgs, config, ... }: { boot.postBootCommands = "rm -f /var/log/messages"; }
