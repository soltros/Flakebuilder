{ config, lib, pkgs, ... }: {
hardware.firmware = [ pkgs.sof-firmware ];
hardware.alsa.enablePersistence = true;
}
