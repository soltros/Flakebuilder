#!/usr/bin/env python3
"""Build deterministic GitHub Wiki Markdown from the Flakebuilder catalog."""
import json
import re
import shutil
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
OUT = Path('/tmp/flakebuilder-wiki')
CATALOG = json.loads((ROOT / 'internal/builder/catalog.json').read_text())
BITS = sorted(CATALOG['bits'], key=lambda b: (b.get('category', ''), b['id']))
BY_ID = {b['id']: b for b in BITS}

def esc(value):
    return str(value or '').replace('|', '\\|').replace('\n', ' ')

def slug(value):
    return re.sub(r'[^A-Za-z0-9]+', '-', value).strip('-')

def refs(bit):
    wiki = bit.get('wiki')
    if not wiki:
        return 'catalog'
    if isinstance(wiki, list):
        return ', '.join(f'[{esc(title)}](https://wiki.nixos.org/wiki/{title.replace(" ", "_")})' for title in wiki)
    title = wiki.get('title', 'NixOS Wiki')
    url = wiki.get('url', '')
    return f'[{esc(title)}]({url})' if url else esc(title)

def bit_row(bit, link=False):
    ident = f"[{bit['id']}](Bits#{bit['id']})" if link else bit['id']
    req = ', '.join(bit.get('requires', [])) or '—'
    inp = ', '.join(bit.get('inputs', [])) or '—'
    return f"| `{ident}` | {esc(bit.get('label'))} | {esc(bit.get('description'))} | {esc(req)} | {esc(inp)} | {refs(bit)} |"

if OUT.exists():
    shutil.rmtree(OUT)
OUT.mkdir(parents=True)

categories = {}
for bit in BITS:
    categories.setdefault(bit.get('category', 'Other'), []).append(bit)

counts = len(BITS)
home = f"""# Flakebuilder

Flakebuilder builds a complete, single-file `flake.nix` from selectable configuration bits. It is a menu-driven derivation of [Configbuilder](https://github.com/soltros/configbuilder): choose a preset, add or remove bits, add inputs and overlays, then generate a flake you own.

The catalog currently contains **{counts} bits**, **{len(CATALOG['presets'])} presets**, and **{len(CATALOG['inputs'])} built-in inputs**. Generated flakes are saved in `~/generated_flakes` by default.

## Quick start

```bash
git clone https://github.com/soltros/Flakebuilder.git
cd Flakebuilder
go run .
```

Use the interactive menu to select a preset and bits. The command line can generate non-interactively:

```bash
go run . --preset minimal --bits ssh,pipewire --output ~/generated_flakes/my-flake
```

Generation always writes the flake first. Optional parse, lock, evaluation, and build checks can then be run; a failed check never removes the generated file.

## Wiki navigation

* [Complete bit catalog](Bits)
* [Bits by category](Categories)
* [Presets](Presets)
* [Inputs and overlays](Inputs)
* [Building and validation](Safety-and-Building)
* [Catalog authoring](Catalog-Authoring)
* [NixOS Wiki coverage](Wiki-Source-Coverage)

Flakebuilder emits ordinary Nix code with no runtime dependency on external modules. See the repository [README](https://github.com/soltros/Flakebuilder#readme) for command options and development details.
"""
(OUT / 'Home.md').write_text(home)

rows = '\n'.join(bit_row(b, True) for b in BITS)
(OUT / 'Bits.md').write_text(f"""# Complete bit catalog

Every selectable configuration bit is listed below. Bits are self-contained snippets assembled into the generated flake.

| ID | Label | Description | Requires | Inputs | Evidence |
|---|---|---|---|---|---|
{rows}
""")

cat_links = '\n'.join(f"* [{name} ({len(items)})](Category-{slug(name)})" for name, items in sorted(categories.items()))
(OUT / 'Categories.md').write_text(f"# Bits by category\n\n{cat_links}\n")
for name, items in sorted(categories.items()):
    table = '\n'.join(bit_row(b) for b in sorted(items, key=lambda b: b['id']))
    (OUT / f'Category-{slug(name)}.md').write_text(f"""# {name}

{len(items)} bits are in this category.

| ID | Label | Description | Requires | Inputs | Evidence |
|---|---|---|---|---|---|
{table}
""")

p_rows = []
for name, preset in sorted(CATALOG['presets'].items()):
    bits = preset.get('bits', [])
    p_rows.append(f"| `{name}` | {esc(preset.get('label'))} | `{preset.get('track', 'unstable')}` | {len(bits)} | {', '.join(f'`{b}`' for b in bits)} |")
(OUT / 'Presets.md').write_text("""# Presets

Presets are starting points. Every preset can be edited in the menu or overridden with `--bits`; hardware-specific values remain yours to supply.

| Name | Description | Track | Bits | Included bit IDs |
|---|---|---:|---:|---|
""" + '\n'.join(p_rows) + '\n')

i_rows = []
for name, inp in sorted(CATALOG['inputs'].items()):
    follows = ', '.join(f'`{k}` → `{v}`' for k, v in inp.get('follows', {}).items()) or '—'
    i_rows.append(f"| `{name}` | `{esc(inp.get('url'))}` | {follows} |")
(OUT / 'Inputs.md').write_text("""# Inputs and overlays

Built-in inputs are available from the input menu and can also be selected with `--inputs`. The generated flake declares them in its `inputs` block and passes them to outputs. Add arbitrary repositories with repeated `--input name=url` arguments.

| Name | URL | Follows |
|---|---|---|
""" + '\n'.join(i_rows) + """

## Custom declarations

Use `--input NAME=URL` for an extra flake input and `--overlay PATH` for a Nix overlay file or expression. Bits may declare required inputs; Flakebuilder reports missing selections before writing the final output.
""")

(OUT / 'Safety-and-Building.md').write_text("""# Generation, validation, and building

Flakebuilder generates a normal `flake.nix` in the selected output directory. By default this is `~/generated_flakes`; existing directories are backed up automatically, while `--force` permits replacement of a custom destination.

Generation is independent from validation. The file is saved before any optional checks, so parse errors, lock failures, evaluation failures, or build failures do not prevent you from inspecting and editing the result.

Available checks include Nix syntax parsing, lock-file generation, `nix flake check`, and optional builds. Flakebuilder does not activate or switch the resulting system. Review the generated code and run your preferred `nixos-rebuild` command yourself.
""")

(OUT / 'Catalog-Authoring.md').write_text("""# Catalog authoring

The catalog lives at [`internal/builder/catalog.json`](https://github.com/soltros/Flakebuilder/blob/main/internal/builder/catalog.json). A bit has an ID, label, category, description, and a Nix snippet under `bits/`. Optional `requires`, `conflicts`, and `inputs` fields express composition rules. Presets map names to tracks and ordered bit IDs.

Keep snippets declarative and composable. If a bit depends on an external flake input, list its catalog input key in `inputs`. Run `go test ./...` after edits. The importer at `tools/import_wiki_bits.py` can add conservative bits from the scraped NixOS Wiki corpus; review generated entries before committing.
""")

coverage = ROOT / 'docs/WIKI_COVERAGE.md'
(OUT / 'Wiki-Source-Coverage.md').write_text(coverage.read_text() if coverage.exists() else '# NixOS Wiki coverage\n')
(OUT / '_Sidebar.md').write_text("""* [Home](Home)
* [Complete bit catalog](Bits)
* [Bits by category](Categories)
* [Presets](Presets)
* [Inputs and overlays](Inputs)
* [Generation and building](Safety-and-Building)
* [Catalog authoring](Catalog-Authoring)
* [Wiki source coverage](Wiki-Source-Coverage)
""")
print(f'generated {len(list(OUT.glob("*.md")))} wiki pages in {OUT}')
