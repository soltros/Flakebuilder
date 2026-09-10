{ config, lib, pkgs, ... }: {
virtualisation.docker.enable = true; users.users.${[[ nix .User ]]}.extraGroups = [ "docker" ];
}
