{ config, lib, pkgs, ... }: {
boot.kernelParams = [ "i915.enable_fbc=1" "i915.enable_guc=3" "i915.enable_psr=1" "i915.enable_dc=1" ];
}
