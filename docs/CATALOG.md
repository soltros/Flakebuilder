# Catalog authoring

A catalog is a JSON manifest plus trusted Nix template files. Version 1 example:

```json
{
  "version": 1,
  "inputs": {
    "example": {
      "url": "github:owner/repo",
      "follows": { "nixpkgs": "nixpkgs" }
    }
  },
  "bits": [
    {
      "id": "example-service",
      "label": "Example service",
      "category": "Services",
      "description": "Enable the upstream service",
      "file": "bits/example-service.nix",
      "inputs": ["example"],
      "requires": [],
      "conflicts": [],
      "systems": ["x86_64-linux"],
      "tracks": ["unstable"]
    }
  ],
  "presets": {}
}
```

Example `bits/example-service.nix`:

```nix
{ inputs, ... }: {
  imports = [ inputs.example.nixosModules.default ];
  services.example.enable = true;
}
```

The example service name is illustrative. Use the real upstream module options.

A bit is a **complete NixOS module expression**, including its function header and any `let` bindings. The renderer places each expression in parentheses in the generated flake's `modules` list. Nix merges option definitions, package lists and service settings through its module system. Do not concatenate arbitrary attribute-set bodies or use last-write-wins JSON merging for Nix options.

- `requires` adds other bits transitively, once each.
- `conflicts` rejects another selected bit, including transitive selections.
- `group` implements an exclusive choice such as desktop or shell.
- `inputs` refers to declarations in the input registry; declarations are emitted only when needed (or explicitly added by the user).
- `tracks` and `systems` restrict selection. Empty means unrestricted by the catalog, not universally verified.
- Presets contain `label`, `track`, and `bits`.

Use `[[ nix .User ]]`, `[[ nix .Timezone ]]`, `[[ nix .Workspace ]]`, etc. to insert a correctly escaped Nix string from the selection model. `[[ nix (home) ]]` renders the selected user's home directory. Never insert unescaped user text into Nix syntax. Go templates use `[[` / `]]` so Nix `${...}` is left intact.

All selected configuration must live in the generated file. Do not introduce `./module.nix`, `builtins.readFile /outside/file`, local asset paths, or `<channel-lookups>`. Embed assets with Nix strings or use pinned upstream fetches. The generated hardware file's `modulesPath` import is supplied by nixpkgs itself.

External-input adapters may use `inputs`, supplied via `specialArgs`. Ordinary bits should prefer `{ config, lib, pkgs, ... }`. Inputs are static top-level declarations: do not try to add flake inputs from a module's `config` or make `imports` depend on the module fixpoint.

Source catalogs are trusted code. The loader checks IDs, references and file paths; it is not a Nix sandbox or a full static dependency analyzer. Generated source is saved before syntax checks, so validation failures do not prevent export. Input fetching and module evaluation occur with --lock or --build.

The user-input menu can declare a source and consume exported overlays/modules/packages without authoring a bit. Add a catalog adapter when a service also needs enable flags, user settings, dependencies or compatibility constraints.
