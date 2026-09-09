{ config, lib, pkgs, ... }:
let
  hermesHome = "/var/lib/hermes/.hermes";
  hermesUser = config.services.hermes-agent.user;
  hermesGroup = config.services.hermes-agent.group;
  yaml = pkgs.formats.yaml { };

  codexModel = model: {
    provider = "openai-codex";
    default = model;
    base_url = "https://chatgpt.com/backend-api/codex";
    api_mode = "codex_responses";
  };

  supportAuxiliary = builtins.listToAttrs (map (name: {
    inherit name;
    value = {
      provider = "openai-codex";
      model = "gpt-5.4-mini";
    };
  }) [
    "approval"
    "compression"
    "curator"
    "goal_judge"
    "mcp"
    "profile_describer"
    "skill_search"
    "title_generation"
    "triage_specifier"
    "web_extract"
  ]);

  moaConfig = {
    default_preset = "codex-hard-synthesis";
    privacy_filter = "display";
    presets.codex-hard-synthesis = {
      enabled = true;
      fanout = "user_turn";
      reference_max_tokens = 600;
      reference_models = [
        {
          provider = "openai-codex";
          model = "gpt-5.6-sol";
          reasoning_effort = "medium";
        }
        {
          provider = "openai-codex";
          model = "gpt-5.5";
          reasoning_effort = "low";
        }
        {
          provider = "openai-codex";
          model = "gpt-5.4-mini";
          reasoning_effort = "minimal";
        }
      ];
      aggregator = {
        provider = "openai-codex";
        model = "gpt-5.6-terra";
        reasoning_effort = "high";
      };
    };
  };

  commonConfig = {
    toolsets = [ "all" ];
    terminal = {
      backend = "local";
      cwd = [[ nix .Workspace ]];
      timeout = 180;
    };
    moa = moaConfig;
    auxiliary = supportAuxiliary;
  };

  supportConfig = commonConfig // {
    model = (codexModel "gpt-5.4-mini") // {
      openai_runtime = "auto";
    };
  };

  architectConfig = commonConfig // {
    model = (codexModel "gpt-5.6-terra") // {
      openai_runtime = "auto";
    };
  };

  builderConfig = commonConfig // {
    model = (codexModel "gpt-5.6-sol") // {
      openai_runtime = "codex_app_server";
    };
  };

  reviewConfig = commonConfig // {
    model = (codexModel "gpt-5.6-luna") // {
      openai_runtime = "auto";
    };
  };

  profileYaml = description: displayName: yaml.generate "profile-${displayName}.yaml" {
    inherit description;
    description_auto = false;
    display_name = displayName;
  };

  roleSoul = role: model: purpose: ''
    You are the operator's ${role} Hermes profile.

    Role: ${purpose}
    Provider: openai-codex
    Default model: ${model}

    Follow the operator's global instructions, keep work Codex-first, and stay within this profile's lane unless explicitly asked to switch roles.
  '';

  profiles = {
    architect = {
      config = yaml.generate "hermes-architect-config.yaml" architectConfig;
      profile = profileYaml "Planning, design, synthesis, and architecture" "Architect";
      soul = pkgs.writeText "hermes-architect-SOUL.md" (roleSoul "architect" "gpt-5.6-terra" "system design, architecture decisions, tradeoff analysis, planning multi-step work, and drafting specs/prompts.");
    };
    builder = {
      config = yaml.generate "hermes-builder-config.yaml" builderConfig;
      profile = profileYaml "Implementation, refactors, edits, and code execution" "Builder";
      soul = pkgs.writeText "hermes-builder-SOUL.md" (roleSoul "builder" "gpt-5.6-sol" "coding, implementation, refactors, file edits, patching, and PR-style work. Codex app-server runtime is enabled for this lane.");
    };
    review = {
      config = yaml.generate "hermes-review-config.yaml" reviewConfig;
      profile = profileYaml "Code review, QA, and final checks" "Review";
      soul = pkgs.writeText "hermes-review-SOUL.md" (roleSoul "review" "gpt-5.6-luna" "code review, QA, alternative solutions, final sanity checks, and what-did-we-miss passes.");
    };
  };

  installProfile = name: data: ''
    _profile_dir="${hermesHome}/profiles/${name}"
    install -d -m 2770 -o ${hermesUser} -g ${hermesGroup} "$_profile_dir"
    install -d -m 2770 -o ${hermesUser} -g ${hermesGroup} "$_profile_dir/memories" "$_profile_dir/cron" "$_profile_dir/logs" "$_profile_dir/sessions"
    install -m 0660 -o ${hermesUser} -g ${hermesGroup} ${data.config} "$_profile_dir/config.yaml"
    install -m 0660 -o ${hermesUser} -g ${hermesGroup} ${data.profile} "$_profile_dir/profile.yaml"
    install -m 0660 -o ${hermesUser} -g ${hermesGroup} ${data.soul} "$_profile_dir/SOUL.md"
    ln -sfn ../../.env "$_profile_dir/.env"
    ln -sfn ../../auth.json "$_profile_dir/auth.json"
    chown -h ${hermesUser}:${hermesGroup} "$_profile_dir/.env" "$_profile_dir/auth.json"
    if [ ! -e "$_profile_dir/skills" ] && [ -d "${hermesHome}/skills" ]; then
      cp -a "${hermesHome}/skills" "$_profile_dir/skills"
      chown -R ${hermesUser}:${hermesGroup} "$_profile_dir/skills"
    fi
    chmod -R u+rwX,g+rwX,o-rwx "$_profile_dir"
    find "$_profile_dir" -type d -exec chmod g+s {} +
  '';

  profileAlias = name: pkgs.writeShellScriptBin name ''
    export HERMES_HOME=${hermesHome}
    exec /run/current-system/sw/bin/hermes -p ${name} "$@"
  '';
in
{
  services.hermes-agent = {
    # Keep this setup native. Container mode writes .container-mode and routes
    # host CLI commands through Docker/Podman, which this setup deliberately avoids.
    container.enable = lib.mkForce false;
    settings = lib.mkForce supportConfig;
    extraPackages = [ pkgs.codex ];
  };

  environment.systemPackages = [
    pkgs.codex
    (profileAlias "architect")
    (profileAlias "builder")
    (profileAlias "review")
    (pkgs.writeShellScriptBin "support" ''
      export HERMES_HOME=${hermesHome}
      exec /run/current-system/sw/bin/hermes "$@"
    '')
  ];

  system.activationScripts.hermes-codex-profiles = lib.stringAfter [ "hermes-agent-setup" ] ''
    install -d -m 2770 -o ${hermesUser} -g ${hermesGroup} "${hermesHome}/profiles"
    ${lib.concatStringsSep "\n" (lib.mapAttrsToList installProfile profiles)}
  '';
}
