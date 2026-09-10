{ pkgs, config, ... }: { networking.proxy.httpsProxy = "http://127.0.0.1:3128"; }
