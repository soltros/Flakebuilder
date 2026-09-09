{ config, lib, pkgs, ... }: {
programs.gamemode.settings.general = { renice = 10; softrealtime = "auto"; };
}
