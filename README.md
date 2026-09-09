# Flakebuilder

Build a NixOS `flake.nix` from a searchable menu of configuration bits and inputs.

Flakebuilder is a separate descendant of [ConfigBuilder](https://github.com/soltros/configbuilder)'s menu-based workflow. It embeds selected configuration directly into **one `flake.nix`**. It does not generate a `modules/` directory or import a local hardware file. Nix still generates its normal `flake.lock` when you lock or build. Upstream inputs such as nixpkgs, Hermes and Unity remain external dependencies, pinned by that lock.

The initial catalog contains **83 bits and 12 presets**, covering desktops, hardware support, power, storage, services, shell setup, appearance, applications, gaming and external package sources. Ten presets are based on the branches of `soltros/nixos-config`. They are starting points, not certified replacements for those machines.

## Run

With Go 1.24 or newer and Nix installed:

```sh
CGO_ENABLED=0 go build -o flakebuilder .
./flakebuilder --list
./flakebuilder --state-version 26.05 --host workstation --user alice \
  --hardware /etc/nixos/hardware-configuration.nix --dir ./generated
```

Or from this checkout:

```sh
nix run . -- --list
```

Use the installation's **original** `system.stateVersion`; selecting a newer nixpkgs track does not mean changing that value. Hostname, user, locale, timezone, keyboard and personal paths are CLI settings. Defaults use a generic user and UTC. Hardware comes from your own generated hardware configuration, never from a preset's original machine.

Menu controls:

- Arrow keys / j / k: navigate; Space: select a bit.
- `/`: search by name, category or ID; Esc: clear search.
- `i`: add an external input or remove a custom input by name.
- Enter: preview the single generated file; y: generate; n: return to editing.
- q / Ctrl+C: cancel without changing output files.

`[+]` means a bit is needed by another selected bit. Dependency selections are calculated automatically. Remove the dependent selection to remove an automatically selected bit. Conflicting desktop, shell, bootloader, kernel or power choices must be resolved before generation.

## Add an input, overlay, module or package

The input menu asks for a name, URL, usage, output attribute and whether its nixpkgs input should follow the root nixpkgs. It supports these usages:

| Usage | Generated reference |
|---|---|
| input | Input declaration only |
| overlay | `inputs.NAME.overlays.ATTRIBUTE` in `nixpkgs.overlays` |
| module | `inputs.NAME.nixosModules.ATTRIBUTE` in the system's module list |
| package | `inputs.NAME.packages.SYSTEM.ATTRIBUTE` in system packages |

An upstream module may also need configuration options to enable its service; importing it alone does not necessarily enable anything. The built-in Hermes/Unity bits supply those settings. A custom catalog bit can provide settings for other modules.

The same functionality is available without the menu:

```sh
./flakebuilder --preset minimal --state-version 26.05 --host workstation \
  --input soltros_nixpkgs=github:soltros/soltros_nixpkgs \
  --input-follows soltros_nixpkgs \
  --input-use soltros_nixpkgs=overlay:default \
  --hardware /etc/nixos/hardware-configuration.nix --dir ./generated --yes
```

Built-in integrations already wire their input, follows policy and settings:

```sh
./flakebuilder --state-version 26.05 --host workstation --user alice \
  --bits nix-settings,networkmanager,localization,user,boot-systemd,desktop-plasma,waterfox,voxtype,hermes \
  --hardware /etc/nixos/hardware-configuration.nix --dir ./generated --yes
```

For another source, repeat `--input NAME=URL`, `--input-use NAME=package:ATTRIBUTE` / `NAME=module:ATTRIBUTE` / `NAME=overlay:ATTRIBUTE`, and optionally `--input-follows NAME`. `--input-use` can also refer to a known catalog input without redeclaring its URL. Output attributes may be dotted paths. Local input paths are excluded so generated output does not depend on a source checkout.

Explicit declarations can override a built-in input's URL, for example:

```sh
./flakebuilder --state-version 26.05 --bits voxtype \
  --input voxtype=github:peteonrails/voxtype/v0.7.5 \
  --input-follows voxtype --stdout
```

Preserve required follows mappings when overriding catalog inputs. Changing an input version can change compatibility. Use either a built-in bit or a manual consumer for the same overlay to avoid applying it twice.

## Generate, reopen, validate and build

```sh
# View output only; no Nix commands or output files.
./flakebuilder --preset laptop_plasma --state-version 26.05 --stdout

# Reopen a generated file's saved selections, including embedded hardware.
./flakebuilder --from ./generated/flake.nix --dir ./generated --force

# Lock inputs and evaluate in a temporary staging directory, then write output.
./flakebuilder --from ./generated/flake.nix --dir ./generated --yes --force --lock

# Build the selected system before replacing the generated files.
./flakebuilder --from ./generated/flake.nix --dir ./generated --yes --force --build
```

Normal generation checks Nix syntax. `--lock` additionally locks inputs and runs `nix flake check --no-build`. `--build` also runs a non-activating `nix build` of the system toplevel. **Flakebuilder never runs switch, boot or nixos-install.** Build failures leave existing generated files unchanged. Successful replacements keep backups. Existing lock pins are reused when applicable; changing selections may add or remove inputs.

`--from` reads selection metadata, not arbitrary hand edits to the Nix body. Hand edits remain usable by Nix, but reopening and regenerating produces code from the saved selections. Backups preserve the old file. No output is overwritten without `--force`. Generation can omit hardware for previews; building requires it. Standard hardware files using `modulesPath` are supported; local file imports must first be inlined.

## Add more bits

See [catalog authoring](docs/CATALOG.md). The bundled catalog lives in `internal/builder/catalog.json`; each source bit is a complete NixOS module expression under `internal/builder/bits/`. Those source files are embedded into the executable and inlined into generated output. They are not runtime module files required by the generated flake.

A trusted local catalog can be used with `--catalog DIRECTORY`, without rebuilding the executable. Catalogs and hardware files are executable Nix code and should come from a trusted source. The manifest handles transitive requirements, input dependencies, conflicts, exclusive groups, supported platforms and selected nixpkgs tracks. Nix performs the actual option and package validation.

## Current limits

- `--nixpkgs` supports unstable and 26.05. Per-bit track declarations are compatibility constraints, not evidence that every upstream revision builds.
- The branch-inspired presets deliberately avoid embedding anyone's disks/UUIDs. Some optional personal integrations retain original Durandal persona text and Hermes paths. Review those bits before adopting them.
- The Hermes runtime workaround follows source code that assumes Python 3.12. It is optional and may need updating against a new Hermes revision.
- Muse uses a runtime upstream installer; a locked flake does not pin the software it later downloads.
- Input existence, output attributes and upstream module compatibility are checked by Nix during `--lock`/`--build`, not by guessing from a URL.
- The first version uses one desktop and one normal user per generated host. Multiple GPU support bits can be combined where the hardware warrants it.

## Development

```sh
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go vet ./...
```

With `nix-instantiate` available, tests parse every bit and preset as a generated flake. Additional tests cover dependencies, conflicts, cycles, input integration, string escaping, round trips, cancellation, backups and failed builds.

GPL-3.0. See [LICENSE](LICENSE) and [source attribution](docs/ATTRIBUTION.md).
