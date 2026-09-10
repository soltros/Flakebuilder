{ config, lib, pkgs, ... }: {
boot.initrd.luks.devices = lib.mkDefault {};
}
