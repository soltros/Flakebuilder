{ config, lib, pkgs, ... }: {
environment.systemPackages = with pkgs; [ git python312 nodejs pipx php gh lazygit geany fresh-editor ];
}
