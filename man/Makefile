# Respective MANPATH entries would include /usr/local/share/man or $HOME/man.
# This should be run after make in the ../c directory but before make in the ../docs directory,
# since ../go/mlr is used to autogenerate ./manpage.txt which is used in ../docs.
# See also https://miller.readthedocs.io/en/latest/build.html#creating-a-new-release-for-developers
INSTALLDIR=/usr/local/share/man/man1
INSTALLHOME=$(HOME)/man/man1

# This is normally done only on a development host. Through CI and
# package-installer, mlr.1 should be treated as an already-built artifact,
# needing only to be copied to its install dir.
build: .always
	echo mkman start
	./mkman.rb > mlr.1
	./mkman.rb | groff -man -Tascii | col -b | expand -8 > manpage.txt
	cp manpage.txt ../docs/src/
	echo mkman end

install:
	mkdir -p $(INSTALLDIR)
	cp mlr.1 $(INSTALLDIR)/mlr.1

# ----------------------------------------------------------------
# These targets are only for local dev work.
installhome: build
	mkdir -p $(INSTALLHOME)
	cp mlr.1 $(INSTALLHOME)/mlr.1
# This is an out-of-tree operation so do it only if the user is using out-of-tree manpage setup
# (e.g. their $MANPATH as $HOME/man as a component).
maybeinstallhome: build
	if [ -d $(INSTALLHOME) ]; then cp mlr.1 $(INSTALLHOME)/mlr.1; else echo No $(INSTALLHOME) -- skipping ; fi

# ----------------------------------------------------------------
.always:
	true
