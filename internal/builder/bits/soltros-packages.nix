{ inputs, ... }: { nixpkgs.overlays = [ inputs.soltros-nixpkgs.overlays.default ]; }
