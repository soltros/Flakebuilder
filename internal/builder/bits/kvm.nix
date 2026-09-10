{ config, lib, pkgs, ... }: {
boot.kernelModules = [ "kvm-intel" "kvm-amd" ];
}
