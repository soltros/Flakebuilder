{ config, lib, pkgs, ... }: {
systemd.user.services.flakebuilder-plasma-theme = {
 wantedBy = [ "graphical-session.target" ]; unitConfig.ConditionUser = [[ nix .User ]]; serviceConfig.Type = "oneshot";
 path = with pkgs; [ kdePackages.kconfig kdePackages.kservice coreutils dconf desktop-file-utils ];
 script = ''
      KWRITE="${pkgs.kdePackages.kconfig}/bin/kwriteconfig6"

      # 1. Dark Color Scheme & Plasma Theme
      $KWRITE --file kdeglobals --group General --key ColorScheme "BreezeDark"
      $KWRITE --file plasmarc --group Theme --key name "breeze-dark"

      # 2. Papirus-Dark Icons
      $KWRITE --file kdeglobals --group Icons --key Theme "Papirus-Dark"

      # 3. Cursor theme and size
      $KWRITE --file kcminputrc --group Mouse --key cursorTheme "breeze_cursors"
      $KWRITE --file kcminputrc --group Mouse --key cursorSize 24

      # 4. Fonts
      $KWRITE --file kdeglobals --group General --key font "Inter,9,-1,5,50,0,0,0,0,0"
      $KWRITE --file kdeglobals --group General --key fixed "Roboto Mono,10,-1,5,50,0,0,0,0,0"
      $KWRITE --file kdeglobals --group General --key smallestReadableFont "Inter,8,-1,5,50,0,0,0,0,0"
      $KWRITE --file kdeglobals --group General --key toolBarFont "Inter,9,-1,5,50,0,0,0,0,0"
      $KWRITE --file kdeglobals --group General --key menuFont "Inter,9,-1,5,50,0,0,0,0,0"

      # 5. GTK 3 & 4 dark theme, fonts, and Papirus-Dark icons
      mkdir -p "$HOME/.config/gtk-3.0" "$HOME/.config/gtk-4.0"
      cat << 'EOF' > "$HOME/.config/gtk-3.0/settings.ini"
[Settings]
gtk-theme-name=Breeze-Dark
gtk-icon-theme-name=Papirus-Dark
gtk-font-name=Inter 9
gtk-cursor-theme-name=breeze_cursors
gtk-cursor-theme-size=24
gtk-application-prefer-dark-theme=1
EOF
      cp -f "$HOME/.config/gtk-3.0/settings.ini" "$HOME/.config/gtk-4.0/settings.ini"

      /run/current-system/sw/bin/dconf write /org/gnome/desktop/interface/color-scheme "'prefer-dark'" || true
      /run/current-system/sw/bin/dconf write /org/gnome/desktop/interface/gtk-theme "'Breeze-Dark'" || true
      /run/current-system/sw/bin/dconf write /org/gnome/desktop/interface/icon-theme "'Papirus-Dark'" || true
      /run/current-system/sw/bin/dconf write /org/gnome/desktop/interface/font-name "'Inter 9'" || true
      /run/current-system/sw/bin/dconf write /org/gnome/desktop/interface/monospace-font-name "'Roboto Mono 10'" || true


 '';
};
}
