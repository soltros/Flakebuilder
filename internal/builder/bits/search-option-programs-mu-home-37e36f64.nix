{ pkgs, config, ... }: { programs.mu.home = "\${config.home.homeDirectory}/Maildir/.mu"; }
