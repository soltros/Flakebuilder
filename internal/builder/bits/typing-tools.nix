{ config, lib, pkgs, ... }: {
hardware.uinput.enable = true; programs.ydotool.enable = true;
environment.systemPackages = with pkgs; [ wtype wl-clipboard ydotool dotool ];
users.users.${[[ nix .User ]]}.extraGroups = [ "input" "uinput" ];
}
