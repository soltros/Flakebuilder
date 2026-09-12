using Gtk;
using Granite;

public class FlakebuilderGui : Gtk.Application {
    private Gtk.Entry state_version;
    public FlakebuilderGui () {
        Object (application_id: "io.github.soltros.Flakebuilder", flags: ApplicationFlags.DEFAULT_FLAGS);
    }

    protected override void activate () {
        install_style (); 
        var window = new Gtk.ApplicationWindow (this);
        window.title = "Flakebuilder";
        window.default_width = 820;
        window.default_height = 600;
        window.resizable = true;

        var header = new Gtk.HeaderBar ();
        header.show_title_buttons = true;
        var title = new Gtk.Label ("Flakebuilder");
        title.add_css_class ("title-4");
        header.title_widget = title;
        window.titlebar = header;

        var content = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
        content.add_css_class ("page");
        content.margin_start = 48;
        content.margin_end = 48;
        content.margin_top = 40;
        content.margin_bottom = 32;
        content.spacing = 20;

        var hero = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 18);
        hero.add_css_class ("hero");
        var logo = new Gtk.Image.from_icon_name ("applications-system-symbolic");
        logo.pixel_size = 46;
        logo.add_css_class ("hero-icon");
        hero.append (logo);
        var hero_text = new Gtk.Box (Gtk.Orientation.VERTICAL, 4);
        var eyebrow = new Gtk.Label ("NIXOS CONFIGURATION BUILDER");
        eyebrow.halign = Gtk.Align.START;
        eyebrow.add_css_class ("eyebrow");
        hero_text.append (eyebrow);
        var welcome = new Gtk.Label ("Build your NixOS flake");
        welcome.add_css_class ("title-1");
        welcome.halign = Gtk.Align.START;
        hero_text.append (welcome);
        var intro = new Gtk.Label ("Assemble configuration bits, packages and inputs into one reviewable flake.nix.");
        intro.wrap = true;
        intro.halign = Gtk.Align.START;
        intro.add_css_class ("dim-label");
        hero_text.append (intro);
        hero.append (hero_text);
        content.append (hero);

        var card = new Gtk.Box (Gtk.Orientation.VERTICAL, 12);
        card.add_css_class ("card");
        var card_title = new Gtk.Label ("Start a configuration session");
        card_title.halign = Gtk.Align.START;
        card_title.add_css_class ("title-4");
        card.append (card_title);
        card.append (new Gtk.Label ("The selector opens in a terminal so it can provide the complete interactive workflow, including NUR browsing, package shopping, review and dry-run checks.") { wrap = true, halign = Gtk.Align.START });
        var state_row = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 10);
        state_row.margin_top = 8;
        var state_label = new Gtk.Label ("Original NixOS release");
        state_label.halign = Gtk.Align.START;
        state_label.hexpand = true;
        state_row.append (state_label);
        state_version = new Gtk.Entry ();
        state_version.placeholder_text = "e.g. 26.05";
        state_version.width_chars = 10;
        state_version.max_width_chars = 10;
        state_version.add_css_class ("release-entry");
        state_row.append (state_version);
        card.append (state_row);
        content.append (card);

        var actions = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 12);
        var configure = new Gtk.Button.with_label ("Configure flake");
        configure.add_css_class ("suggested-action");
        configure.add_css_class ("pill-button");
        configure.clicked.connect (() => launch_builder ());
        actions.append (configure);
        var open = new Gtk.Button.with_label ("Open generated flakes");
        open.add_css_class ("pill-button");
        open.clicked.connect (() => open_output ());
        actions.append (open);
        content.append (actions);

        var status = new Gtk.Label ("  Output  ·  ~/generated_flakes");
        status.halign = Gtk.Align.START;
        status.add_css_class ("status-pill");
        content.append (status);
        var footer = new Gtk.Label ("Generated files are saved locally. Flakebuilder never activates a system automatically.");
        footer.halign = Gtk.Align.START;
        footer.add_css_class ("dim-label");
        footer.margin_top = 4;
        content.append (footer);
        window.child = content;
        window.present ();
    }

    private void install_style () {
        var css = new Gtk.CssProvider ();
        css.load_from_data ("""
            .page { background: @theme_bg_color; }
            .hero { padding-bottom: 4px; }
            .hero-icon { color: @accent_color; }
            .eyebrow { font-size: 11px; font-weight: 700; letter-spacing: 1px; color: @accent_color; }
            .card { padding: 24px; border-radius: 16px; border: 1px solid alpha(@theme_fg_color, 0.10); background: alpha(@theme_fg_color, 0.045); }
            .card > label { line-height: 1.35; }
            .release-entry { min-height: 34px; border-radius: 9px; }
            .pill-button { min-height: 40px; padding-left: 18px; padding-right: 18px; border-radius: 10px; }
            .status-pill { padding: 7px 12px; border-radius: 8px; background: alpha(@accent_color, 0.12); color: @accent_color; font-size: 12px; }
        """.data);
        Gtk.StyleContext.add_provider_for_display (Gdk.Display.get_default (), css, Gtk.STYLE_PROVIDER_PRIORITY_APPLICATION);
    }

    private void launch_builder () {
        try {
            var release = state_version.text.strip ();
            if (release == "") {
                state_version.grab_focus ();
                return;
            }
            // Bubble Tea needs a real terminal. Prefer the Pantheon terminal when
            // available, then fall back to kgx and common terminal emulators.
            string? terminal = null;
            foreach (var candidate in new string[] { "io.elementary.terminal", "kgx", "gnome-terminal", "xterm" }) {
                if (Environment.find_program_in_path (candidate) != null) {
                    terminal = candidate;
                    break;
                }
            }
            if (terminal == null) {
                warning ("No terminal emulator found");
                return;
            }
            string[] argv;
            if (terminal == "gnome-terminal")
                argv = new string[] { terminal, "--", "flakebuilder", "--state-version", release };
            else if (terminal == "xterm")
                argv = new string[] { terminal, "-e", "flakebuilder", "--state-version", release };
            else
                argv = new string[] { terminal, "--", "flakebuilder", "--state-version", release };
            new GLib.Subprocess.newv (argv, GLib.SubprocessFlags.NONE);
        } catch (Error e) {
            warning ("Unable to start flakebuilder: %s", e.message);
        }
    }

    private void open_output () {
        try {
            var path = Path.build_filename (Environment.get_home_dir (), "generated_flakes");
            AppInfo.launch_default_for_uri ("file://" + path, null);
        } catch (Error e) {
            warning ("Unable to open generated flakes: %s", e.message);
        }
    }
}

int main (string[] args) {
    return new FlakebuilderGui ().run (args);
}
