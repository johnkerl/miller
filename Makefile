MANDIR  ?= /usr/share/man
DESTDIR ?=

all: c manpage
devall: c install doc
# TODO: the install target exists to put most-recent mlr executable in the
# path to be picked up by the mlr-execs in the docs dir. better would be to
# export PATH here with ./c at its head.
c: .always
	make -C c top
doc: .always
	cd doc && poki
install: .always
	make -C c install
	install -d -m 0755 $(DESTDIR)/$(mandir)
	install -m 0644 doc/miller.1 $(DESTDIR)/$(mandir)
clean: .always
	make -C c clean
.PHONY: manpage
manpage: doc/miller.1.txt
	( cd doc && a2x a2x -d manpage -f manpage miller.1.txt )
.always:
	@true
