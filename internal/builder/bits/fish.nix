{ config, lib, pkgs, ... }: {
programs.fish.enable = true;
users.users.${[[ nix .User ]]}.shell = pkgs.fish;
}
