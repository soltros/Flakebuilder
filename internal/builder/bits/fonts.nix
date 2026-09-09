{ config, lib, pkgs, ... }: {
fonts = {
 fontDir.enable = true;
 packages = with pkgs; [ inter open-sans roboto roboto-mono noto-fonts noto-fonts-cjk-sans noto-fonts-color-emoji dejavu_fonts liberation_ttf hack-font fira-code font-awesome nerd-fonts.symbols-only ];
 fontconfig = { enable = true; defaultFonts = {
 sansSerif = [ "Inter" "Noto Sans" "DejaVu Sans" ]; serif = [ "Noto Serif" "DejaVu Serif" ];
 monospace = [ "Roboto Mono" "Symbols Nerd Font" "Hack" "DejaVu Sans Mono" ]; emoji = [ "Noto Color Emoji" ];
 }; subpixel.rgba = "rgb"; hinting.style = "slight"; };
};
}
