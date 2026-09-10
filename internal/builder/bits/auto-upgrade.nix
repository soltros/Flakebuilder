{ config, lib, pkgs, ... }: {
system.autoUpgrade = { enable = true; dates = "weekly"; allowReboot = false; };
}
