{ pkgs, inputs, ... }:
{
  # Hermes tool installers sometimes fetch generic Linux binaries (uv-managed
  # Python tools, agent-browser/Chromium, cua-driver). NixOS needs nix-ld for
  # those ELF interpreters and envfs for hard-coded shebangs like /bin/bash.
  # Sources: https://wiki.nixos.org/wiki/Nix-ld and https://wiki.nixos.org/wiki/Playwright
  services.envfs.enable = true;

  programs.nix-ld = {
    enable = true;
    libraries = with pkgs; [
      acl
      alsa-lib
      at-spi2-atk
      at-spi2-core
      atk
      attr
      bzip2
      cairo
      coreutils
      cups
      curl
      dbus
      dbus-glib
      expat
      ffmpeg
      fontconfig
      freetype
      gdk-pixbuf
      glib
      gtk3
      icu
      libdrm
      libGL
      libgbm
      libnotify
      libpulseaudio
      libva
      libxcrypt
      libxkbcommon
      nspr
      nss
      openssl
      pango
      pipewire
      stdenv.cc.cc
      systemd
      util-linux
      vulkan-loader
      xz
      zlib
      zstd
      libice
      libsm
      libx11
      libxscrnsaver
      libxcomposite
      libxcursor
      libxdamage
      libxext
      libxfixes
      libxi
      libxinerama
      libxrandr
      libxrender
      libxtst
      libxcb
      libxshmfence
    ];
  };

  environment.systemPackages = with pkgs; [
    chromium
    codex
    uv
  ];

  services.hermes-agent = {
    package = (inputs.hermes-agent.packages.${pkgs.stdenv.hostPlatform.system}.default.overrideAttrs (old: {
      postInstall = (old.postInstall or "") + ''
        mkdir -p $out/lib/python3.12/site-packages
        cp ${inputs.hermes-agent}/hermes_state_registry.py $out/lib/python3.12/site-packages/
        cp ${inputs.hermes-agent}/mini_swe_runner.py $out/lib/python3.12/site-packages/
        if [ -f "${inputs.hermes-agent}/hermes_state_holders.py" ]; then
          cp ${inputs.hermes-agent}/hermes_state_holders.py $out/lib/python3.12/site-packages/
        fi
        for bin in $out/bin/*; do
          wrapProgram "$bin" --prefix PYTHONPATH : "$out/lib/python3.12/site-packages"
        done
      '';
    }));
    extraPackages = with pkgs; [
      bash
      chromium
      codex
      curl
      nodejs
      uv
    ];
    environment = {
      AGENT_BROWSER_ENGINE = "chrome";
      PLAYWRIGHT_SKIP_VALIDATE_HOST_REQUIREMENTS = "true";
    };
  };
}
