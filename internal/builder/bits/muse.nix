{ pkgs, ... }:
let updater = pkgs.writeShellScript "muse-update" ''
 export MUSE_NO_MODIFY_PATH=1
 ${pkgs.curl}/bin/curl -fsSL https://dev.meta.ai/install.sh | ${pkgs.bash}/bin/bash
'';
in {
 environment.systemPackages = [ pkgs.bubblewrap ];
 systemd.user.services.muse-code-update = {
 description = "Update Muse Code CLI";
 unitConfig.ConditionUser = [[ nix .User ]];
 serviceConfig = { Type = "oneshot"; ExecStart = updater; };
 path = [ pkgs.curl pkgs.bash pkgs.coreutils ];
 };
 systemd.user.timers.muse-code-update = { wantedBy = [ "timers.target" ]; timerConfig = { OnCalendar = "weekly"; Persistent = true; }; };
}
