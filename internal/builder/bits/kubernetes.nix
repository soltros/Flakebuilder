{ config, lib, pkgs, ... }: {
services.kubernetes.kubelet.enable = true;
}
