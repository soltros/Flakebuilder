{ config, lib, pkgs, ... }: {
boot.kernelParams = [ "amd_iommu=on" "intel_iommu=on" "iommu=pt" ];
}
