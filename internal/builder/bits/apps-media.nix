{ config, lib, pkgs, ... }: {
environment.systemPackages = with pkgs; [ vlc gimp gthumb pinta yt-dlp ];
}
