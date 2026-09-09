{ config, lib, pkgs, ... }: {
environment.systemPackages = with pkgs; [ eza bat ripgrep fd duf ];
}
