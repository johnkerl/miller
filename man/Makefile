# Respective MANPATH entries would include /usr/local/share/man or $HOME/man.
# This should be run after make in the ../c directory but before make in the ../docs directory,
# since ../go/mlr is used to autogenerate ./manpage.txt which is used in ../docs.
# See also https://miller.readthedocs.io/en/latest/build.html#creating-a-new-release-for-developers
INSTALLDIR=/usr/local/share/man/man1
INSTALLHOME=$(HOME)/man/man1

top: .always
	echo mkman start
	./mkman.rb > mlr.1
	./mkman.rb | groff -man -Tascii | col -b | expand -8 > manpage.txt
	cp manpage.txt ../docs/src/
	echo mkman end

# These targets are only for local dev work.
install: top
	mkdir -p $(INSTALLDIR)
	cp mlr.1 $(INSTALLDIR)/mlr.1
installhome: top
	mkdir -p $(INSTALLHOME)
	cp mlr.1 $(INSTALLHOME)/mlr.1
# This is an out-of-tree operation so do it only if the user is using out-of-tree manpage setup
# (e.g. their $MANPATH as $HOME/man as a component).
maybeinstallhome: top
	if [ -d $(INSTALLHOME) ]; then cp mlr.1 $(INSTALLHOME)/mlr.1; else echo No $(INSTALLHOME) -- skipping ; fi

.always:
	true
