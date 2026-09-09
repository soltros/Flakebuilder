{ config, lib, pkgs, ... }: {
users.users.${[[ nix .User ]]}.extraGroups = [ config.services.hermes-agent.group ];
systemd.tmpfiles.rules = [
 ("z /var/lib/hermes/.hermes/auth.json 0660 " + config.services.hermes-agent.user + " " + config.services.hermes-agent.group + " -")
 ("d " + builtins.toJSON [[ nix .Workspace ]] + " 0755 " + [[ nix .User ]] + " users -")
];
}
