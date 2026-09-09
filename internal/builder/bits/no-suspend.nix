{ config, lib, pkgs, ... }: {
services.logind.settings.Login = {
 IdleAction = "ignore"; HandleHibernateKey = "ignore"; HandleLidSwitch = "ignore";
 HandleLidSwitchExternalPower = "ignore"; HandleSuspendKey = "ignore";
};
systemd.sleep.settings.Sleep = { AllowHibernation = "no"; AllowHybridSleep = "no"; AllowSuspend = "no"; AllowSuspendThenHibernate = "no"; };
}
