{ pkgs, config, ... }: { programs.lutris.winePackages = "[ pkgs.wineWow64Packages.full ]"; }
