{ pkgs, config, ... }: { i18n.inputMethod.ibus.panel = "${pkgs.kdePackages.plasma-desktop}/libexec/kimpanel-ibus-panel"; }
