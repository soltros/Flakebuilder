{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.adaptivecppWithCuda ]; }
