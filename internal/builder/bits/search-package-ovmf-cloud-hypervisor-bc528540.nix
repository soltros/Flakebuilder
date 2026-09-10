{ pkgs, config, ... }: { environment.systemPackages = [ pkgs.OVMF-cloud-hypervisor ]; }
