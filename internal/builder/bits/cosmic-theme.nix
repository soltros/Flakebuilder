{ pkgs, ... }:

let
  catppuccinCosmicRevision = "95e81098042dd2102f0b258f6990f886c5759692";

  cosmicAppearance = pkgs.fetchurl {
    url = "https://raw.githubusercontent.com/catppuccin/cosmic-desktop/${catppuccinCosmicRevision}/themes/cosmic-settings/catppuccin-mocha-mauve%2Bround.ron";
    hash = "sha256-/XUIANPnrO/nRHXMfjIWb+sNmSjNGIaLU4EJ/HP7IMI=";
  };

  cosmicTerminal = pkgs.fetchurl {
    url = "https://raw.githubusercontent.com/catppuccin/cosmic-desktop/${catppuccinCosmicRevision}/themes/cosmic-term/catppuccin-mocha.ron";
    hash = "sha256-jk2gZKsdgtENjvRKvWkwbb3MvCBU+1wyi+1ovREspVQ=";
  };
in
{
  # COSMIC currently imports themes through its Appearance UI. Keeping the
  # pinned files in /etc makes the selected theme reproducible and easy to
  # re-import after a fresh install.
  environment.etc = {
    "cosmic/themes/catppuccin-mocha-mauve+round.ron".source = cosmicAppearance;
    "cosmic/themes/catppuccin-mocha-terminal.ron".source = cosmicTerminal;
  };

  environment.systemPackages = with pkgs; [
    catppuccin-cursors.mochaMauve
    papirus-icon-theme
  ];

  environment.sessionVariables = {
    XCURSOR_THEME = "catppuccin-mocha-mauve-cursors";
    XCURSOR_SIZE = "24";
  };
}
