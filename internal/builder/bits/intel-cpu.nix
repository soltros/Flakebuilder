{ config, lib, pkgs, ... }: {
hardware.enableRedistributableFirmware = true;
hardware.cpu.intel.updateMicrocode = true;
}
