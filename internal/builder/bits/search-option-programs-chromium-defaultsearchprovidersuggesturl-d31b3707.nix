{ pkgs, config, ... }: { programs.chromium.defaultSearchProviderSuggestURL = "https://encrypted.google.com/complete/search?output=chrome&q={searchTerms}"; }
