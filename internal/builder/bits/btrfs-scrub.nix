{ config, lib, pkgs, ... }: {
services.btrfs.autoScrub = { enable = true; interval = "weekly"; };
}
