{ pkgs, config, ... }: { programs.astroid.externalEditor = "nvim-qt -- -c 'set ft=mail' '+set fileencoding=utf-8' '+set ff=unix' '+set enc=utf-8' '+set fo+=w' %1"; }
