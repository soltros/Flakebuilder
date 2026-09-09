{ inputs, ... }: {
 imports = [ inputs.hermes-agent.nixosModules.default ];
 services.hermes-agent = {
 enable = true; addToSystemPackages = true; container.enable = false;
 settings = { model = { provider = "openai-codex"; default = [[ nix .HermesModel ]]; }; toolsets = [ "all" "skills" ]; terminal = { backend = "local"; timeout = 180; }; };
 environmentFiles = [ [[ nix .HermesEnv ]] ];
 };
}
