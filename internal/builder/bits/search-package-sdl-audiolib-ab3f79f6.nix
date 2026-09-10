{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.SDL_audiolib ]; }
