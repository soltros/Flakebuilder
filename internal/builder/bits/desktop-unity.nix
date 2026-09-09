{ inputs, ... }: {
 imports = [ inputs.unity-on-nix.nixosModules.default ];
 services.desktopManager.unity.enable = true;
}
