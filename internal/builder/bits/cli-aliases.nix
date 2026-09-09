{ config, lib, pkgs, ... }: {
programs.zsh.shellAliases = { cat = "bat --paging=never"; grep = "rg"; find = "fd"; df = "duf"; ls = "eza --icons=auto"; ll = "eza -la --icons=auto --git"; tree = "eza --tree --icons=auto"; };
}
