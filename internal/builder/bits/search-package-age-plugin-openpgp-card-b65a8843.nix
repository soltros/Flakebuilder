{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.age-plugin-openpgp-card ]; }
