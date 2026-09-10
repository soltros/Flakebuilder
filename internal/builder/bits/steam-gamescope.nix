{ config, lib, pkgs, ... }: {
programs.gamescope = { enable = true; capSysNice = true; }; programs.steam.gamescopeSession.enable = true;
}
