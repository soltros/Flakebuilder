{ pkgs, config, ... }: { networking.proxy.noProxy = "127.0.0.1,localhost,.localdomain"; }
