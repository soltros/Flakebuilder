# Source attribution

Flakebuilder derives its purpose and menu-based NixOS configuration workflow from
[soltros/configbuilder](https://github.com/soltros/configbuilder), inspected at
`df0b9e6dd0495c2ab39e18620865760525e25b9c`. ConfigBuilder's GPL-3.0 license is retained.
The application and resolver were implemented separately so ConfigBuilder's
channel-based behavior remains independent.

The initial catalog is adapted from the user's `soltros/nixos-config` branches,
reviewed on 2026-09-09:

| Branch | Commit |
|---|---|
| master | 53e4f82dca0bb9a13797b81a8791dced62a34be5 |
| desktop_pantheon | b0dd3a50a092b1588a8f0a3a3658db4d00951f55 |
| desktop_lomiri | f5bf1d376d135ca88104972f635d71159b55ab3b |
| desktop_plasma | 899fe1896a33e78e2d04dd60ff4d98781de4a0ed |
| desktop_unity_testing | c652dffd8ece2891bb0b1bcf70bb48fe2767bbfe |
| laptop | 6734805af1b75752156f889b2e7fba850afc00b0 |
| laptop_pantheon | d97941257e75a4a10fe9d8c6d62b1eb2331b5529 |
| laptop_plasma | c1dd4a0b560cf96795158665a7046571995612e7 |
| laptop_unity_testing | 9296786369caefbc68938627a78751caa060c578 |
| mediacenter_plasma | fac5499c47a2bbcac10b05ffb12b95d0276114c5 |

Presets are inspired by these branches, not exact migrations. Hardware UUIDs and
machine-specific hardware files are deliberately not bundled. Desktop appearance,
Hermes compatibility/profiles, Durandal assets and COSMIC/Plymouth theme sources
retain substantial code or data from that repository. External dependencies retain
their respective licenses.
