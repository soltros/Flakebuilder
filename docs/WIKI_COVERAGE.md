# NixOS Wiki coverage

The complete scrape was read on 2026-09-10: **1,128 article records**, **5,578 chunks**, and **2,888,930 article-markdown characters**. This pass added **156** conservative, declarative bits.

The importer deliberately excludes installation commands, destructive operations, secrets, machine UUIDs, local paths, and prose. Every promoted bit has a `wiki` evidence field and is represented as inline NixOS code. A wiki example is guidance; `--lock`/`--build` performs authoritative nixpkgs evaluation.

Source digests:

- `nixos_wiki_articles.json`: 3,177,474 bytes, SHA-256 `7cb07d51779919e950b9773a7f043235d55d25095e101e84c8163b6df8b4e7cb`
- `nixos_wiki_chunks.jsonl`: 4,179,410 bytes, SHA-256 `5ebd5c9f9e5db4196eb60a4bd89cdece5db5c01f54eb4091d9261c889144092e`
