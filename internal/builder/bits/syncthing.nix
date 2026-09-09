{ config, lib, pkgs, ... }: {
services.syncthing = { enable = true; user = [[ nix .User ]]; dataDir = [[ nix (home) ]]; };
}
