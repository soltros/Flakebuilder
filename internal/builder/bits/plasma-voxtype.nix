{ config, lib, pkgs, ... }: {
systemd.user.services.flakebuilder-voxtype-shortcut = {
 wantedBy = [ "graphical-session.target" ]; unitConfig.ConditionUser = [[ nix .User ]]; serviceConfig.Type = "oneshot";
 path = with pkgs; [ kdePackages.kconfig kdePackages.kservice coreutils dconf desktop-file-utils ];
 script = ''
      KWRITE="${pkgs.kdePackages.kconfig}/bin/kwriteconfig6"
      # 6. Voxtype Shortcut (Numpad + / KP_Add)
      mkdir -p "$HOME/.local/share/applications"
      cat << 'EOF' > "$HOME/.local/share/applications/voxtype-toggle.desktop"
[Desktop Entry]
Name=Voxtype Toggle Dictation
Exec=voxtype record toggle
Type=Application
Terminal=false
Icon=audio-input-microphone
NoDisplay=true
StartupNotify=false
EOF

      $KWRITE --file kglobalshortcutsrc --group "voxtype-toggle.desktop" --key "_k_friendly_name" "Voxtype Toggle Dictation"
      $KWRITE --file kglobalshortcutsrc --group "voxtype-toggle.desktop" --key "_launch" "KP_Add,none,Voxtype Toggle Dictation"


 '';
};
}
