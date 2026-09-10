{ config, lib, pkgs, ... }: {
services.avahi = { enable = true; publish.enable = true; publish.userServices = true; };
}
