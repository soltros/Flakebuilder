{ pkgs, config, ... }: { programs.foot.server.systemdTarget = "sway-session.target"; }
