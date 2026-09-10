{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.abi-compliance-checker ]; }
