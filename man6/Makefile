# Respective MANPATH entries would include /usr/local/share/man or $HOME/man.
# This should be run after make in the ../c directory but before make in the ../docs directory,
# since ../go/mlr is used to autogenerate ./manpage.txt which is used in ../docs.
# See also https://miller.readthedocs.io/en/latest/build.html#creating-a-new-release-for-developers
INSTALLDIR=/usr/local/share/man/man1
INSTALLHOME=$(HOME)/man/man1

top: .always
	echo mkman start
	mkman.rb > mlr6.1
	mkman.rb | groff -man -Tascii | col -b | expand -8 > manpage.txt
	echo mkman end

# These targets are only for local dev work.
install:
	mkdir -p $(INSTALLDIR)
	cp mlr6.1 $(INSTALLDIR)/mlr6.1
installhome:
	mkdir -p $(INSTALLHOME)
	cp mlr6.1 $(INSTALLHOME)/mlr6.1

.always:
	true
