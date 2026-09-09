{ config, lib, pkgs, ... }: {
time.timeZone = [[ nix .Timezone ]];
i18n.defaultLocale = [[ nix .Locale ]];
services.xserver.xkb.layout = [[ nix .Keyboard ]];
}
