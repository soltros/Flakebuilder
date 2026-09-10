{
  description = "Flakebuilder — assemble a single NixOS flake from selectable bits";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  outputs = { self, nixpkgs }: let
    systems = [ "x86_64-linux" "aarch64-linux" ];
    forAllSystems = nixpkgs.lib.genAttrs systems;
  in {
    packages = forAllSystems (system: let
      pkgs = import nixpkgs { inherit system; };
    in {
      default = pkgs.buildGoModule {
        pname = "flakebuilder";
        version = "0.1.1";
        src = self;
        vendorHash = "sha256-uwBJAqN4sIepiiJf9lCDumLqfKJEowQO2tOiSWD3Fig=";
        env.CGO_ENABLED = 0;
        nativeCheckInputs = [ pkgs.nix ];
        preCheck = ''
          export NIX_STATE_DIR="$TMPDIR/nix-state"
          export NIX_LOG_DIR="$TMPDIR/nix-log"
          export NIX_CONF_DIR="$TMPDIR/nix-conf"
          export NIX_STORE_DIR="$TMPDIR/nix-store"
        '';
        nativeBuildInputs = [ pkgs.makeWrapper ];
        postInstall = ''
          mv $out/bin/Flakebuilder $out/bin/flakebuilder
          wrapProgram $out/bin/flakebuilder --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.nix ]}
        '';
        meta = {
          description = "Menu-driven single-file NixOS flake generator";
          homepage = "https://github.com/soltros/Flakebuilder";
          license = pkgs.lib.licenses.gpl3Only;
          mainProgram = "flakebuilder";
        };
      };
      gui = pkgs.stdenv.mkDerivation {
        pname = "flakebuilder-gui"; version = "0.1.0"; src = ./gui;
        nativeBuildInputs = with pkgs; [ meson ninja pkg-config vala wrapGAppsHook4 glib ];
        buildInputs = [ pkgs.gtk4 pkgs.pantheon.granite7 ];
        meta = { description = "Pantheon-style GTK frontend for Flakebuilder"; license = pkgs.lib.licenses.gpl3Only; mainProgram = "flakebuilder-gui"; };
      };
    });
    apps = forAllSystems (system: {
      default = { type = "app"; program = "${self.packages.${system}.default}/bin/flakebuilder"; };
      gui = { type = "app"; program = "${self.packages.${system}.gui}/bin/flakebuilder-gui"; };
    });
    checks = forAllSystems (system: { app = self.packages.${system}.default; });
    devShells = forAllSystems (system: let pkgs = import nixpkgs { inherit system; }; in {
      default = pkgs.mkShell { packages = [ pkgs.go pkgs.nix ]; CGO_ENABLED = "0"; };
    });
  };
}
