# Respective MANPATH entries would include /usr/local/share/man or $HOME/man.
# This should be run after make in the base directory but before make in the ../docs directory,
# since ../mlr is used to autogenerate ./manpage.txt which is used in ../docs.
# See also ../Makefile and
# https://miller.readthedocs.io/en/latest/build.html#creating-a-new-release-for-developers
PREFIX=/usr/local
INSTALLDIR=$(PREFIX)/share/man/man1

# The manpage.txt pipeline in the build target needs 'pipefail': without it, a
# pipeline's exit status is that of its last command (expand), so a missing or
# failing groff or mkman.rb goes unnoticed and an empty manpage.txt gets copied
# into the docs. 'pipefail' is a bash feature, hence the SHELL setting.
SHELL = bash

# This is normally done only on a development host. Through CI and
# package-installer, mlr.1 should be treated as an already-built artifact,
# needing only to be copied to its install dir.
# Copying to man1/mlr1 is so that you can then put the Miller dir in your
# $MANPATH and do "man mlr".
#
# The groff check gives a clear message up front, since macOS no longer ships
# groff (try 'brew install groff').
build:
	echo mkman start
	groff --version > /dev/null || { echo "man/Makefile: groff is required to build manpage.txt; please install it, e.g. 'brew install groff'." 1>&2; exit 1; }
	./mkman.rb > mlr.1
	set -o pipefail; ./mkman.rb | groff -man -Tascii | col -b | expand -8 > manpage.txt
	cp manpage.txt ../docs/src/
	echo mkman end
	mkdir -p ./man1
	cp mlr.1 ./man1

# DESTDIR is for package installs; nominally blank when this is run interactively.
# See also https://www.gnu.org/prep/standards/html_node/DESTDIR.html
install:
	mkdir -p ./man1
	cp mlr.1 ./man1
	mkdir -p $(DESTDIR)/$(INSTALLDIR)
	cp mlr.1 $(DESTDIR)/$(INSTALLDIR)/mlr.1

.PHONY: build install
