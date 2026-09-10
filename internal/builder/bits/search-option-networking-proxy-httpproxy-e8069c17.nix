{ pkgs, config, ... }: { networking.proxy.httpProxy = "http://127.0.0.1:3128"; }
