{ config, lib, pkgs, ... }: {
nixpkgs.overlays = [ (final: prev: { pipx = prev.pipx.overridePythonAttrs (_: { doCheck = false; }); }) ];
}
