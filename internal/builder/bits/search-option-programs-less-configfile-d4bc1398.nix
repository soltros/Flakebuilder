{ pkgs, config, ... }: { programs.less.configFile = "${pkgs.my-configs}/lesskey"; }
