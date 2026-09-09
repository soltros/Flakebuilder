{ config, lib, pkgs, ... }: {
users.users.${[[ nix .User ]]} = {
 isNormalUser = true;
 description = [[ nix .Description ]];
 extraGroups = [ "wheel" ] ++ lib.optional config.networking.networkmanager.enable "networkmanager";
};
}
