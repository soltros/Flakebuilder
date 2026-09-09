{ config, lib, pkgs, ... }: {
boot.kernelParams = [ "radeon.si_support=0" "radeon.cik_support=0" "amdgpu.si_support=1" "amdgpu.cik_support=1" ];
hardware.graphics.extraPackages = [ pkgs.rocmPackages.clr.icd ];
environment.variables.ROC_ENABLE_PRE_VEGA = "1";
environment.systemPackages = [ pkgs.clinfo ];
}
