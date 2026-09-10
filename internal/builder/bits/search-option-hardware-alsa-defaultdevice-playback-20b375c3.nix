{ pkgs, config, ... }: { hardware.alsa.defaultDevice.playback = "dmix:CARD=1,DEV=0"; }
