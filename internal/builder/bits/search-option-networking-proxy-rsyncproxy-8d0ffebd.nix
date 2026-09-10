{ pkgs, config, ... }: { networking.proxy.rsyncProxy = "http://127.0.0.1:3128"; }
