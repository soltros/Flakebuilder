{ pkgs, config, ... }: { programs.less.lessopen = "|${pkgs.lesspipe}/bin/lesspipe.sh %s"; }
