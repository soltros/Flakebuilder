project('flakebuilder-gui', 'vala', version: '0.1.0')
deps = [dependency('gtk4'), dependency('granite-7'), dependency('gio-unix-2.0')]
executable('flakebuilder-gui', 'application.vala', dependencies: deps, install: true)
install_data('io.github.soltros.Flakebuilder.desktop', install_dir: get_option('datadir') / 'applications')
