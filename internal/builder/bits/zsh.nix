{ config, lib, pkgs, ... }: {
programs.zsh = { enable = true; autosuggestions.enable = true; ohMyZsh = { enable = true; plugins = [ "git" "sudo" ]; }; };
users.users.${[[ nix .User ]]}.shell = pkgs.zsh;
}
