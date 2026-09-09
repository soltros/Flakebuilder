{ config, lib, pkgs, ... }: {
programs.zsh.shellAliases = {
 nrb = "sudo nixos-rebuild switch --flake " + lib.escapeShellArg ([[ nix .ConfigDir ]] + "#" + [[ nix .Host ]]);
 nrb-test = "sudo nixos-rebuild test --flake " + lib.escapeShellArg ([[ nix .ConfigDir ]] + "#" + [[ nix .Host ]]);
 nrb-boot = "sudo nixos-rebuild boot --flake " + lib.escapeShellArg ([[ nix .ConfigDir ]] + "#" + [[ nix .Host ]]);
 nfu = "nix flake update --flake " + lib.escapeShellArg [[ nix .ConfigDir ]];
};
}
