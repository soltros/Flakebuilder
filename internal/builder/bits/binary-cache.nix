{ config, lib, pkgs, ... }: {
nix.settings.substituters = [ "https://cache.nixos.org/" ];
}
