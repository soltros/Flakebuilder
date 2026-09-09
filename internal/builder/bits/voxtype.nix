{ inputs, pkgs, ... }: { environment.systemPackages = with inputs.voxtype.packages.${pkgs.stdenv.hostPlatform.system}; [ vulkan osd-gtk4 ]; }
