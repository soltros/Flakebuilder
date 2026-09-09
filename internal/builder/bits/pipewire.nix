{ config, lib, pkgs, ... }: {
services.pulseaudio.enable = false;
security.rtkit.enable = true;
services.pipewire = { enable = true; alsa.enable = true; alsa.support32Bit = pkgs.stdenv.hostPlatform.isx86_64; pulse.enable = true; };
}
