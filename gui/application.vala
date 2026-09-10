using Gtk;
using Granite;

public class FlakebuilderGui : Gtk.Application {
    private Gtk.Entry state_version;
    public FlakebuilderGui () {
        Object (application_id: "io.github.soltros.Flakebuilder", flags: ApplicationFlags.DEFAULT_FLAGS);
    }

    protected override void activate () {
        var window = new Gtk.ApplicationWindow (this);
        window.title = "Flakebuilder";
        window.default_width = 760;
        window.default_height = 520;
        window.resizable = true;

        var header = new Gtk.HeaderBar ();
        header.show_title_buttons = true;
        var title = new Gtk.Label ("Flakebuilder");
        title.add_css_class ("title-4");
        header.title_widget = title;
        window.titlebar = header;

        var content = new Gtk.Box (Gtk.Orientation.VERTICAL, 0);
        content.margin_start = 42;
        content.margin_end = 42;
        content.margin_top = 36;
        content.margin_bottom = 36;
        content.spacing = 18;

        var welcome = new Gtk.Label ("Build your NixOS flake");
        welcome.add_css_class ("title-1");
        welcome.halign = Gtk.Align.START;
        content.append (welcome);
        var intro = new Gtk.Label ("Choose configuration bits, shop for packages and NUR sources, then review and generate one flake.nix.");
        intro.wrap = true;
        intro.halign = Gtk.Align.START;
        intro.add_css_class ("dim-label");
        content.append (intro);

        var card = new Gtk.Box (Gtk.Orientation.VERTICAL, 12);
        card.add_css_class ("card");
        card.margin_top = 12;
        card.margin_bottom = 12;
        card.margin_start = 4;
        card.margin_end = 4;
        card.append (new Gtk.Label ("The complete selector is available in the terminal workflow today. This native frontend keeps the same generator, validation and output behavior.") { wrap = true, halign = Gtk.Align.START });
        var state_row = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 10);
        state_row.append (new Gtk.Label ("Original NixOS release") { halign = Gtk.Align.START });
        state_version = new Gtk.Entry ();
        state_version.placeholder_text = "e.g. 26.05";
        state_version.width_chars = 10;
        state_row.append (state_version);
        card.append (state_row);
        content.append (card);

        var actions = new Gtk.Box (Gtk.Orientation.HORIZONTAL, 12);
        var configure = new Gtk.Button.with_label ("Configure flake");
        configure.add_css_class ("suggested-action");
        configure.clicked.connect (() => launch_builder ());
        actions.append (configure);
        var open = new Gtk.Button.with_label ("Open generated flakes");
        open.clicked.connect (() => open_output ());
        actions.append (open);
        content.append (actions);

        var status = new Gtk.Label ("Output directory: ~/generated_flakes");
        status.halign = Gtk.Align.START;
        status.add_css_class ("dim-label");
        content.append (status);
        window.child = content;
        window.present ();
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
