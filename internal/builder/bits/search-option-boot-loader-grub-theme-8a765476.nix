{ pkgs, config, ... }: { boot.loader.grub.theme = "${pkgs.kdePackages.breeze-grub}/grub/themes/breeze"; }
