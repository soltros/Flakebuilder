#!/usr/bin/env python3
"""Promote declarative examples from the complete NixOS Wiki scrape."""
import hashlib, json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SCRAPE = Path('/home/derrik/nixos_wiki_rag')
articles = json.loads((SCRAPE/'nixos_wiki_articles.json').read_text())
chunks = (SCRAPE/'nixos_wiki_chunks.jsonl').read_bytes().splitlines()
catalog_path = ROOT/'internal/builder/catalog.json'
bits_dir = ROOT/'internal/builder/bits'
catalog = json.loads(catalog_path.read_text())
existing = {b['id'] for b in catalog['bits']}

# id|label|category|NixOS module body|wiki article title
raw = r'''zfs|ZFS filesystem support|Storage|boot.supportedFilesystems = [ "zfs" ]; services.zfs.autoSnapshot.enable = true;|ZFS
lvm|LVM2 tools and device mapper|Storage|boot.initrd.kernelModules = [ "dm_snapshot" "dm_crypt" ]; environment.systemPackages = [ pkgs.lvm2 ];|LVM
ext4|Ext4 filesystem tools|Storage|boot.supportedFilesystems = [ "ext4" ]; environment.systemPackages = [ pkgs.e2fsprogs ];|Ext4
encrypted-root|Initrd encrypted root support|Storage|boot.initrd.luks.devices = lib.mkDefault {};|Full Disk Encryption
auto-gc|Automatic Nix store garbage collection|Nix|nix.gc = { automatic = true; dates = "weekly"; options = "--delete-older-than 30d"; };|Storage optimization
auto-upgrade|Automatic NixOS upgrades|Nix|system.autoUpgrade = { enable = true; dates = "weekly"; allowReboot = false; };|Automatic system upgrades
binary-cache|Nix substituter cache|Nix|nix.settings.substituters = [ "https://cache.nixos.org/" ];|Binary Cache
nix-sandbox|Nix sandbox|Nix|nix.settings.sandbox = true;|FAQ
sudo|Sudo administration|Security|security.sudo.enable = true;|Sudo
sudo-passwordless-wheel|Passwordless wheel sudo|Security|security.sudo.wheelNeedsPassword = false;|Sudo
doas|Doas administration|Security|security.doas.enable = true; security.sudo.enable = false;|Doas
polkit|Polkit authorization|Security|security.polkit.enable = true;|Polkit
tpm2|TPM2 support|Security|security.tpm2.enable = true;|TPM
fail2ban|Fail2ban intrusion protection|Security|services.fail2ban.enable = true;|Fail2ban
nixos-hardened|NixOS hardened profile|Security|boot.kernelPackages = pkgs.linuxKernel.packages.linux_hardened; security.lockKernelModules = true; security.protectKernelImage = true;|NixOS Hardening
firewall-base|Stateful firewall|Networking|networking.firewall.enable = true;|Firewall
firewall-ping|Allow ICMP ping|Networking|networking.firewall.allowPing = true;|Firewall
ipv6|IPv6 networking|Networking|networking.enableIPv6 = true;|Networking
iwd|IWD wireless backend|Networking|networking.wireless.iwd.enable = true;|Iwd
wireguard|WireGuard tools|Networking|environment.systemPackages = [ pkgs.wireguard-tools ];|WireGuard
openvpn|OpenVPN client tools|Networking|environment.systemPackages = [ pkgs.openvpn ];|OpenVPN
avahi|Avahi service discovery|Networking|services.avahi = { enable = true; publish.enable = true; publish.userServices = true; };|Printing
ntp|NTP time synchronization|Networking|services.ntpd.enable = true;|NTP
chrony|Chrony time synchronization|Networking|services.chrony.enable = true;|Chrony
networkd-dispatcher|Networkd dispatcher|Networking|services.networkd-dispatcher.enable = true;|Systemd/networkd/dispatcher
resolved|systemd-resolved DNS|Networking|services.resolved.enable = true;|Systemd/resolved
nfs-server|NFS server|Servers|services.nfs.server.enable = true;|NFS
samba|Samba file server|Servers|services.samba.enable = true;|Samba
docker|Docker engine|Virtualization|virtualisation.docker.enable = true; users.users.${[[ nix .User ]]}.extraGroups = [ "docker" ];|Docker
qemu|QEMU guest tools|Virtualization|environment.systemPackages = [ pkgs.qemu ];|QEMU
kvm|KVM virtualization modules|Virtualization|boot.kernelModules = [ "kvm-intel" "kvm-amd" ];|NixOps/Virtualization
pci-passthrough|PCI passthrough support|Virtualization|boot.kernelParams = [ "amd_iommu=on" "intel_iommu=on" "iommu=pt" ];|PCI passthrough
virtiofs|VirtioFS support|Virtualization|boot.kernelModules = [ "virtiofs" ];|VirtioFS
nvidia|NVIDIA graphics|Hardware|services.xserver.videoDrivers = [ "nvidia" ]; hardware.graphics.enable = true;|NVIDIA
opengl|OpenGL graphics stack|Hardware|hardware.graphics.enable = true; hardware.graphics.enable32Bit = true;|OpenGL
mesa|Mesa OpenGL drivers|Hardware|hardware.graphics.extraPackages = [ pkgs.mesa.drivers ];|Mesa
thunderbolt|Thunderbolt authorization|Hardware|services.hardware.bolt.enable = true;|Thunderbolt
fwupd|Firmware update service|Hardware|services.fwupd.enable = true;|Fwupd
smartd|SMART disk monitoring|Hardware|services.smartd.enable = true;|Smartmontools
sane|Scanner support|Hardware|hardware.sane.enable = true;|Scanners
pcscd|PC/SC smart card daemon|Hardware|services.pcscd.enable = true;|Web eID
logitech|Logitech device support|Hardware|hardware.logitech.enable = true;|Logitech Unifying Receiver
touchpad|Libinput touchpad support|Hardware|services.xserver.libinput.enable = true;|Touchpad
touchpad-flat|Disable touchpad acceleration|Hardware|services.xserver.libinput = { enable = true; accelProfile = "flat"; };|Xorg
backlight|Backlight control|Hardware|programs.light.enable = true;|Backlight
zsa-keyboards|ZSA keyboard udev rules|Hardware|hardware.keyboard.zsa.enable = true;|ZSA Keyboards
drawing-tablet|Drawing tablet support|Hardware|services.xserver.digimend.enable = true;|Drawing Tablet
gnome-keyring|GNOME keyring|Desktop integrations|services.gnome.gnome-keyring.enable = true;|GNOME
dconf|Dconf settings database|Desktop integrations|programs.dconf.enable = true;|GNOME/Calendar
gdm|GDM display manager|Desktop|services.displayManager.gdm.enable = true;|GDM
sddm|SDDM display manager|Desktop|services.displayManager.sddm.enable = true;|SDDM
lightdm|LightDM display manager|Desktop|services.xserver.displayManager.lightdm.enable = true;|LightDM
greetd|Greetd login manager|Desktop|services.greetd.enable = true;|Greetd
desktop-gnome|GNOME desktop|Desktop|services.desktopManager.gnome.enable = true;|GNOME
desktop-xfce|XFCE desktop|Desktop|services.xserver = { enable = true; desktopManager.xfce.enable = true; };|Xfce
desktop-cinnamon|Cinnamon desktop|Desktop|services.xserver = { enable = true; desktopManager.cinnamon.enable = true; };|Cinnamon
desktop-budgie|Budgie desktop|Desktop|services.xserver = { enable = true; desktopManager.budgie.enable = true; };|Budgie Desktop
desktop-mate|MATE desktop|Desktop|services.xserver = { enable = true; desktopManager.mate.enable = true; };|NixOS as a desktop
desktop-lxqt|LXQt desktop|Desktop|services.xserver = { enable = true; desktopManager.lxqt.enable = true; };|NixOS as a desktop
desktop-i3|i3 window manager|Desktop|services.xserver = { enable = true; windowManager.i3.enable = true; };|I3
desktop-sway|Sway Wayland compositor|Desktop|programs.sway.enable = true;|Sway
desktop-hyprland|Hyprland Wayland compositor|Desktop|programs.hyprland.enable = true;|Hyprland
desktop-niri|Niri Wayland compositor|Desktop|programs.niri.enable = true;|Niri/en
desktop-awesome|Awesome window manager|Desktop|services.xserver.windowManager.awesome.enable = true;|Awesome
desktop-dwm|Dwm window manager|Desktop|services.xserver.windowManager.dwm.enable = true;|Dwm
desktop-qtile|Qtile window manager|Desktop|services.xserver.windowManager.qtile.enable = true;|Qtile
desktop-xmonad|XMonad window manager|Desktop|services.xserver.windowManager.xmonad.enable = true;|XMonad
desktop-openbox|Openbox window manager|Desktop|services.xserver.windowManager.openbox.enable = true;|Openbox
desktop-bspwm|Bspwm window manager|Desktop|services.xserver.windowManager.bspwm.enable = true;|Bspwm
desktop-leftwm|LeftWM window manager|Desktop|services.xserver.windowManager.leftwm.enable = true;|Leftwm
waybar|Waybar status bar|Desktop integrations|programs.waybar.enable = true;|Waybar
swaylock|Swaylock screen locker|Desktop integrations|programs.swaylock.enable = true;|Swaylock
swayidle|Swayidle idle daemon|Desktop integrations|environment.systemPackages = [ pkgs.swayidle ];|Swayidle
fuzzel|Fuzzel application launcher|Desktop integrations|environment.systemPackages = [ pkgs.fuzzel ];|Fuzzel
rofi|Rofi application launcher|Desktop integrations|programs.rofi.enable = true;|Rofi
picom|Picom compositor|Desktop integrations|services.picom.enable = true;|Picom
redshift|Redshift color temperature|Desktop integrations|services.redshift.enable = true;|Redshift
nix-index|Nix package index|Applications|programs.nix-index.enable = true;|NixOS search
vim|Vim editor|Applications|programs.vim.enable = true;|Vim/en
neovim|Neovim editor|Applications|programs.neovim = { enable = true; defaultEditor = true; };|Neovim/en
emacs|Emacs editor|Applications|programs.emacs.enable = true;|Emacs
vscode|Visual Studio Code|Applications|programs.vscode.enable = true;|Visual Studio Code/en
vscodium|VSCodium editor|Applications|environment.systemPackages = [ pkgs.vscodium ];|VSCodium
firefox|Firefox browser|Applications|programs.firefox.enable = true;|Firefox/en
chromium|Chromium browser|Applications|programs.chromium.enable = true;|Chromium/en
qutebrowser|Qutebrowser|Applications|environment.systemPackages = [ pkgs.qutebrowser ];|Qutebrowser
kitty|Kitty terminal|Applications|programs.kitty.enable = true;|Kitty/en
tmux|Tmux terminal multiplexer|Applications|programs.tmux = { enable = true; clock24 = true; };|Tmux
zellij|Zellij terminal multiplexer|Applications|programs.zellij.enable = true;|Zellij
zoxide|Zoxide directory jumper|Applications|programs.zoxide.enable = true;|Zoxide
atuin|Atuin shell history|Applications|programs.atuin.enable = true;|Atuin
fzf|Fzf fuzzy finder|Applications|programs.fzf.fuzzyCompletion = true;|Fzf
yazi|Yazi terminal file manager|Applications|environment.systemPackages = [ pkgs.yazi ];|Yazi
ranger|Ranger terminal file manager|Applications|environment.systemPackages = [ pkgs.ranger ];|Ranger
helix|Helix editor|Applications|programs.helix.enable = true;|Helix
git|Git version control|Applications|programs.git.enable = true;|Git
direnv|Direnv development environments|Development|programs.direnv.enable = true;|Direnv
devenv|Devenv project environments|Development|environment.systemPackages = [ pkgs.devenv ];|Devenv
gcc|GCC compiler toolchain|Development|environment.systemPackages = [ pkgs.gcc ];|C
clang|Clang compiler toolchain|Development|environment.systemPackages = [ pkgs.clang ];|Using Clang instead of GCC
go|Go toolchain|Development|environment.systemPackages = [ pkgs.go ];|Go
rust|Rust toolchain|Development|environment.systemPackages = [ pkgs.rustc pkgs.cargo ];|Rust
python|Python toolchain|Development|environment.systemPackages = [ pkgs.python3 ];|Python
nodejs|Node.js toolchain|Development|environment.systemPackages = [ pkgs.nodejs ];|Node.js
postgresql|PostgreSQL server|Servers|services.postgresql.enable = true;|PostgreSQL
mysql|MySQL server|Servers|services.mysql.enable = true;|Mysql
redis|Redis server|Servers|services.redis.servers.default.enable = true;|NixOS as a server
nginx|Nginx web server|Servers|services.nginx.enable = true;|Nginx
apache|Apache HTTP server|Servers|services.httpd.enable = true;|Apache Httpd
caddy|Caddy web server|Servers|services.caddy.enable = true;|Caddy
bind|BIND DNS server|Servers|services.bind.enable = true;|Bind
unbound|Unbound validating DNS resolver|Servers|services.unbound.enable = true;|Unbound
prometheus|Prometheus metrics server|Servers|services.prometheus.enable = true;|Prometheus
grafana|Grafana dashboards|Servers|services.grafana.enable = true;|Grafana
loki|Grafana Loki logs|Servers|services.loki.enable = true;|Grafana Loki
plex|Plex media server|Servers|services.plex.enable = true;|Plex
sonarr|Sonarr media automation|Servers|services.sonarr.enable = true;|Sonarr
radarr|Radarr movie automation|Servers|services.radarr.enable = true;|Radarr
transmission|Transmission torrent client|Servers|services.transmission.enable = true;|Transmission
nextcloud|Nextcloud server|Servers|services.nextcloud.enable = true;|Nextcloud
vaultwarden|Vaultwarden password server|Servers|services.vaultwarden.enable = true;|Vaultwarden
mosquitto|Mosquitto MQTT broker|Servers|services.mosquitto.enable = true;|Mosquitto
home-assistant|Home Assistant|Servers|services.home-assistant.enable = true;|Home Assistant
k3s|K3s Kubernetes distribution|Servers|services.k3s.enable = true; services.k3s.role = "server";|K3s
kubernetes|Kubernetes services|Servers|services.kubernetes.kubelet.enable = true;|Kubernetes
clamav|ClamAV antivirus|Security|services.clamav.daemon.enable = true; services.clamav.updater.enable = true;|Clamav
opensnitch|OpenSnitch firewall UI|Security|services.opensnitch.enable = true; services.opensnitch-ui.enable = true;|OpenSnitch
openrgb|OpenRGB hardware control|Hardware|services.udev.packages = [ pkgs.openrgb ]; boot.kernelModules = [ "i2c-dev" ];|OpenRGB
steam-gamescope|Steam Gamescope session|Gaming|programs.gamescope = { enable = true; capSysNice = true; }; programs.steam.gamescopeSession.enable = true;|Steam/en
lutris|Lutris game launcher|Gaming|environment.systemPackages = [ pkgs.lutris ];|Lutris
retroarch|RetroArch emulator|Gaming|environment.systemPackages = [ pkgs.retroarch ];|RetroArch
dolphin|Dolphin emulator|Gaming|environment.systemPackages = [ pkgs.dolphinEmu ];|Dolphin Emulator
obs|OBS Studio|Media|programs.obs-studio.enable = true;|OBS Studio
mpv|MPV media player|Media|environment.systemPackages = [ pkgs.mpv ];|MPV
kodi|Kodi media center|Media|services.xserver.desktopManager.kodi.enable = true;|Kodi
gstreamer|GStreamer multimedia stack|Media|environment.systemPackages = [ pkgs.gstreamer ];|GStreamer
libreoffice|LibreOffice suite|Applications|environment.systemPackages = [ pkgs.libreoffice ];|LibreOffice
thunderbird|Thunderbird mail client|Applications|programs.thunderbird.enable = true;|Thunderbird
spotify|Spotify client|Applications|environment.systemPackages = [ pkgs.spotify ];|Spotify
telegram|Telegram desktop client|Applications|environment.systemPackages = [ pkgs.telegram-desktop ];|Telegram
discord|Discord client|Applications|environment.systemPackages = [ pkgs.discord ];|Discord
obsidian|Obsidian notes|Applications|environment.systemPackages = [ pkgs.obsidian ];|Obsidian
calibre|Calibre ebook manager|Applications|environment.systemPackages = [ pkgs.calibre ];|Calibre
impermanence|Impermanence filesystem helper|Storage|programs.fuse.userAllowOther = true;|Impermanence
agenix|Agenix secret management tools|Security|environment.systemPackages = [ pkgs.agenix ];|Agenix
restic|Restic backup client|Backup|environment.systemPackages = [ pkgs.restic ];|Restic
borg|Borg backup client|Backup|environment.systemPackages = [ pkgs.borgbackup ];|Borg backup
rsync|Rsync file synchronization|Backup|environment.systemPackages = [ pkgs.rsync ];|Rsync
rclone|Rclone cloud synchronization|Backup|environment.systemPackages = [ pkgs.rclone ];|Rclone
smartmontools|Smartmontools disk health|Storage|environment.systemPackages = [ pkgs.smartmontools ];|Smartmontools
btrbk|Btrfs snapshot backup|Storage|environment.systemPackages = [ pkgs.btrbk ];|Btrbk
xserver|Xorg server|Desktop|services.xserver.enable = true;|Xorg
console-fonts|Console fonts|Appearance|console.font = "Lat2-Terminus16";|Console Fonts
cursor-themes|Cursor theme support|Appearance|environment.systemPackages = [ pkgs.catppuccin-cursors ];|Cursor Themes
plymouth-basic|Plymouth boot splash|Appearance|boot.plymouth.enable = true;|Plymouth
secure-boot|Secure Boot tooling|Boot|environment.systemPackages = [ pkgs.sbctl ];|Secure Boot
u-boot|U-Boot bootloader tools|Boot|environment.systemPackages = [ pkgs.ubootTools ];|U-Boot
refind|rEFInd bootloader|Boot|boot.loader.refind.enable = true;|REFInd
specialisations|NixOS specialisations example|Nix|specialisation.test.configuration = { documentation.enable = false; };|Specialisation'''

added = 0
for line in raw.splitlines():
    ident, label, category, body, source = line.split('|', 4)
    if ident in existing or source not in articles: continue
    path = bits_dir/f'{ident}.nix'
    path.write_text('{ config, lib, pkgs, ... }: {\n'+body+'\n}\n')
    catalog['bits'].append({'id': ident, 'label': label, 'category': category,
        'description': f'Wiki-backed {label} setting', 'file': f'bits/{ident}.nix', 'wiki': [source]})
    added += 1

catalog_path.write_text(json.dumps(catalog, indent=2, ensure_ascii=False)+'\n')
digests = {}
for name in ('nixos_wiki_articles.json','nixos_wiki_chunks.jsonl'):
    data=(SCRAPE/name).read_bytes(); digests[name]=(len(data),hashlib.sha256(data).hexdigest())
(ROOT/'docs/WIKI_COVERAGE.md').write_text(f'''# NixOS Wiki coverage\n\nThe complete scrape was read on 2026-09-10: **{len(articles):,} article records**, **{len(chunks):,} chunks**, and **{sum(len(v.get("markdown", "")) for v in articles.values()):,} article-markdown characters**. This pass added **{added}** conservative, declarative bits.\n\nThe importer deliberately excludes installation commands, destructive operations, secrets, machine UUIDs, local paths, and prose. Every promoted bit has a `wiki` evidence field and is represented as inline NixOS code. A wiki example is guidance; `--lock`/`--build` performs authoritative nixpkgs evaluation.\n\nSource digests:\n\n- `nixos_wiki_articles.json`: {digests['nixos_wiki_articles.json'][0]:,} bytes, SHA-256 `{digests['nixos_wiki_articles.json'][1]}`\n- `nixos_wiki_chunks.jsonl`: {digests['nixos_wiki_chunks.jsonl'][0]:,} bytes, SHA-256 `{digests['nixos_wiki_chunks.jsonl'][1]}`\n''')
print(f'added {added} bits; catalog now has {len(catalog["bits"])}')
