{ pkgs, config, ... }: { hardware.alsa.defaultDevice.capture = "dsnoop:CARD=0,DEV=2"; }
