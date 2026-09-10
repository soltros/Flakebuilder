{ pkgs, config, ... }: { environment.wvdial.dialerDefaults = "Init1 = AT+CGDCONT=1,\"IP\",\"internet.t-mobile\""; }
