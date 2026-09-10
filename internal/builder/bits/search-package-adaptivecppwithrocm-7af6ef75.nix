{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.adaptivecppWithRocm ]; }
