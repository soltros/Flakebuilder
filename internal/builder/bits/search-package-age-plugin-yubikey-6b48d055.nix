{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.age-plugin-yubikey ]; }
