{ inputs, pkgs, ... }: { environment.systemPackages = with inputs.antigravity-nix.packages.${pkgs.stdenv.hostPlatform.system}; [ default google-antigravity-ide google-antigravity-cli ]; }
