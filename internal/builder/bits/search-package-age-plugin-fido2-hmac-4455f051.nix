{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.age-plugin-fido2-hmac ]; }
